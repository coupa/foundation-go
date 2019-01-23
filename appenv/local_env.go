package appenv

import (
	"fmt"
	"os"
	"io/ioutil"
	"crypto/rsa"
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

func (e LocalEnv) LoadSslCertificate() error {
	for _, file := range e.sslCertFiles {
		if err := FileExists(file); err != nil {
			return err
		}
	}
	return nil
}

// Does not support database ssl connection
func (e LocalEnv) LoadDbPublicKey(envSecretFile, notUse string) (*rsa.PublicKey, error) {
	return nil, nil
	file := os.Getenv(envSecretFile)
	if file == "" {
		return nil, fmt.Errorf("'%s' not found", envSecretFile)
	}

	if err := FileExists(file); err != nil {
		return nil, err
	}	
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Fail to read from '%s': %v", file, err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("No key found in '%s'", file)
	}
	return DecodeRsaKey(data)
}

func FileExists(file string) (err error) {
	if _, err = os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("File '%s' not Found", file)
		}
	}
	return
}
