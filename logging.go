/**
 * Maintains common logging data and logic
 */
package common_go

import (
	log "github.com/sirupsen/logrus"
)

var LoggingAppVersion = ""
var LoggingProject string = ""
var LoggingApp string = ""

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
}

// ----------------------------------------------------------------------------
func InitLogging(projectName string, appName string, appVersion string) {
	LoggingProject = projectName
	LoggingApp = appName
	LoggingAppVersion = appVersion
}
