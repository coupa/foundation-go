package middleware

import (
	"github.com/coupa/foundation-go/metrics"
	"github.com/gin-gonic/gin"
)

//MetricsMiddleware is a middleware that will send metrics for every request
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		var timing *metrics.StatsdTiming

		tags := map[string]string{"path": c.Request.URL.Path}
		metrics.Increment("requests", tags)
		timing = metrics.NewTiming("requests", tags)

		c.Next()

		timing.Send()
	}
}
