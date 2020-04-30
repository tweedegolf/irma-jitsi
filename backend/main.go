package main

import "encoding/json"
import "fmt"
import "io/ioutil"
import "net/http"

import "github.com/privacybydesign/irmago"
import flag "github.com/spf13/pflag"

type Configuration struct {
	ListenAddress         string                             `json:"listen-address,omitempty"`
	IrmaServerURL         string                             `json:"irma-server,omitempty"`
	RoomToAttributes      map[string]irma.AttributeConDisCon `json:"room-map,omitempty"`
	DefaultRoomAttributes *irma.AttributeConDisCon           `json:"default-room,omitempty"`
	BackendName           string                             `json:"backend-name,omitempty"`
	BackendSecret         string                             `json:"backend-secret,omitempty"`
	JitsiName             string                             `json:"jitsi-name,omitempty"`
	JitsiSecret           string                             `json:"jitsi-secret,omitempty"`
	JitsiDomain           string                             `json:"jitsi-domain,omitempty"`
}

func main() {
	var cfg Configuration

	configuration := flag.StringP("config", "c", "", `The file to read configuration from. Further options override.`)
	listenAddress := flag.String("listen-address", "", `The address to listen for external requests, e.g. ":8080".`)
	irmaServer := flag.String("irma-server", "", `The address of the IRMA server to use for disclosure.`)
	roomMap := flag.String("room-map", "", `The map from rooms to attribute condiscons.`)
	defaultRoom := flag.String("default-room", "", `If provided, supplies the attribute condiscons for all unspecified rooms. If not provided, unspecified rooms are not allowed.`)
	backendName := flag.String("backend-name", "", `The name this backend uses to produce JWT (i.e. the 'issuer').`)
	backendSecret := flag.String("backend-secret", "", `The HS256 secret used by this backend to sign & verify own JWT messages.`)
	jitsiSecret := flag.String("jitsi-secret", "", `The HS256 secret used by Jitsi to verify our JWT messages.`)
	jitsiName := flag.String("jitsi-name", "", `The name the Jitsi Authentication Module uses to consume JWT (i.e. the 'audience').`)
	jitsiDomain := flag.String("jitsi-domain", "", `The XMPP domain in use by Jitsi (i.e. the 'subject').`)

	flag.Parse()

	if *configuration != "" {
		contents, err := ioutil.ReadFile(*configuration)
		if err != nil {
			panic(fmt.Sprintf("configuration file not found: %v", *configuration))
		}
		err = json.Unmarshal(contents, &cfg)
		if err != nil {
			panic(fmt.Sprintf("could not parse configuration file: %v", err))
		}
	}

	if *listenAddress != "" {
		cfg.ListenAddress = *listenAddress
	}
	if *irmaServer != "" {
		cfg.IrmaServerURL = *irmaServer
	}
	if *roomMap != "" {
		err := json.Unmarshal([]byte(*roomMap), &cfg.RoomToAttributes)
		if err != nil {
			panic(fmt.Sprintf("could not parse room map: %v", err))
		}
	}
	if *defaultRoom != "" {
		err := json.Unmarshal([]byte(*defaultRoom), &cfg.DefaultRoomAttributes)
		if err != nil {
			panic(fmt.Sprintf("could not parse default room: %v", err))
		}
	}
	if *backendName != "" {
		cfg.BackendName = *backendName
	}
	if *backendSecret != "" {
		cfg.BackendSecret = *backendSecret
	}
	if *jitsiName != "" {
		cfg.JitsiName = *jitsiName
	}
	if *jitsiSecret != "" {
		cfg.JitsiSecret = *jitsiSecret
	}
	if *jitsiDomain != "" {
		cfg.JitsiDomain = *jitsiDomain
	}

	if cfg.ListenAddress == "" {
		panic("option required: listen-address")
	}
	if cfg.IrmaServerURL == "" {
		panic("option required: irma-server")
	}
	if cfg.RoomToAttributes == nil {
		panic("option required: room-map")
	}
	if cfg.BackendName == "" {
		panic("option required: backend-name")
	}
	if cfg.BackendSecret == "" {
		panic("option required: backend-secret")
	}
	if cfg.JitsiName == "" {
		panic("option required: jitsi-name")
	}
	if cfg.JitsiSecret == "" {
		panic("option required: jitsi-secret")
	}
	if cfg.JitsiDomain == "" {
		panic("option required: jitsi-domain")
	}

	externalMux := http.NewServeMux()
	externalMux.HandleFunc("/session", cfg.handleSession)
	externalMux.HandleFunc("/disclose", cfg.handleDisclose)

	externalServer := http.Server{
		Addr:    cfg.ListenAddress,
		Handler: externalMux,
	}
	externalServer.ListenAndServe()
}
