/**
 * Maintains common health check logic
 */
package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	OK   = "OK"
	WARN = "WARN"
	CRIT = "CRIT"
)

/*
 * Health information struct.
 */
type HealthInfo struct {
	Status              string      `json:"status"`
	Version             string      `json:"version"`
	Revision            string      `json:"revision"`
	Uptime              int         `json:"uptime"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	Host                string      `json:"host"`
	Project             ProjectInfo `json:"project"`
	DBDependencies      []DBDependency
	ServiceDependencies []DependencyInfo
}

/*
 * Project information struct.
 */
type ProjectInfo struct {
	Repo   string   `json:"repo"`
	Home   string   `json:"home"`
	Owners []string `json:"owners"`
	Logs   []string `json:"logs"`
	Stats  []string `json:"stats"`
}

/*
 * Dependency information struct.
 */
type DependencyInfo struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Version      string    `json:"version"`
	Revision     string    `json:"revision"`
	State        StateInfo `json:"state"`
	ResponseTime float64   `json:"responseTime"`
}

/*
 * State information struct.
 */
type StateInfo struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

// Global variables
var (
	serverStartTime   time.Time  // Used to calculate server uptime
	HealthInformation HealthInfo // Globally shared HelathInformation instance.
)

type HealthCheckHandler struct {
	version  string
	revision string
}

type DetailedHealthCheckHandler struct {
	dBDependencies                []DBDependency
	httpEnpointHealthCheckService HTTPEndPointHealthCheckService
}

/*
 * Factory method for health check handler
 */
func NewHealthCheckHandler(version string, revision string) HealthCheckHandler {
	h := HealthCheckHandler{version, revision}
	return h
}

/*
 * Factory method for detailed health check handler
 */
func NewDetailedHealthCheckHandler(dBDependencies []DBDependency, serviceDependencies []ServiceDependencyInfo) DetailedHealthCheckHandler {
	httpEnpointHealthCheckService := NewHTTPEndPointHealthCheckService(serviceDependencies)
	h := DetailedHealthCheckHandler{dBDependencies, httpEnpointHealthCheckService}
	return h
}

/*
 * This health check endpoint can be plugged directly as a gin handler.
 */
func (handler HealthCheckHandler) HealthCheckHandler(gc *gin.Context) {
	gc.JSON(http.StatusOK, gin.H{"status": "OK", "version": handler.version, "revision": handler.revision})
}

/*
 * Detailed health check function
 */
func (handler DetailedHealthCheckHandler) DetailedHealthCheckHandler(gc *gin.Context) {
	var healthInfo HealthInfo
	for i := range handler.dBDependencies {
		dbStatusCheck(&handler.dBDependencies[i])
	}
	healthInfo.DBDependencies = handler.dBDependencies
	dependencyInfo := handler.httpEnpointHealthCheckService.CheckHttpServiceStatus()
	healthInfo.ServiceDependencies = dependencyInfo
	gc.JSON(http.StatusOK, healthInfo)
}
