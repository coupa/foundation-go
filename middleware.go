/**
 * Various middlewares are defined here for metrics and monitoring.
 * e.g. Logging for each API call on webserver, or stats on each API call
 */
package common_go

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xhandler"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

var REQUEST_ID_HEADER = "X-Request-Id"
var COUPA_INSTANCE_HEADER = "X-COUPA-Instance"

var ID_INDEX_IN_URL_VAR = "Id_Index_In_Url"

func RequestIdMiddlewareFunc(obj ...interface{}) {
	if len(obj) == 0 {
		return // Defensive check
	}

	c, ok := obj[0].(*gin.Context)
	if ok {
		c.Writer.Header().Set(REQUEST_ID_HEADER, uuid.NewV4().String())
		c.Next()
	}
}

func LogMiddlewareFunc(obj ...interface{}) {
	if len(obj) == 0 {
		return // Defensive check
	}

	c, ok := obj[0].(*gin.Context)
	if ok {
		// before request
		t := time.Now()

		c.Next()

		// after request
		latency := time.Since(t)

		// access the status we are sending
		status := c.Writer.Status()
		request_id := c.Writer.Header().Get(REQUEST_ID_HEADER)
		instance := c.Request.Header.Get(COUPA_INSTANCE_HEADER)

		log.WithFields(log.Fields{
			"duration":             latency.Seconds(),
			"status":               status,
			"correlation_id":       request_id,
			"EnterpriseInstanceId": instance,
			"path":                 c.Request.URL.Path,
			"method":               c.Request.Method,
			"remote_ip":            c.Request.RemoteAddr,
			"app":                  LoggingApp,
			"project":              LoggingProject,
			"version":              LoggingAppVersion,
		}).Info("API call")
	}

	// NOTE: In next iteration would like to move this to separate packages alltogether
	handler, ok := obj[0].(xhandler.HandlerC)
	if ok {
		if len(obj) < 3 {
			return // Defensive check
		}
		// before request
		t := time.Now()

		handler.ServeHTTPC(obj[1].(context.Context), obj[2].(http.ResponseWriter), obj[3].(*http.Request))

		log.WithFields(log.Fields{
			"duration":  time.Since(t).Seconds(),
			"path":      c.Request.URL.Path,
			"method":    c.Request.Method,
			"remote_ip": c.Request.RemoteAddr,
			"app":       LoggingApp,
			"project":   LoggingProject,
			"version":   LoggingAppVersion,
		}).Info("API call")
	}
}

func MetricsMiddlewareFunc(obj ...interface{}) {
	if len(obj) == 0 {
		return // Defensive check
	}

	c, ok := obj[0].(*gin.Context)
	if ok {
		instance := c.Request.Header.Get(COUPA_INSTANCE_HEADER)
		url := getDelimitedURLForStats(c)

		tags := Tags{EnterpriseInstanceId: instance, Name: url}
		StatsIncrementWithTags("events.api.requests", tags)
		StatsTimeWithTags(c.Next /* Func to invoke */, "transactions.api.latency", tags)
	}
}

func getDelimitedURLForStats(c *gin.Context) string {
	url := strings.Replace(c.Request.URL.Path, "/", ".", -1)
	if strings.HasPrefix(url, ".") {
		url = url[1:len(url)]
	}

	// If there is an ID in the index, we want to stop generating unique URLs from that point onwards
	idx := c.GetInt(ID_INDEX_IN_URL_VAR)
	if idx > 0 {
		urlSplice := strings.Split(url, ".")
		if len(urlSplice) > idx {
			urlSplice = urlSplice[0:idx] // Do not include index itself
			url = strings.Join(urlSplice, ".")
		}
	}

	return url
}
