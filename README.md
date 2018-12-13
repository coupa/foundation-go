# Foundational Go library

This foundational Go library conforms to Coupa's microservice standard and can be used in Go based applications.

## Getting Started

### Prerequisites
* Go version 1.10.* or higher

## Structure
Foundation lets you set up your application to use logging, health checks, and metrics conforming to the microservice standard.

To set up logging, you need to call `logging.InitStandardLogger` early in your application. To set up statsd metrics, you need to define a factory function `func() *metrics.Statsd` and set it using `metrics.SetFactory` function (please see the details in the example below).

Foundation also provides Gin-backed (https://github.com/gin-gonic/gin) server, several middlewares, and health check mechanisms to support the microservice standard.

An example of usage:
```
import (
  "github.com/coupa/foundation-go/health"
  "github.com/coupa/foundation-go/logging"
  "github.com/coupa/foundation-go/metrics"
  "github.com/coupa/foundation-go/middleware"
  "github.com/coupa/foundation-go/server"
  "github.com/gin-gonic/gin"
  log "github.com/sirupsen/logrus"
)

func main() {

  //***************************** Logging ***********************************

  //logging.InitStandardLogger must be called early in your app.
  logging.InitStandardLogger("v1.0.0")

  //The InitStandardLogger above will make logrus' standard logger to use the standard format
  //So this log.Info here will have the required fields for the standard.
  log.Info("Starting the server on :8080")

  //***************************** Metrics ***********************************

  //These enable metrics to know how to initialize the Statsd client
  factory := func() *metrics.Statsd {
    return metrics.NewStatsd("address", "prefix.", "version", "app", 1.0)
  }
  metrics.SetFactory(factory)

  //****************************** Server ***********************************

	svr := server.Server{
		Engine:               gin.New(),
		AppInfo:              &health.AppInfo{...},
		ProjectInfo:          &health.ProjectInfo{...},
    AdditionalHealthData: map[string]*health.AdditionalHealthData{},
	}

  //*** Middlewares ***

  //Register the middlewares first before registering routes
  svr.UseMiddleware(middleware.Correlation())
  svr.UseMiddleware(middleware.Metrics())
  svr.UseMiddleware(middleware.RequestLogger(false))

  //*** Health ***

  //Declare or create the health checks that you want
  dbCheck := health.SQLCheck{
  	Name: "mysql",
  	Type: "internal",
    DB: ...,  //Some *sql.DB
  }
  serviceCheck1 := health.WebCheck{
  	Name: "some web 1",
  	Type: "service",
  	URL:  "https://some.web/health",
  }
  serviceCheck2 := health.WebCheck{
  	Name: "some web 2",
  	Type: "service",
  	URL:  "https://some.web2/health",
  }

  //Register 3 versions of the detailed health. Note that they are different as "/v1" has additional custom data but does not have `serviceCheck2` dependency check.
  ahd1 := health.AdditionalHealthData{
    DependencyChecks: []HealthCheck{dbCheck, serviceCheck1},
    DataProvider:    func(c *gin.Context) map[string]interface{}{
      return map[string]interface{}{
        "custom": "data",
      }
    },
  }

  adh2 := health.AdditionalHealthData{
    DependencyChecks: []HealthCheck{dbCheck, serviceCheck1, serviceCheck2},
  }

  svr.RegisterDetailedHealth("/v1", "v1 of app detailed health", ahd1)
  svr.RegisterDetailedHealth("/v2", "v2 of app detailed health", ahd2)
  svr.RegisterDetailedHealth("/v3", "v3 without custom data or dependency check", nil)
  svr.RegisterSimpleHealth()

  someHandler := func(c *gin.Context) {
    //Emitting statsd metric
    metrics.Increment("interesting.metric")

    logging.RL(c).Info("Logging with correlation ID")
    c.JSON(200, `{"he":"llo"}`)
  }

  //Register routes
  svr.Engine.GET("/test", someHandler)

  svr.Engine.Run(":80") //svr.Engine.Run() without address parameter will run on ":8080"
}
```

### Env variables

`health.AppInfo{}.FillFromENV` and `health.ProjectInfo{}.FillFromENV` will by default load from these Env variables:

* APPLICATION_NAME
* HOSTNAME
* SERVER_BIND_ADDRESS
* PROJECT_REPO
* PROJECT_HOME
* PROJECT_OWNERS: Comma-separated owner names.
* PROJECT_LOG_URLS: Comma-separated log URLs.
* PROJECT_STATS_URLS: Comma-separated metrics URLs.
