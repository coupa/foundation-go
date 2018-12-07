/**
 * Maintains common health check logic
 */
package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthInfo struct {
	Status       string           `json:"status"`
	Version      string           `json:"version"`
	Revision     string           `json:"revision"`
	Uptime       int              `json:"uptime"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Host         string           `json:"host"`
	Project      ProjectInfo      `json:"project"`
	Dependencies []DependencyInfo `json:"dependencies"`
}

type ProjectInfo struct {
	Repo   string   `json:"repo"`
	Home   string   `json:"home"`
	Owners []string `json:"owners"`
	Logs   []string `json:"logs"`
	Stats  []string `json:"stats"`
}

type DependencyInfo struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Version      string    `json:"version"`
	Revision     string    `json:"revision"`
	State        StateInfo `json:"state"`
	ResponseTime float64   `json:"responseTime"`
}

type StateInfo struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

var (
	serverStartTime   time.Time // Used to calculate server uptime
	HealthInformation HealthInfo
)

func HealthCheckHandler(gc *gin.Context) {
	healthInfo := HealthInformation
	gc.JSON(http.StatusOK, gin.H{"status": "OK", "version": healthInfo.Version, "revision": healthInfo.Revision})
}
