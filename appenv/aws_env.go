package appenv

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/coupa/foundation-go/config"
)

const (
	awsRegion        = "AWS_REGION"
	defaultAwsRegion = "us-east-1"
	awsSmName        = "AWSSM_NAME"
	sslSecretName    = "SSL_SECRET_NAME"
)

type AwsEnv struct {
	sslCertFileMap map[string]string
}

func NewAwsEnv() AppEnvironment {
	var ae AwsEnv
	ae.sslCertFileMap = map[string]string{
		"app_ssl_certificate":     SslCertFile,
		"app_ssl_certificate_key": SslKeyFile,
	}
	return ae
}

func (ae AwsEnv) LoadEnv() (err error) {
	s := os.Getenv(awsRegion)
	if s == "" {
		os.Setenv(awsRegion, defaultAwsRegion)
	}
	if s = os.Getenv(awsSmName); s == "" {
		err = fmt.Errorf("Environment '%s' not set", awsSmName)
	} else {
		err = config.WriteSecretsToENV(s)
	} 
	return err
}

func(ae AwsEnv) ConfigureServer(confFile string, conf config.AppConfiguration) (err error) {
	if err = config.LoadJsonConfigFile(confFile, conf); err == nil {
		if conf.IsSslEnabled() {
			err = ae.LoadSSL(conf)
		}
	}
	return err
}

func (ae AwsEnv) LoadSSL(c config.AppConfiguration) error {
	snName := c.GetSslSecretName()
	if snName == "" {
		return fmt.Errorf("Environment parameter 'SSL_SECRET_NAME' not set")
	}
	data, _, err := config.GetSecrets(snName)
	if err != nil || len(data) == 0 {
		return fmt.Errorf("Fail to get secret from '%s': %v", snName, err)
	}
	for k, file := range ae.sslCertFileMap {
		if v, ok := data[k]; ok {
			decoded := strings.Replace(v, `\n`, "\n", -1)
			if err = ioutil.WriteFile(file, []byte(decoded), 0644); err != nil {
				return fmt.Errorf("Failed to save file '%s': %v", file, err)
			}
		}
	}
	return nil
}
