package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var (
	//Regexp for standard environment variable names, which are all uppercase characters
	//and digits connected with underscores. No other symbols or lowercase characters
	//will match.
	regexENVName = regexp.MustCompile(`^[A-Z]([A-Z\d_]*[A-Z\d])?$`)
)

//WriteSecretsToENV will take ENV format keys (all uppercase characters
//and digits connected with underscores) and set them as environment variables.
func WriteSecretsToENV(smName string) error {
	sess := session.Must(session.NewSession())
	return WriteSecretsToENVWithSession(smName, sess)
}

//WriteSecretsToENVWithSession will take ENV format keys (all uppercase characters
//and digits connected with underscores) and set them as environment variables.
//Use this function if you want to provide your own aws session.
func WriteSecretsToENVWithSession(smName string, sess *session.Session) error {
	data, _, err := GetSecretsWithSession(smName, sess)
	if err != nil || len(data) == 0 {
		return err
	}
	for k, v := range data {
		if regexENVName.MatchString(k) {
			os.Setenv(k, v)
		}
	}
	return nil
}

//GetSecrets gets the secrets from AWS secrets manager. If the secrets are a string
//map, it will be returned as the first return value. If the secrets are binary
//data, it will be returned as the second return value.
//
//Use this if your AWS credentials and options are already configured in either
//environment variables or as IAM roles
func GetSecrets(smName string) (map[string]string, []byte, error) {
	sess := session.Must(session.NewSession())
	return GetSecretsWithSession(smName, sess)
}

//GetSecretsWithSession gets the secrets from AWS secrets manager. If the secrets
//are a string map, it will be returned as the first return value. If the secrets
//are binary data, it will be returned as the second return value.
//Use this function if you want to provide your own aws session.
func GetSecretsWithSession(smName string, sess *session.Session) (map[string]string, []byte, error) {
	svc := secretsmanager.New(sess)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(smName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, nil, handleAWSError(err)
	}

	if result.SecretString == nil {
		// The secret is binary
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			return nil, nil, fmt.Errorf("Error decoding binary secret:%v", err)
		}
		return nil, decodedBinarySecretBytes[:len], nil
	}

	var jsonResult map[string]string
	if err = json.Unmarshal([]byte(*result.SecretString), &jsonResult); err != nil {
		return nil, nil, fmt.Errorf("Unable to unmarshall AWS SM secrets: %v", err)
	}
	return jsonResult, nil, nil
}

func handleAWSError(err error) error {
	if aerr, yes := err.(awserr.Error); yes {
		switch aerr.Code() {
		case secretsmanager.ErrCodeDecryptionFailure:
			// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
			err = fmt.Errorf("%s: %s", secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

		case secretsmanager.ErrCodeInternalServiceError:
			// An error occurred on the server side.
			err = fmt.Errorf("%s: %s", secretsmanager.ErrCodeInternalServiceError, aerr.Error())

		case secretsmanager.ErrCodeInvalidParameterException:
			// You provided an invalid value for a parameter.
			err = fmt.Errorf("%s: %s", secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

		case secretsmanager.ErrCodeInvalidRequestException:
			// You provided a parameter value that is not valid for the current state of the resource.
			err = fmt.Errorf("%s: %s", secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

		case secretsmanager.ErrCodeResourceNotFoundException:
			// We can't find the resource that you asked for.
			err = fmt.Errorf("%s: %s", secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
		}
	}
	return err
}
