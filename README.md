# Foundational Go library

This foundational Go library conforms to Coupa's microservice standard and can be used in Go based applications.

## Getting Started

### Prerequisites
* Go version 1.10.* or higher

## Structure
Foundation lets you set up your application to use logging, health checks, and metrics conforming to the microservice standard.

To set up logging, you need to call `logging.InitStandardLogger` early in your application. To set up statsd metrics, you need to define a factory function `func() *metrics.Statsd` and set it using `metrics.SetFactory` function (please see the details in the example below).

Foundation also provides Gin-backed (https://github.com/gin-gonic/gin) server, several middlewares, and health check mechanisms to support the microservice standard.

### Logging
```
import (
  "github.com/coupa/foundation-go/logging"
  log "github.com/sirupsen/logrus"
)

func main() {
  //logging.InitStandardLogger must be called early in your app.
  logging.InitStandardLogger("v1.0.0")

  //The InitStandardLogger above will make logrus' standard logger use the standard format
  //So any call on logrus package (log) will have the standard fields.
  log.Info("this log will have the required standard fields in JSON format")
}
```
### Metrics (Statsd)
The design of the metrics package is to simplify the code of adding required and custom tags, while setting the standard measurement names (events/transactions) automatically.
```
import (
  "github.com/coupa/foundation-go/metrics"
)

func main() {
  //Setting the metrics factory so that it can initialize the Statsd client
  factory := func() *metrics.Statsd {
    return metrics.NewStatsd("address", "prefix.", "version", "app", 1.0)
  }
  metrics.SetFactory(factory)

  //"interesting.metric" will become the "name" tag of the metric
  metrics.Increment("interesting.metric")

  //Add custom tags to the metric. All metrics methods can take optional maps as custom tags.
  tags := map[string]string{"tag1": "value1", "tag2": "value2"}
  metrics.Increment("with.custom.tags", tags)
}

func someFunc() {
  //This measures the execution time of this function and emits the metric at the end of the method.
  defer metrics.NewTiming("timing.someFunc").Send()

  t := metrics.NewTiming("timing.particular.code")
  ...
  t.Send()
}
```
### Secrets Manager
```
import (
  "github.com/coupa/foundation-go/config"
)

func main() {
  //To use the GetSecrets or WriteSecretsToENV functions in config/aws_secrets_manager.go,
  //you need to have AWS credentials configured as environment variables or configure
  //AWS EC2 IAM roles to allow access (see https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).
  //If you have other ways of authentication, you can provide your own session with
  //config.GetSecretsWithSession or config.WriteSecretsToENVWithSession.

  //GetSecrets will download all keys and values from the specific secret name
  //`secrets` is of map[string]string
  secrets, err := config.GetSecrets("dev/application/some_secret")

  //WriteSecretsToENV will grab values from secret whose key names are in the format of an
  //environment variable name (all uppercase characters with digits and connected
  //with underscores. No dash or other symbols) and set them to corresponding
  //environment variables.
  err = config.WriteSecretsToENV("dev/application/for_testing")

  //For example, if "dev/application/for_testing" has {"GOOD_ENV1":"good","bad-env2":"bad"},
  //WriteSecretsToENV will set only GOOD_ENV1="good" and not the other one
}
```
### Server, Middlewares, and Health Checks
```
import (
  "github.com/coupa/foundation-go/health"
  "github.com/coupa/foundation-go/logging"
  "github.com/coupa/foundation-go/middleware"
  "github.com/coupa/foundation-go/server"
  "github.com/gin-gonic/gin"
)

func main() {
  svr := server.Server{
    Engine:               gin.New(),
    AppInfo:              &health.AppInfo{...},     //You should fill in this info
    ProjectInfo:          &health.ProjectInfo{...}, //You should fill in this info
    AdditionalHealthData: map[string]*health.AdditionalHealthData{},
  }

  //*** Middlewares ***

  //Register the middlewares first before registering routes

  //This forwards or generates correlation ID in the response header
  svr.UseMiddleware(middleware.Correlation())

  //This enables metrics collection of request counting and timing
  svr.UseMiddleware(middleware.Metrics())

  //This writes access logs in the microservice standard format.
  //If the parameter is true: `middleware.RequestLogger(true)`, it won't write access logs
  //for simple health check requests, which some applications may prefer this
  //in order to have a cleaner log
  svr.UseMiddleware(middleware.RequestLogger(false))

  //*** Health ***

  //Create the detailed health checks that your service depends on
  dbCheck := health.SQLCheck{
    Name: "mysql",
    Type: "internal",
    DB: ...,      //Some *sql.DB
  }
  serviceCheck := health.WebCheck{
  	Name: "some web 1",
  	Type: "service",
  	URL:  "https://some.web/health", //If the target is a Coupa service, make sure to use the "simple health" endpoint
  }
  redisCheck := health.RedisCheck{
  	Name: "redis",
  	Type: "internal",
  	Client: ...,    //Some *redis.Client
  }

  ahdV2 := health.AdditionalHealthData{
    DependencyChecks: []HealthChecker{dbCheck, serviceCheck},
    DataProvider:    func(c *gin.Context) map[string]interface{}{
      return map[string]interface{}{
        "custom": "data",
      }
    },
  }

  adhV3 := health.AdditionalHealthData{
    DependencyChecks: []HealthChecker{dbCheck, serviceCheck, redisCheck},
  }

  //Register 3 versions of the detailed health. Note that they are different as
  //"/v1" has additional custom data but does not have `serviceCheck2` dependency check,
  //and "/v3" has no dependency or additional data.
  svr.RegisterDetailedHealth("/v1", "v1 without custom data or dependency check", nil)
  svr.RegisterDetailedHealth("/v2", "v2 with mysql and service checks", adhV2)
  svr.RegisterDetailedHealth("/v3", "v3 with mysql, service, and redis checks", adhV3)

  //Simple health is at /health
  svr.RegisterSimpleHealth()

  someHandler := func(c *gin.Context) {
    logging.RL(c).Info("This is logging with correlation ID")

    c.JSON(200, `{"he":"llo"}`)
  }

  //Register routes
  svr.Engine.GET("/test", someHandler)

  svr.Engine.Run(":80") //This will run on port 80
  //svr.Engine.Run() without address parameter will run on ":8080"
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

### Enabling Secrets Manager Tests

This is for foundation-go contributors to run secrets manager tests in foundation-go.

Foundation-go has tests that actually connect to AWS secrets manager. These tests are not run in regular `go test ./...` command.

To run these tests, you need to pre-setup the AWS credentials (mostly likely in environment variables or config files, see https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html), make sure the target secret of name `dev/application/for_testing` exists and have:
```
{
  "TEST_DESCRIPTION": "For testing",
  "not_set_to_env":"Will not be set to your ENV variables"
}
```
as data. Then run test with this environment variable and command: `TEST_SECRETS_MANAGER=true go test ./...`.

### Enabling Redis Tests

This is for foundation-go contributors to run redis tests in foundation-go.

Foundation-go has tests that actually connect to a local Redis server. These tests are not run in regular `go test ./...` command.

To run these tests, you need to run a local Redis server *without password* on the port 6379, then run test with this environment variable and command: `TEST_REDIS=true go test ./...`.
