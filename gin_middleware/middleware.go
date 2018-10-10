/**
 * Various middlewares are defined here for metrics and monitoring.
 * e.g. Logging for each API call on webserver, or stats on each API call
 */
package gin_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
	"common-go/logging"
	"common-go/metrics"
)

var REQUEST_ID_HEADER = "X-Request-Id"
var COUPA_INSTANCE_HEADER = "X-COUPA-Instance"

var ID_INDEX_IN_URL_VAR = "Id_Index_In_Url"

// --------------------------------------------------------------------------
func RequestIdMiddleware(c *gin.Context) {
	c.Writer.Header().Set(REQUEST_ID_HEADER, uuid.NewV4().String())
	c.Next()
}

func LoggingMiddleware(c *gin.Context) {
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
		"app":                  logging.LoggingApp,
		"project":              logging.LoggingProject,
		"version":              logging.LoggingAppVersion,
	}).Info("API call")
}

func MetricsMiddleware(c *gin.Context) {
	instance := c.Request.Header.Get(COUPA_INSTANCE_HEADER)
	url := getDelimitedURLForStats(c)

	tags := metrics.Tags{EnterpriseInstanceId: instance, Name: url}
	metrics.StatsIncrementWithTags("events.api.requests", tags)
	metrics.StatsTimeWithTags(c.Next /* Func to invoke */, "transactions.api.latency", tags)
}

// --------------------------------------------------------------------------
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
