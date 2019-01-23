package appenv

import (
	"fmt"
	"crypto/rsa"
	)

const (
	// Local certificate files
	CloudProvider = "CLOUD_PROVIDER"
	AwsProvider   = "AWS"
	LocalProvider = "LOCAL"

	defaultSslCertFile = "./server.crt"
	defaultSslKeyFile  = "./server.key"

	dbSecretName     = "DB_SECRET_NAME"

	)

var (
	SslCertFile = defaultSslCertFile
	SslKeyFile  = defaultSslKeyFile
)

type AppEnvironment interface {
	LoadEnv() error
	LoadSslCertificate() error
	LoadDbPublicKey(entryEnv, keyEnv string) (*rsa.PublicKey, error)
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
