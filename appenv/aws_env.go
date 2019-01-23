package appenv

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/coupa/foundation-go/config"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const (
	awsRegion        = "AWS_REGION"
	defaultAwsRegion = "us-east-1"
	awsSmName        = "AWSSM_NAME"
	sslSecretName    = "SSL_SECRET_NAME"
	dbSecretName     = "DB_SECRET_NAME"
	awsDBPublictKey  = "rds_sslca"
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
		return fmt.Errorf("Environment parameter '%s' not set", sslSecretName)
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

func (ae AwsEnv) LoadDbPublicKey(c config.AppConfiguration) (*rsa.PublicKey, error) {
	snName := c.GetDbSecretName()
	if snName == "" {
		return nil, fmt.Errorf("Environment parameter '%s' not set", dbSecretName)
	}
	data, _, err := config.GetSecrets(snName)
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("Fail to get secret from '%s': %v", snName, err)
	}
	if v, ok := data[awsDBPublictKey]; ok {
		decoded := strings.Replace(v, `\n`, "\n", -1)
		block, _ := pem.Decode([]byte(decoded))
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
	} else {
		return nil, fmt.Errorf("Failed to find key '%s' from secret '%s'", awsDBPublictKey, snName)
	}
}
