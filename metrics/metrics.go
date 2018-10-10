/**
 * Exposes APIs for consumers to use standardized StatsD metrics.
 * Maintains and injects common application data with stats
 */
package metrics

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
	"os"
)

type Tags struct {
	Name                 string
	EnterpriseInstanceId string
	ErrorCode            string
}

var version = ""
var project = ""
var app = ""
var component = ""

// ----------------------------------------------------------------------------
func InitMetrics(projectName string, appName string, appVersion string, componentName string) {
	project = projectName
	app = appName
	component = componentName
	version = appVersion
}

func getStatsUrl() string {
	var url = "172.17.0.1:8125"

	// For local testing or else, statsD url can be overridden using env variable
	if os.Getenv("STATSD_URL") != "" {
		url = os.Getenv("STATSD_URL")
	}
	return url
}

func getStatsHost() string {
	var host = "app.io.coupahost.com"

	if os.Getenv("STATSD_HOST") != "" {
		host = os.Getenv("STATSD_HOST")
	}
	return host
}

func GetStatsD(t Tags) *statsd.Client {
	statsDClient, err := statsd.New(
		statsd.Address(getStatsUrl()),
		statsd.TagsFormat(statsd.InfluxDB),
		statsd.Prefix(getStatsHost()+"./."+project+"."),
		statsd.Tags("version", version, "app", app, "component", component,
			"enterprise_instance_id", t.EnterpriseInstanceId, "name", t.Name,
			"error_code", t.ErrorCode),
	)
	if err != nil {
		// If nothing is listening on the target port, an error is returned and
		// the returned client does nothing but is still usable. So we can
		// just log the error and go on.
		log.WithFields(log.Fields{
			"errorCode": "STATSD",
		}).Error(err.Error())
		return nil
	}
	return statsDClient
}

// ----------------------------------------------------------------------------
func StatsIncrement(key string) {
	StatsIncrementWithTags(key, Tags{})
}

func StatsIncrementWithTags(key string, tags Tags) {
	statsD := GetStatsD(tags)
	if statsD == nil {
		return
	}
	defer statsD.Close()

	statsD.Increment(key)
}

func StatsTime(callback func(), key string) {
	StatsTimeWithTags(callback, key, Tags{})
}

func StatsTimeWithTags(callback func(), key string, tags Tags) {
	statsD := GetStatsD(tags)
	if statsD == nil {
		return
	}

	defer statsD.Close()

	func() {
		defer statsD.NewTiming().Send(key)
		callback()
	}()
}
