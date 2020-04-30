package main

// Note: Although most API calls specify their intended HTTP method, they
// currently accept every HTTP method.

import "encoding/json"
import "fmt"
import "log"
import "net/http"
import "strings"
import "time"

import "github.com/dgrijalva/jwt-go"
import "github.com/privacybydesign/irmago"
import "github.com/privacybydesign/irmago/server"
import _ "golang.org/x/net/websocket"

type DTMF = string
type Secret = string

// The response following the request for a session,
// containing the IRMA session pointer for creating a QR code,
// and some trusted facts used by this backend for continuing.
type SessionResponse struct {
	SessionPtr   *irma.Qr `json:"sessionPtr,omitempty"`
	TrustedFacts string   `json:"trustedFacts,omitempty"`
}

// The response following a disclosure request, containing
// everything to access the Jitsi room with your released credentials.
type DiscloseResponse struct {
	Name string `json:"name,omitempty"`
	Room string `json:"room,omitempty"`
	Jwt  string `json:"jwt,omitempty"`
}

// Trusted facts required by this backend to finalize the session.
// These trusted facts are signed as a JWT.
type SessionTrustedFacts struct {
	Token string `json:"token,omitempty"`
	Room  string `json:"room,omitempty"`
	jwt.StandardClaims
}

// Note do not emit empty fields; Jitsi explicitly states that they are required.
type JitsiUser struct {
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Id     string `json:"id"`
}

// The context for a set of Jitsi claims, detailing the user and group.
type JitsiContext struct {
	Group string    `json:"group"`
	User  JitsiUser `json:"user"`
}

// The set of Jitsi claims used to enter a Jitsi room securely.
type JitsiClaims struct {
	Room    string       `json:"room"`
	Context JitsiContext `json:"context"`
	jwt.StandardClaims
}

// Get the Condiscon for any room from the room mapping.
// If the room does not occur in the room mapping, revert to the default room attributes.
// If the default room attributes are not set, yield an error.
func (cfg Configuration) getAttributesForRoom(room string) (irma.RequestorRequest, error) {
	condiscon, ok := cfg.RoomToAttributes[room]
	if !ok {
		if cfg.DefaultRoomAttributes != nil {
			condiscon = *cfg.DefaultRoomAttributes
		} else {
			return nil, fmt.Errorf("unknown room: %#v", room)
		}
	}

	disclosure := irma.NewDisclosureRequest()
	disclosure.Disclose = condiscon

	request := &irma.ServiceProviderRequest{
		Request: disclosure,
	}

	return request, nil
}

// Request a new session from this backend to sign in to a specific Jitsi room.
// If the `getAttributesForRoom` method yields a corresponding IRMA condiscon,
// a IRMA session is started. For the continuation of the session some trusted
// facts are also passed back along.
func (cfg Configuration) handleSession(w http.ResponseWriter, r *http.Request) {
	room := r.FormValue("room")

	request, err := cfg.getAttributesForRoom(room)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transport := irma.NewHTTPTransport(cfg.IrmaServerURL)
	var pkg server.SessionPackage
	err = transport.Post("session", &pkg, request)
	if err != nil {
		log.Printf("failed to request irma session: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var trustedFacts SessionTrustedFacts
	trustedFacts.Token = pkg.Token
	trustedFacts.Room = room

	trustedFactsJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, trustedFacts)
	trustedFactsJwtString, err := trustedFactsJwt.SignedString([]byte(cfg.BackendSecret))
	if err != nil {
		log.Printf("failed to generate JWT: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var session SessionResponse
	session.SessionPtr = pkg.SessionPtr
	session.TrustedFacts = trustedFactsJwtString

	// Update the request server URL to include the session token.
	transport.Server += fmt.Sprintf("session/%s/", pkg.Token)
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		log.Printf("failed to marshal QR code: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(sessionJSON)
}

// Finalize the session by disclosing the appropriate JWT intended for Jitsi,
// so that we can join a room in an authenticated manner.
//
// Calling this endpoint requires the IRMA session to be already completed succesfully.
//
// Requires the 'trustedFacts' parameter to be set, and will yield the parameters
// required to join the Jitsi room, foremost being a JWT authenticating your person
// with the Jitsi authentication module.
//
// The nickname used to join the Jitsi room consists of a concatenation of the chosen attribute values
// (depending on which condis is chosen by the IRMA app user) separated by spaces.
//
// For specific usecases this can be adapted as part of a fork, or you might want to suggest
// a more robust interface for this to this project.
func (cfg Configuration) handleDisclose(w http.ResponseWriter, r *http.Request) {
	trustedFactsJwtString := r.FormValue("trustedFacts")
	if trustedFactsJwtString == "" {
		http.Error(w, "disclosure needs trustedFacts", http.StatusBadRequest)
		return
	}

	trustedFactsJwt, err := jwt.ParseWithClaims(
		trustedFactsJwtString,
		&SessionTrustedFacts{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.BackendSecret), nil
		})

	if err != nil {
		log.Printf("failed to parse trusted facts: %v", err)
		http.Error(w, "internal server error", http.StatusBadRequest)
		return
	}

	trustedFacts, ok := trustedFactsJwt.Claims.(*SessionTrustedFacts)
	if !ok || !trustedFactsJwt.Valid {
		log.Printf("failed to verify trusted facts: %v", err)
		http.Error(w, "internal server error", http.StatusBadRequest)
		return
	}

	transport := irma.NewHTTPTransport(cfg.IrmaServerURL)
	// Update the request server URL to include the session token.
	transport.Server += fmt.Sprintf("session/%s/", trustedFacts.Token)

	// At this point, the IRMA session is done.
	result := &server.SessionResult{}
	err = transport.Get("result", result)
	if err != nil {
		log.Printf("failed to get irma session result: %v", err)
		return
	}

	status := string(result.Status)
	if status != "DONE" {
		log.Printf("unexpected irma session status: %#v", status)
		return
	}

	// The nickname consists of all condis values concatenated, separated by spaces.
	var nameParts []string
	for _, con := range result.Disclosed {
		for _, attr := range con {
			nameParts = append(nameParts, *attr.RawValue)
		}
	}

	name := strings.Join(nameParts, " ")

	duration, err := time.ParseDuration("1h")
	if err != nil {
		log.Printf("failed to parse ExpiresAt duration: %v", err)
		return
	}

	jitsiClaims := &JitsiClaims{
		Room: trustedFacts.Room,
		Context: JitsiContext{
			Group: "", // Leave empty
			User: JitsiUser{
				Avatar: "", // Leave empty
				Name:   name,
				Email:  "", // Leave empty
				Id:     "", // Leave empty
			},
		},
		StandardClaims: jwt.StandardClaims{
			Audience:  cfg.JitsiName,
			Issuer:    cfg.BackendName,
			Subject:   cfg.JitsiDomain,
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	jitsiJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jitsiClaims)
	jitsiJwtString, err := jitsiJwt.SignedString([]byte(cfg.JitsiSecret))
	if err != nil {
		log.Printf("failed to generate JWT: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var response DiscloseResponse
	response.Room = trustedFacts.Room
	response.Name = name
	response.Jwt = jitsiJwtString

	responseJSON, err := json.Marshal(response)

	w.Write(responseJSON)
}
