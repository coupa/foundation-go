package appenv

import (
	"fmt"
	"github.com/coupa/foundation-go/config"
	"io/ioutil"
	"os"
	"strings"
)

const (
	awsRegion        = "AWS_REGION"
	defaultAWSRegion = "us-east-1"
	awsSmName        = "AWSSM_NAME"
	sslSecretName    = "SSL_SECRET_NAME"
)

type AWSEnv struct {
	*BaseEnv
	sslCertPathMap map[string]string
}

func NewAWSEnv(certPath, keyPath string) AppEnvironment {
	return &AWSEnv{
		BaseEnv: &BaseEnv{SSLCertPath: certPath, SSLKeyPath: keyPath},
		sslCertPathMap: map[string]string{
			"app_ssl_certificate":     certPath,
			"app_ssl_certificate_key": keyPath,
		},
	}
}

//LoadEnv for AWS will download environment variables from Secrets Manager and
//set/overwrite those in the current process' environment variables.
func (ae *AWSEnv) LoadEnv() (err error) {
	if os.Getenv(awsRegion) == "" {
		os.Setenv(awsRegion, defaultAWSRegion)
	}
	if s := os.Getenv(awsSmName); s == "" {
		err = fmt.Errorf("Environment '%s' not set", awsSmName)
	} else {
		err = config.WriteSecretsToENV(s)
	}
	if os.Getenv(SSLEnabled) == "true" {
		ae.SSLEnabled = true
	}
	return
}

func (ae *AWSEnv) PrepareSSLCertificate() error {
	snName := os.Getenv(sslSecretName)
	if snName == "" {
		return fmt.Errorf("%s not found", sslSecretName)
	}
	data, _, err := config.GetSecrets(snName)
	if err != nil || len(data) == 0 {
		return fmt.Errorf("Fail to get secret from '%s': %v", snName, err)
	}
	for k, file := range ae.sslCertPathMap {
		if v, ok := data[k]; ok {
			decoded := strings.Replace(v, `\n`, "\n", -1)
			if err = ioutil.WriteFile(file, []byte(decoded), 0644); err != nil {
				return fmt.Errorf("Failed to save file '%s': %v", file, err)
			}
		}
	}
	return nil
}
