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
	"foundation-go/logging"
	"foundation-go/metrics"
)

var CORRELATION_ID_HEADER = "X-CORRELATION-ID"

// Obsolete instance ID header
var COUPA_INSTANCE_ID_HEADER = "X-COUPA-Instance"
// New instance ID header as per latest standards
var ENTERPRISE_INSTANCE_ID_HEADER = "X-ENTERPRISE-INSTANCE-ID"

// The variable to set on to gin.Context to denote the index
// at which to terminate the URL at. This is useful for URLs where
// we have variable entity IDs being passed e.g. /v1/entity/3445
// For metrics it would be desirable to use /v1/entity rather than /v1/entity/3445
// To achieve the same client should pass index aa 2 via following code:
//     gin.Context.Set(gin_middleware.URL_TERMINATION_INDEX_VAR, 2)
var URL_TERMINATION_INDEX_VAR = "URL_TERMINATION_IDX"

// --------------------------------------------------------------------------
func CorrelationIdMiddleware(c *gin.Context) {
	correlation_id := c.Writer.Header().Get(CORRELATION_ID_HEADER)
	if (correlation_id == "") {
		c.Writer.Header().Set(CORRELATION_ID_HEADER, uuid.NewV4().String())
	}
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
	request_id := c.Writer.Header().Get(CORRELATION_ID_HEADER)
	instance := getInstanceId(c)

	log.WithFields(log.Fields{
		"duration":             latency.Seconds()*1000,
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
	instance := getInstanceId(c)
	url := getDelimitedURLForStats(c)

	tags := metrics.Tags{EnterpriseInstanceId: instance, Name: url}
	metrics.StatsIncrementWithTags("events.api.requests", tags)
	metrics.StatsTimeWithTags(c.Next /* Func to invoke */, "transactions.api.latency", tags)
}

// --------------------------------------------------------------------------
func getInstanceId(c *gin.Context) string {
	instance := c.Request.Header.Get(ENTERPRISE_INSTANCE_ID_HEADER)
	if (instance == "") {
		// Fallback to obsolete instance in case not found
		instance = c.Request.Header.Get(COUPA_INSTANCE_ID_HEADER)
	}
	return instance
}

func getDelimitedURLForStats(c *gin.Context) string {
	url := strings.Replace(c.Request.URL.Path, "/", ".", -1)
	if strings.HasPrefix(url, ".") {
		url = url[1:len(url)]
	}

	// If there is an ID in the index, we want to stop generating unique URLs from that point onwards
	idx := c.GetInt(URL_TERMINATION_INDEX_VAR)
	if idx > 0 {
		urlSplice := strings.Split(url, ".")
		if len(urlSplice) > idx {
			urlSplice = urlSplice[0:idx] // Do not include index itself
			url = strings.Join(urlSplice, ".")
		}
	}

	return url
}
