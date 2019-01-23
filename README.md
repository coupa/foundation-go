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
  "github.com/coupa/foundation-go/config"
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
  log.Info("this log will have the required standard fields")

  //************************* Secrets Manager *******************************

  //To use the GetSecrets or WriteSecretsToENV functions in config/aws_secrets_manager.go,
  //you need to have AWS credentials configured as environment variables or configure
  //AWS EC2 IAM roles to allow access (see https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).
  //If you have other ways of authentication, you can provide your own session with
  //GetSecretsWithSession or WriteSecretsToENVWithSession.

  //This will grab values from secrets whose key names are in the format of an
  //environment variable name (all uppercase characters with digits and connected
  //with underscores. No dash or other symbols) and set them to corresponding
  //environment variables.
  err := config.WriteSecretsToENV("dev/application/for_testing")

  //So if "dev/application/for_testing" has {"GOOD_ENV1":"good","bad-env2":"bad"},
  //the above method will set only GOOD_ENV1="good" and not the other one

  //***************************** Metrics ***********************************

  //These enable metrics to know how to initialize the Statsd client
  factory := func() *metrics.Statsd {
    return metrics.NewStatsd("address", "prefix.", "version", "app", 1.0)
  }
  metrics.SetFactory(factory)

  //****************************** Server ***********************************

  svr := server.Server{
    Engine:               gin.New(),
    AppInfo:              &health.AppInfo{...},     //You should fill in this info
    ProjectInfo:          &health.ProjectInfo{...}, //You should fill in this info
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

  ahd1 := health.AdditionalHealthData{
    DependencyChecks: []HealthChecker{dbCheck, serviceCheck1},
    DataProvider:    func(c *gin.Context) map[string]interface{}{
      return map[string]interface{}{
        "custom": "data",
      }
    },
  }

  adh2 := health.AdditionalHealthData{
    DependencyChecks: []HealthChecker{dbCheck, serviceCheck1, serviceCheck2},
  }

  //Register 3 versions of the detailed health. Note that they are different as
  //"/v1" has additional custom data but does not have `serviceCheck2` dependency check,
  //and "/v3" has no dependency or additional data.
  svr.RegisterDetailedHealth("/v1", "v1 of app detailed health", ahd1)
  svr.RegisterDetailedHealth("/v2", "v2 of app detailed health", ahd2)
  svr.RegisterDetailedHealth("/v3", "v3 without custom data or dependency check", nil)
  svr.RegisterSimpleHealth()

  someHandler := func(c *gin.Context) {
    //Emitting statsd metric
    metrics.Increment("interesting.metric")

    logging.RL(c).Info("This is logging with correlation ID")

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

### Environment loading:
* Note: Application should set the prefered default environment.
* Note: It supports with/without Json configuration file (not other format). Developer needs to define the environment bindings in the Configure structure.
```
  // Implement two functions for Configuration struct (where defined environment bindings):
    func (c Configuration) IsSslEnabled() bool {return c.SslEnabled == "true"}
    func (c Configuration) GetSslSecretName() string {return c.SslSecretName}
  //
    main(){
      ...
      conf := &Configuration{}
        // Application should set the prefered default environment.
      appEnv, err := appenv.NewAppEnv(provider)
        // It overwrites existed environment
      err = appEnv.LoadEnv()
        // Load config file to conf
      err = config.LoadJsonConfigFile(configFile, conf)
        //Overwrite conf from environment parameters
      config.PopulateEnvConfig(conf)
      if conf.SslEnabled == "true" {
        err = appEnv.LoadSslCertificate()
      }
      ...
      //Application make decision on running on TLS or not.
      if conf.SslEnabled == "true" {
        err = svr.Engine.RunTLS(conf.BindAddress, appenv.SslCertFile, appenv.SslKeyFile)
      } else {
        err = svr.Engine.Run(conf.BindAddress)
      }
   }
```

* Environments:
```
  CLOUD_PROVIDER: AWS:   Pulls environments from secret manager, overwrites existed environment.
                  LOCAL: Use existed environment (default).
  SSL_ENABLED: true or false (files: server.crt & server.key)
               When CLOUD_PROVIDER set to AWS, downloads certificats from secret (SSL_SECRET_NAME) and
               save to the working directory.
               When set to LOCAL, uses ceritificate and key file under the application work dir.

  // AWS specific environment defined in task definition (sample):
  AWS_REGION: us-east-1 (default)
  // Secret for application enviornment/configurations:
  AWSSM_NAME: {dev/application/app_name}
  // Secret for SSL certificates
  SSL_SECRET_NAME: dev/application/appcerts
```
