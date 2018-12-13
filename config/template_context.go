package config

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type TemplateContext struct {
}

// Returns an environment variable
func (c *TemplateContext) Env(key string) string {
	return os.Getenv(key)
}

// Expands environment variables
func (c *TemplateContext) ExpandEnv(s string) string {
	return os.ExpandEnv(s)
}

// Returns the contents of the specified file. Useful for pulling secrets
// that are mounted on the filesystem into the config file.
func (c *TemplateContext) Cat(file string) string {
	bytes, _ := ioutil.ReadFile(file)
	// Ignore error
	return string(bytes)
}

// Returns a Base64 representation of the specified string.
func (c *TemplateContext) ToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Returns a lower-case representation of the specified string.
func (c *TemplateContext) ToLower(s string) string {
	return strings.ToLower(s)
}

// Returns a upper-case representation of the specified string.
func (c *TemplateContext) ToUpper(s string) string {
	return strings.ToUpper(s)
}

// Returns an underscored representation of the specified string.
func (c *TemplateContext) ToUnderscore(s string) string {
	return strings.Replace(s, "-", "_", -1)
}

// Returns a upper-case representation of the specified string.
func (c *TemplateContext) Scanln(prompt string) string {
	var ret string
	fmt.Printf("%s: ", prompt)
	fmt.Scanln(&ret)
	return ret
}

func (c *TemplateContext) Default(def string, data string) string {
	if len(data) > 0 {
		return data
	} else {
		return def
	}
}
