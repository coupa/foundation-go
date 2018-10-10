/**
 * Maintains common logging data and logic
 */
package logging

import (
	log "github.com/sirupsen/logrus"
)

var LoggingAppVersion = ""
var LoggingProject = ""
var LoggingApp = ""

// ----------------------------------------------------------------------------
func InitLogging(projectName string, appName string, appVersion string) {
	LoggingProject = projectName
	LoggingApp = appName
	LoggingAppVersion = appVersion
}

func EnableJsonFormat() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
}

