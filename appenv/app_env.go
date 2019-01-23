package appenv

import (
	"fmt"
	"crypto/rsa"
	"github.com/coupa/foundation-go/config"
	)

const (
	// Local certificate files
	CloudProvider = "CLOUD_PROVIDER"
	AwsProvider   = "AWS"
	LocalProvider = "LOCAL"

	defaultSslCertFile = "./server.crt"
	defaultSslKeyFile  = "./server.key"
	)

var (
	SslCertFile = defaultSslCertFile
	SslKeyFile  = defaultSslKeyFile
)

type AppEnvironment interface {
	LoadEnv() error
	ConfigureServer(confFile string, conf config.AppConfiguration) error
	LoadSSL(config.AppConfiguration) error
	LoadDbPublicKey(config.AppConfiguration) (*rsa.PublicKey, error)
}

func NewAppEnv(provider string) (appEnv AppEnvironment, err error) {
	switch provider {
	case LocalProvider:
		appEnv =  NewLocalEnv()
	case AwsProvider:
		appEnv =  NewAwsEnv()
	default:
		err = fmt.Errorf("Unknown cloud provider: '%s'", provider)
	}
	return appEnv, err
}
