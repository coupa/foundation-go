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
	logging.InitLogging(projectName, appName, appVersion)
	logging.EnableJsonFormat()

	metrics.InitMetrics(projectName, appName, appVersion)

	log.WithFields(log.Fields{
		"project":    projectName,
		"app":        appName,
		"appVersion": appVersion,
	}).Info("MetricsMonitoring initialized")

}
