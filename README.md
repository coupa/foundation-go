# Common Golang library for Coupa

This repository hosts common golang code for use in all Golang based Coupa applications.

## Getting Started

Any go IDE would work well for development.

### Prerequisites
* Go version 1.9.* or higher
* (Optional) Docker

## Building and testing

A makefile is added to quickly build code and run tests. The following are the commands respectively:
```
make dist
make test
```

## Usage

Importing the library:
```
import "github.com/coupa/common-go"
```

Metrics and monitoring usage initialization:
```
common_go.InitMetricsMonitoring(<project-name>, <app-name>, <app-version>, <component-name>)
```
NOTE: Project name, App name and Component Name are required arguments

### Middleware support

Few middleware functions are exposed by the library for use in any of the applications that provide web API support.
These functions are initially geared towards use by applications using Go gin web framework (https://github.com/gin-gonic/gin).
Plans are to onboard xMux as well (https://github.com/rs/xmux)

RequestIdMiddlewareFunc:- Adds request ID to each incoming request (if absent). This ID acts as a correlation ID for correlating
info between applications.
```
func mware() gin.HandlerFunc {
	return func(c *gin.Context) { common_go.RequestIdMiddlewareFunc(c) }
}
rengine := gin.New()
rengine.Use(mware())
```

LogMiddlewareFunc:- Emits standardized logs to each web API request, along with the latency info, in JSON format.
```
func mware() gin.HandlerFunc {
	return func(c *gin.Context) { common_go.LogMiddlewareFunc(c) }
}
rengine := gin.New()
rengine.Use(mware())
```

MetricsMiddlewareFunc:- Emits statsD counters to each web API request, along with the latency time durations.
```
func mware() gin.HandlerFunc {
	return func(c *gin.Context) { common_go.MetricsMiddlewareFunc(c) }
}
rengine := gin.New()
rengine.Use(mware())
```

To include all middleware provided by library:
```
func convertToHandler(funcDef func(interface{})) gin.HandlerFunc {
	return func(c *gin.Context) { funcDef(c) }
}
for _, funcDef := range common_go.GetAllMiddlewareFunc() {
    rengine.Use(convertToHandler(funcDef))
}
```

### Logging

Assumption is that application is using Sirupsen for logging. JSON format will be enabled by the library. Example log:
```
	log.WithFields(log.Fields{
		"field1":  val1,
		"field2": val2,
	}).Info("Message")
```

### StatsD Metrics

Following statsD APIs are available to use by the application in it's code:
```
StatsIncrement(key string)
StatsIncrementWithTags(key string, tags Tags)
StatsTime(callback func(), key string)
StatsTimeWithTags(callback func(), key string, tags Tags)
```

The key should of format: <metric-name>.* and should not contain project, application or host info (The same is appended by the library)

### Env variables

* STATSD_HOST: Determines the statsD prefix for the metrics. Default 'app.io.coupahost.com'
* STATSD_URL: The statsD url and port. Useful for e.g. for local dev testing. Defaults to '172.17.0.1:8125'

