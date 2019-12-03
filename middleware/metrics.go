package middleware

import (
	"fmt"

	"github.com/coupa/foundation-go/metrics"
	"github.com/gin-gonic/gin"
)

//Metrics is a middleware that will send metrics for every request
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

//MetricsWithSampleRate is a middleware with adjustable sample rate that will
//send metrics for every request
func MetricsWithSampleRate(rate float32) (gin.HandlerFunc, error) {
	if rate < 0 {
		return nil, fmt.Errorf("Negative metrics sample rate '%f' is invalid. Use the default 1.0", rate)
	}
	return func(c *gin.Context) {
		var timing *metrics.StatsdTiming
		m := metrics.WithSampleRate(rate, map[string]string{"path": c.Request.URL.Path})

		m.Increment("requests")
		timing = m.NewTiming("requests")

		c.Next()

		timing.Send()
	}, nil
}
