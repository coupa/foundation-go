package config

import (
	"bytes"
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"fmt"
	"encoding/json"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

// - Generic SSL configuration interace
// 
type AppConfiguration interface {
	IsSslEnabled() bool
	GetSslSecretName() string
	IsDbSslEnabled() bool
	GetDbSecretName() string
}

//- ConfigBytes Functions ------------------------------------------------------
type ConfigBytes []byte

// Reads the specified config file. Note that bedrock.Application will process
// the config file, using text/template, with the following extra functions:
//
//     {{.Env "ENVIRONMENT_VARIABLE"}}
//     {{.Cat "File name"}}
//     {{.Base64 "a string"}}
func ReadConfigFile(file string) (ConfigBytes, error) {
	if _, err := os.Stat(file); err != nil {
		return nil, errors.New("config path not valid")
	}

	tmpl, err := template.New(path.Base(file)).ParseFiles(file)
	if err != nil {
		return nil, err
	}

	var configBytes bytes.Buffer
	tc := TemplateContext{}
	err = tmpl.Execute(&configBytes, &tc)
	if err != nil {
		return nil, err
	}

	return ConfigBytes(configBytes.Bytes()), nil
}

func LoadJsonConfigFile(configFile string, conf AppConfiguration) (err error) {
	var file io.Reader
	if configFile == "" {
		file = strings.NewReader("{}")
	} else if file, err = os.Open(configFile); err != nil {
		return fmt.Errorf("Failed to open file '%s'. Error: %v", configFile, err)
	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(conf); err != nil {
		return fmt.Errorf("Failed to decode file '%s'. Error: %v", configFile, err)
	}
	PopulateEnvConfig(conf)
	return nil
}

func (c ConfigBytes) Unmarshal(dst interface{}) error {
	return yaml.Unmarshal(c, dst)
}

//UnmarshalAt unmarshals a specific key in the config into dst
func (c ConfigBytes) UnmarshalAt(dst interface{}, key string) error {
	var full = make(map[interface{}]interface{})
	if err := c.Unmarshal(&full); err != nil {
		return err
	}
	d, err := yaml.Marshal(full[key])
	if err != nil {
		return err
	}

	return yaml.Unmarshal([]byte(d), dst)
}

//PopulateEnvConfig uses the "env" tag for struct fields to load environment
//variable values into respective struct fields.
func PopulateEnvConfig(c interface{}) {
	configType := reflect.TypeOf(c).Elem()
	configValue := reflect.ValueOf(c).Elem()

	for i := 0; i < configType.NumField(); i++ {
		configField := configType.Field(i)
		envValue := os.Getenv(configField.Tag.Get("env"))

		if envValue != "" {
			switch configValue.Field(i).Type().String() {
			case "bool":
				v, _ := strconv.ParseBool(envValue)
				configValue.FieldByName(configField.Name).SetBool(v)
			case "int", "int64":
				v, _ := strconv.ParseInt(envValue, 10, 64)
				configValue.FieldByName(configField.Name).SetInt(v)
			case "float32":
				v, _ := strconv.ParseFloat(envValue, 32)
				configValue.FieldByName(configField.Name).SetFloat(v)
			case "float64":
				v, _ := strconv.ParseFloat(envValue, 64)
				configValue.FieldByName(configField.Name).SetFloat(v)
			default:
				configValue.FieldByName(configField.Name).SetString(envValue)
			}
		}
	}
}

//SplitByCommaSpace converts strings separated by comma's into a string slice.
//It can be used to convert comma-separated env values into a string slice.
func SplitByCommaSpace(s string) []string {
	f := func(c rune) bool {
		return c == ',' || c == ' '
	}
	return strings.FieldsFunc(s, f)
}
