module github.com/tweedegolf/irmabellen/backend

go 1.14

require (
	github.com/bwesterb/go-atum v1.0.3 // indirect
	github.com/certifi/gocertifi v0.0.0-20200211180108-c7c1fbc02894 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fxamacker/cbor v1.5.1 // indirect
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.4 // indirect
	github.com/jasonlvhit/gocron v0.0.0-20191228163020-98b59b546dee // indirect
	github.com/jinzhu/gorm v1.9.12 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/privacybydesign/gabi v0.0.0-20200306134149-18dd7a01d765 // indirect
	github.com/privacybydesign/irmago v0.0.0-20200306135745-b0ce9a706e71
	github.com/spf13/pflag v1.0.4-0.20190111213756-a45bfec10d59
	github.com/timshannon/bolthold v0.0.0-20200308034358-09aaf76b2c32 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
)

replace astuart.co/go-sse => github.com/sietseringers/go-sse v0.0.0-20200223201439-6cc042ab6f6d

replace github.com/spf13/pflag => github.com/sietseringers/pflag v1.0.4-0.20190111213756-a45bfec10d59

replace github.com/spf13/viper => github.com/sietseringers/viper v1.0.1-0.20200205174444-d996804203c7
