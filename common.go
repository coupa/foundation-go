/**
 * Entry point for initializing common modules and or common utility functions
 */
package common_go

import (
	log "github.com/sirupsen/logrus"
)

/**
 * Though individual middleware can be registered by themselves, provides
 * a utility to return them all
 */
func GetAllMiddlewareFunc() []func(c ...interface{}) {
	allFuncs := []func(c ...interface{}){
		RequestIdMiddlewareFunc,
		LogMiddlewareFunc,
		MetricsMiddlewareFunc,
	}

	return allFuncs
}

func InitMetricsMonitoring(projectName string, appName string, appVersion string, componentName string) {
	InitLogging(projectName, appName, appVersion)
	InitMetrics(projectName, appName, appVersion, componentName)

	log.WithFields(log.Fields{
		"project":    project,
		"app":        app,
		"component":  component,
		"appVersion": version,
	}).Info("MetricsMonitoring initialized")

}
