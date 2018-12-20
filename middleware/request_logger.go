package middleware

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	userIDHeader        = "X-USER-ID"
	clientNameHeader    = "X-CLIENT-NAME"
	clientVersionHeader = "X-CLIENT-VERSION"
)

func RequestLogger(excludeSimpleHealth bool, excludePathsInRegexp ...string) gin.HandlerFunc {
	//Compile the exclusion regexps
	var regexps []*regexp.Regexp
	for _, re := range excludePathsInRegexp {
		if regex, err := regexp.Compile(re); err == nil {
			regexps = append(regexps, regex)
		}
	}

	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t).Seconds()

		if c.Request.URL.Path == "/health" && excludeSimpleHealth || matches(regexps, c.Request.URL.Path) {
			return
		}

		fields := log.Fields{
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   latency,
			"method":     c.Request.Method,
			"remote_ip":  c.ClientIP(),
			"parameters": c.Request.URL.Query(),
		}
		value := c.Writer.Header().Get(correlationHeader)
		if value != "" {
			fields["correlation_id"] = value
		}
		value = c.Request.Header.Get(userIDHeader)
		if value != "" {
			fields["user_id"] = value
		}
		value = c.Request.Header.Get(clientNameHeader)
		if value != "" {
			fields["client"] = value
		}
		value = c.Request.Header.Get(clientVersionHeader)
		if value != "" {
			fields["client_version"] = value
		}

		log.WithFields(fields).Info("")
	}
}

func matches(regexps []*regexp.Regexp, str string) bool {
	for _, regex := range regexps {
		if regex.MatchString(str) {
			return true
		}
	}
	return false
}
