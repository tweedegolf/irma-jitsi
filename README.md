# IRMA + Jitsi

Enter a [Jitsi](https://jitsi.org/) videoconference securely with an immutable authenticated identity.
The authentication is done using [IRMA](https://irma.app), a digital passport on your phone.

Visit [jitsi-demo.tweede.golf](https://jitsi-demo.tweede.golf/) for a short demo.

This project is partially funded using public european funds. Therefore this project is licensed under the EUPLv1.2.

## Structure

This project requires 4 components:
* An [IRMA](https://irma.app) server with the private-facing socket accessible to the backend. (JWT authentication is not (yet) implemented for the backend)
* A [Jitsi](https://jitsi.org/) server running the [Token Authentication module](https://github.com/jitsi/lib-jitsi-meet/blob/master/doc/tokens.md) and a small custom module found in `/modules`.
* A backend server that takes the IRMA response and subsequently yields a Jitsi authentication token.
* A frontend website to allow users to start an IRMA session and enter the Jitsi room.

Essentially what this project contributes to Jitsi is a small module (specifically for the Prosody XMPP server) that enforces the nickname of any room attendees, and a process to craft authentication tokens using IRMA to securely determine the proper nickname.

## Production

This guide does not help you to install the project for public use. Merge requests for guides on how to deploy the backend and frontend are welcome. See the [Jitsi server setup section](#Jitsi%20server%20setup) on how to set up your Jitsi server with the appropriate modules configured.

## Development

**Note**: in order to run this application in development you unfortunately require a fully running Jitsi instance with the [Token Authentication module](https://github.com/jitsi/lib-jitsi-meet/blob/master/doc/tokens.md) already running, and an IRMA Go server also fully running. Merge requests to get a Jitsi instance configured thusly running in docker are welcome. You may use the following non-authenticated IRMA instance: `https://irma-noauth.demo.sarif.nl`. (Please be mindful that you do not disclose production attributes unless you do not mind) Tweede golf also provides a Jitsi instance with a publicly known JWT keyphrase, which is already configured in `docker-compose.example.yml`. Please do not use this instance for anything other than a demo.

### Requirements

In order to run the development project you require `docker` (version 19.03) and `docker-compose` (version 1.25) installed. It has been tested to work on any plain Ubuntu 18.04.1 LTS.
Before continuing your will require to configure your application to use the aforementioned IRMA server and Jitsi instance.

### Setup

Copy `docker-compose.example.yml` to `docker-compose.yml` and change the various environment variable and most notably the command-line arguments for the backend to their proper values.

Then run:
```bash
bin/setup.sh
```

This will build the docker images and download any runtime dependencies for both the backend (Go) and frontend (Javascript with yarn).

### Running

```bash
bin/up.sh
```

This starts the frontend, backend, and an nginx server using self-signed certificates. Ensure that the self-signed certificates of both domains are accepted in your browser:
* https://veiligjitsi.test.tweede.golf/
* https://backend.veiligjitsi.test.tweede.golf/

Then visit https://veiligjitsi.test.tweede.golf/ and press the button. Release your IRMA credentials to the server to proceed.

You will need the configured demo attributes to proceed with this demo. Fill in the form on the bottom of these pages and scan the QR codes using your IRMA app:
* [irma-demo.MijnOverheid.fullName.firstname](https://privacybydesign.foundation/attribute-index/en/irma-demo.MijnOverheid.fullName.html)
* [irma-demo.MijnOverheid.fullName.familyname](https://privacybydesign.foundation/attribute-index/en/irma-demo.MijnOverheid.fullName.html)
* [irma-demo.MijnOverheid.birthCertificate.dateofbirth](https://privacybydesign.foundation/attribute-index/en/irma-demo.MijnOverheid.birthCertificate.html)

### Backend commandline options

* `--config`: the file to read configuration from. Further options override these values.
* `--listen-address`: the address to listen for external requests, e.g. `:8080`.
* `--irma-server`: the address of the IRMA server to use for disclosure, e.g. https://irma-noauth.demo.sarif.nl. You should really set up your own IRMA server, but you can use this IRMA server for development. Make sure that you do not use production attributes during development.
* `--room-map`: the mapping from rooms to attribute [condiscons](https://irma.app/docs/condiscon/), e.g. ```{"roomName": [[["irma-demo.MijnOverheid.fullName.firstname"]]]}```
* `--default-room`: if provided, supplies the attribute condiscons for all unspecified rooms. If not provided, unspecified rooms are not allowed.
* `--backend-name`: the name this backend uses to produce JWT (i.e. the 'issuer'). At Jitsi this corresponds to `asap_accepted_issuers` in `/etc/prosody/prosody.cfg.lua`.
* `--backend-secret`: the HS256 secret used by this backend to sign & verify own JWT messages. This can be anything you choose.
* `--jitsi-secret`: the HS256 secret used by Jitsi to verify our JWT messages. This corresponds to `app_secret` in `/etc/prosody/conf.d/HOSTNAME.cfg.lua`.
* `--jitsi-name`: the name the Jitsi Authentication Module uses to consume JWT (i.e. the 'audience'). At Jitsi this corresponds to `asap_accepted_audiences` in `/etc/prosody/prosody.cfg.lua` and `app_id` in `/etc/prosody/conf.d/HOSTNAME.cfg.lua`.
* `--jitsi-domain`: the XMPP domain in use by Jitsi (i.e. the 'subject'). Configured in `/etc/jitsi/jicofo/config` as `JICOFO_HOSTNAME`.

### Jitsi server setup

Starting with a blank Debian Buster (10) machine with a public IP address, hostname and disclosing:
* 80, 443 TCP for HTTP & HTTPS
* 4443, 5347 TCP for Jitsi control
* 10000 - 20000 UDP for Jitsi RTC audio & video

We basically follow the [Jitsi Meet quick install guide](https://github.com/jitsi/jitsi-meet/blob/master/doc/quick-install.md) and the [Token module guide](https://github.com/jitsi/lib-jitsi-meet/blob/master/doc/tokens.md), and will specify here where we differ.

1. Set up FQDN in `/etc/hosts` by adding `127.0.0.1 localhost meet.example.org`.
2. Add the Jitsi package repository by executing:
    ```bash
    echo 'deb https://download.jitsi.org stable/' >> /etc/apt/sources.list.d/jitsi-stable.list
    wget -qO -  https://download.jitsi.org/jitsi-key.gpg.key | sudo apt-key add -
    ```
3. Install Jitsi meet:
    ```bash
    # Ensure support is available for apt repositories served via HTTPS
    apt-get install apt-transport-https

    # Retrieve the latest package versions across all repositories
    apt-get update

    # Perform jitsi-meet installation
    apt-get -y install jitsi-meet
    ```
    During installation you will be asked to enter your FQDN, and whether you want to generate a new certificate. That is what I want, and as such will continue as follows.
4. Generate a certificate:
    ```bash
    /usr/share/jitsi-meet/scripts/install-letsencrypt-cert.sh
    ```
5. Set up STUN for users that are behind NAT (you certainly are at least some of the time).
    Change `/etc/jitsi/videobridge/sip-communicator.properties` to include:
    ```
    org.ice4j.ice.harvest.NAT_HARVESTER_LOCAL_ADDRESS=<Local.IP.Address>
    org.ice4j.ice.harvest.NAT_HARVESTER_PUBLIC_ADDRESS=<Public.IP.Address>
    ```
6. Confirm that the server is working as a standard Jitsi server by visiting your domain.
7. Install the Token authentication module:
    ```bash
    apt-get install jitsi-meet-tokens
    ```
    It will ask you for an application ID and secret, but this appears not to do anything.
8. Open `/etc/prosody/prosody.cfg.lua` and
    * Set `plugin_paths` to `/usr/share/jitsi-meet/prosody-plugins/`.
    * Add and replace `BACKEND_NAME` and `JITSI_NAME` with the values as passed to the backend:
    ```
    asap_accepted_issuers = { "BACKEND_NAME" }
    asap_accepted_audiences = { "JITSI_NAME" }
    ```
    * Set `c2s_require_encryption` to `false`.
9. Open `/etc/prosody/conf.d/HOSTNAME.cfg.lua` and
    * Under the virtualhost, change authentication from `anonymous` to `token`.
    * And set the aforementioned `app_id` and `app_secret`. These will correspond to the backend `--jitsi-name` and `--jitsi-secret`.
    * And add `presence_identity` to `modules_enabled`. This is the module exposing the identity data.
    * And add `drop_namechanges` to `modules_enabled`. This is our module enforcing the nicknames.
    * Under `Component "conference.HOSTNAME" "muc"` enable `token_verification` at `modules_enabled`. This is the token authentication module.
10. Copy `mod_drop_namechanges.lua` in the `modules` folder of this repository to `/usr/share/jitsi-meet/prosody-plugins`.
11. Restart prosody:
    ```bash
    systemctl restart prosody
    ```
12. Read the prosody log at `/var/log/prosody/prosody.err` to confirm that it does not start due to a missing LUA dependency:
    ```
    modulemanager   error   Error initializing module 'token_verification' on 'conference.HOSTNAME': /usr/lib/prosody/util/startup.lua:144: module 'basexx' not found:No LuaRocks module found for basexx
    ```
13. Install some dependencies:
    ```
    apt purge liblua5.1 liblua5.1-dev
    apt install git cmake liblua5.2 liblua5.2-dev
    ```
13. Build and install the `luacrypto` dependency:
    ```
    cd /tmp
    git clone https://github.com/Wassasin/luacrypto.git
    luarocks make
    ```
13. Install the LUA dependencies that are OK:
    ```bash
    luarocks install lua-cjson 2.1.0-1
    luarocks install luajwtjitsi
    luarocks install basexxx
    ```
14. Restart prosody and notice that it is running:
    ```bash
    systemctl restart prosody
    systemctl status prosody
    ```

You can now start using authenticated Jitsi.