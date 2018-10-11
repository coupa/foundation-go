/**
 * Entry point for initializing common modules and or common utility functions
 */
package foundation

import (
	log "github.com/sirupsen/logrus"
	"github.com/coupa/foundation-go/logging"
	"github.com/coupa/foundation-go/metrics"
)

func InitMetricsMonitoring(projectName string, appName string, appVersion string) {
	InitMetricsMonitoringOnComponent(projectName, appName, appVersion, "")
}

func InitMetricsMonitoringOnComponent(projectName string, appName string, appVersion string, componentName string) {
	logging.InitLogging(projectName, appName, appVersion)
	logging.EnableJsonFormat()

	metrics.InitMetrics(projectName, appName, appVersion, componentName)

	log.WithFields(log.Fields{
		"project":    projectName,
		"app":        appName,
		"component":  componentName,
		"appVersion": appVersion,
	}).Info("MetricsMonitoring initialized")

}
