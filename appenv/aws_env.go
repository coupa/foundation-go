package appenv

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/coupa/foundation-go/config"
	"crypto/rsa"
)

const (
	awsRegion        = "AWS_REGION"
	defaultAwsRegion = "us-east-1"
	awsSmName        = "AWSSM_NAME"
	sslSecretName    = "SSL_SECRET_NAME"
	
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
		err = fmt.Errorf("'%s' not found", awsSmName)
	} else {
		err = config.WriteSecretsToENV(s)
	} 
	return
}

func (ae AwsEnv) LoadSslCertificate() error {
	snName := os.Getenv(sslSecretName)
	if snName == "" {
		return fmt.Errorf("%s not found", sslSecretName)
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

// There may be multiple public keys for accessing to different endpoints.
func (ae AwsEnv) LoadDbPublicKey(secretEnv, dataKeyEnv string) (*rsa.PublicKey, error) {
	snName := os.Getenv(secretEnv)
	if snName == "" {
		return nil, fmt.Errorf("'%s' not found", secretEnv)
	}
	dataKey := os.Getenv(dataKeyEnv)
	if dataKey == "" {
		return nil, fmt.Errorf("'%s' not found", dataKey)
	}

	data, _, err := config.GetSecrets(snName)
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("Fail to get secret from '%s': %v", snName, err)
	}
	if v, ok := data[dataKey]; ok {
		decoded := strings.Replace(v, `\n`, "\n", -1)
		return DecodeRsaKey([]byte(decoded))
	} else {
		return nil, fmt.Errorf("Failed to find key '%s' from secret '%s'", dataKey, snName)
	}
}
