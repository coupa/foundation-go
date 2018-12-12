# Foundational Golang library for Coupa

This repository hosts foundational golang code for use in all Golang based Coupa applications.

## Getting Started

Any go IDE would work well for development.

### Prerequisites
* Go version 1.9.* or higher
* (Optional) Docker

## Building and testing

A makefile is added to quickly build code (to ensure no compile errors) and run tests. The following are the commands respectively:
```
make dist
make test
```

## Usage

Importing the library:
```
import "github.com/coupa/foundation-go"
```

Metrics and monitoring usage initialization:
```
import "github.com/coupa/foundation-go"

foundation.InitMetricsMonitoring(<project-name>, <app-name>, <app-version>, <component-name>)
```
NOTE: Project name, App name and Component Name are required arguments

### Middleware support

Few middleware functions are exposed by the library for use in any of the applications that provide web API support.
These functions are initially geared towards use by applications using Go gin web framework (https://github.com/gin-gonic/gin).
Plans are to onboard xMux as well (https://github.com/rs/xmux)

RequestIdMiddleware:- Adds request ID to each incoming request (if absent). This ID acts as a correlation ID for correlating
info between applications.
```
import "github.com/coupa/foundation-go/gin_middleware"

rengine := gin.New()
rengine.Use(gin_middleware.RequestIdMiddleware)
```

LoggingMiddleware:- Emits standardized logs to each web API request, along with the latency info, in JSON format.
```
import "github.com/coupa/foundation-go/gin_middleware"

rengine := gin.New()
rengine.Use(gin_middleware.LoggingMiddleware)
```

MetricsMiddleware:- Emits statsD counters to each web API request, along with the latency time durations.
```
import "github.com/coupa/foundation-go/gin_middleware"

rengine := gin.New()
rengine.Use(gin_middleware.MetricsMiddleware)
```

### Logging

Assumption is that application is using Sirupsen for logging. JSON format will be enabled by the library. Example log:
```
import log "github.com/sirupsen/logrus"
import "github.com/coupa/foundation-go/logging"

// Below two lines not needed if using foundation.InitMetricsMonitoring...
logging.InitLogging(<project-name>, <app-name>, <app-version>)
logging.EnableJsonFormat()

log.WithFields(log.Fields{
    "field1":  val1,
    "field2": val2,
}).Info("Message")
```

### StatsD Metrics

Following statsD APIs are available to use by the application in it's code:
```
import "github.com/coupa/foundation-go/metrics"

StatsIncrement(key string)
StatsIncrementWithTags(key string, tags Tags)
StatsTime(callback func(), key string)
StatsTimeWithTags(callback func(), key string, tags Tags)
NewTimingWithTags(tags Tags) (t *statsd.Timing)
StatsTiming(t *statsd.Timing, key string)
TimingWithTagsValue(key string, tags Tags, value interface{})

```

Two ways of using timing:
```
func rates () {
	defer metrics.StatsTiming(metrics.NewTimingWithTags(metrics.Tags{}), "rates")
	...
}

metrics.TimingWithTagsValue("rates"", metrics.Tags{}, int64(my_value)) // my_value in milliseconds

```


The key should of format: <metric-name>.* and should not contain project, application or host info (The same is appended by the library)

### Env variables

* STATSD_INSTANCE_NAME: Determines the statsD prefix for the metrics. Default to 'unknown_instance_name'
* STATSD_URL: The statsD url and port. Useful for e.g. for local dev testing. Defaults to '172.17.0.1:8125'

