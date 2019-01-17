package appenv

import (
	"fmt"
	"os"
	"github.com/coupa/foundation-go/config"
)

type LocalEnv struct {
	sslCertFiles []string
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
