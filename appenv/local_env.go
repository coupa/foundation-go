package appenv

import (
	"fmt"
	"os"
)

type LocalEnv struct {
	*BaseEnv
}

func NewLocalEnv(certPath, keyPath string) AppEnvironment {
	return &LocalEnv{
		&BaseEnv{SSLCertPath: certPath, SSLKeyPath: keyPath},
	}
}

func (e *LocalEnv) LoadEnv() error {
	if os.Getenv(SSLEnabled) == "true" {
		e.SSLEnabled = true
	}
	return nil
}

func (e *LocalEnv) PrepareSSLCertificate() error {
	for _, file := range e.certs() {
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

func (e *LocalEnv) certs() []string {
	return []string{e.SSLCertPath, e.SSLKeyPath}
}
