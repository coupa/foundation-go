package appenv

import (
	"fmt"
	"os"
	"io/ioutil"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/coupa/foundation-go/config"
)

const (
	localDBPublicKeyFile = "./db_public_key.pem"
)

type LocalEnv struct {
	sslCertFiles []string
	DbPublicKey  *rsa.PublicKey
}

func NewLocalEnv() AppEnvironment {
	var e LocalEnv
	e.sslCertFiles = []string{SslCertFile, SslKeyFile}
	return e
}

func (e LocalEnv) LoadEnv() error {
	return nil
}

func(e LocalEnv) ConfigureServer(confFile string, conf config.AppConfiguration) (err error) {
	if err = config.LoadJsonConfigFile(confFile, conf); err == nil {
		if conf.IsSslEnabled() {
			err = e.LoadSSL(conf)
		}
	}
	return err
}

func (e LocalEnv) LoadSSL(c config.AppConfiguration) error {
	for _, file := range e.sslCertFiles {
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("File '%s' not Found.", file)
			} else {
				return fmt.Errorf("File '%s' not Found. Error: %v", file, err)
			}
		}
	}
	return nil
}

func (e LocalEnv) LoadDbPublicKey(c config.AppConfiguration) (*rsa.PublicKey, error) {
	file := c.GetDbSecretName() // file name
	if file == "" {
		return nil, fmt.Errorf("Environment parameter '%s' not set", dbSecretName)
	}
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("File '%s' not Found.", file)
		} else {
			return nil, fmt.Errorf("File '%s' not Found. Error: %v", file, err)
		}
	}

	data, err := ioutil.ReadFile(file)
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("Fail to get secret from '%s': %v", file, err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("Failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse PKI public key: %v", err)
	}
	if rsaPubKey, ok := pub.(*rsa.PublicKey); ok {
		return rsaPubKey, nil
	} else {
		return nil, fmt.Errorf("Failed to dereference value.")
	}
}