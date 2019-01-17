package appenv

import (
	"fmt"
	"os"
)

var (
	// Local certificate files
	CloudProvider = "CLOUD_PROVIDER"
	AWSProvider   = "AWS"
	LocalProvider = "LOCAL"
	SSLEnabled    = "SSL_ENABLED"

	DefaultSSLCertPath = "./foundation-server.crt"
	DefaultSSLKeyPath  = "./foundation-server.key"
)

type AppEnvironment interface {
	LoadEnv() error
	PrepareSSLCertificate() error
	GetSSLCertPath() string
	GetSSLKeyPath() string
	IsSSLEnabled() bool
}

func NewAppEnv(provider string) (appEnv AppEnvironment, err error) {
	return NewAppEnvWithCerts(provider, "", "")
}

//NewAppEnvWithCerts creates an AppEnvironment based on the provider and accepts
//customer file paths for SSL certs. If the cert parameters are empty strings,
//the default cert file paths will be used.
func NewAppEnvWithCerts(provider, sslCertPath, sslKeyPath string) (appEnv AppEnvironment, err error) {
	if sslCertPath == "" {
		sslCertPath = DefaultSSLCertPath
	}
	if sslKeyPath == "" {
		sslKeyPath = DefaultSSLKeyPath
	}
	if provider == "" {
		provider = os.Getenv(CloudProvider)
	}
	switch provider {
	case LocalProvider:
		appEnv = NewLocalEnv(sslCertPath, sslKeyPath)
	case AWSProvider:
		appEnv = NewAWSEnv(sslCertPath, sslKeyPath)
	default:
		err = fmt.Errorf("Unknown cloud provider: '%s'", provider)
	}
	return appEnv, err
}

type BaseEnv struct {
	SSLEnabled  bool
	SSLCertPath string
	SSLKeyPath  string
}

func (be *BaseEnv) GetSSLCertPath() string {
	return be.SSLCertPath
}

func (be *BaseEnv) GetSSLKeyPath() string {
	return be.SSLKeyPath
}

func (be *BaseEnv) IsSSLEnabled() bool {
	return be.SSLEnabled
}
