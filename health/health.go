package health

import (
	"time"

	"github.com/coupa/foundation-go/config"
	"github.com/gin-gonic/gin"
)

var (
	startTime time.Time
)

const (
	OK   = "OK"
	WARN = "WARN"
	CRIT = "CRIT"
)

type AppInfo struct {
	Version  string
	Revision string

	AppName  string `env:"APPLICATION_NAME"`
	Hostname string `env:"HOSTNAME"`
}

//FillFromENV will load pre-dfined values for env tags from environment variables
func (ai *AppInfo) FillFromENV(version, revision string) *AppInfo {
	config.PopulateEnvConfig(ai)
	return ai
}

type ProjectInfo struct {
	Repo   string   `json:"repo"   env:"PROJECT_REPO"`
	Home   string   `json:"home"   env:"PROJECT_HOME"`
	Owners []string `json:"owners"`
	Logs   []string `json:"logs"`
	Stats  []string `json:"stats"`

	OwnersStr string `json:"-" env:"PROJECT_OWNERS"`
	LogsStr   string `json:"-" env:"PROJECT_LOG_URLS"`
	StatsStr  string `json:"-" env:"PROJECT_STATS_URLS"`
}

//FillFromENV will load pre-dfined values for env tags from environment variables
func (pi *ProjectInfo) FillFromENV() *ProjectInfo {
	config.PopulateEnvConfig(pi)
	pi.convertStrToSlices()
	return pi
}

//Coverts comma-separated string fields to respective string slices
func (pi *ProjectInfo) convertStrToSlices() {
	pi.Owners = config.SplitByCommaSpace(pi.OwnersStr)
	pi.Logs = config.SplitByCommaSpace(pi.LogsStr)
	pi.Stats = config.SplitByCommaSpace(pi.StatsStr)
}

type DependencyInfo struct {
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	Version      string          `json:"version"`
	Revision     string          `json:"revision"`
	State        DependencyState `json:"state"`
	ResponseTime float64         `json:"responseTime"`
}

type DependencyState struct {
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
}

func init() {
	startTime = time.Now()
}

//UpTime is the application up time.
func UpTime() int64 {
	return int64(time.Since(startTime).Seconds())
}

/**
  // To make data for detailed health check:
  h := NewDetailedHealth(c, "some...")
  d := DependencyInfo{...}
  h.AddDependency(d) // Add as many dependencies as needed.
*/
type Health map[string]interface{}

func (h Health) AddDependency(d *DependencyInfo) {
	if d == nil {
		return
	}
	if h["dependencies"] == nil {
		h["dependencies"] = []interface{}{*d}
		return
	}
	h["dependencies"] = append(h["dependencies"].([]interface{}), *d)
}

func (h Health) SetDependencies(d []DependencyInfo) {
	h["dependencies"] = d
}

//NewSimpleHealth creates a health struct that can be rendered for the simple health check.
func NewSimpleHealth(ai *AppInfo, status string) Health {
	if ai == nil {
		ai = &AppInfo{}
	}
	return Health{
		"status":   status,
		"version":  ai.Version,
		"revision": ai.Revision,
	}
}

//NewDetailedHealth creates a health struct without any dependency.
func NewDetailedHealth(ai *AppInfo, pi *ProjectInfo, description string) Health {
	if ai == nil {
		ai = &AppInfo{}
	}
	if pi == nil {
		pi = &ProjectInfo{}
	}
	health := NewSimpleHealth(ai, OK)
	health["project"] = *pi
	health["host"] = ai.Hostname
	health["description"] = description
	health["name"] = ai.AppName
	health["uptime"] = UpTime()
	return health
}

type HealthChecker interface {
	Check() *DependencyInfo
	GetName() string
	GetType() string
}

type AdditionalHealthData struct {
	DependencyChecks []HealthChecker
	//The custom data can override the default health detail values
	DataProvider func(*gin.Context) map[string]interface{}

	//Description is set in server.RegisterDetailedHealth function, so there is no need
	//to set this field when initializing AdditionalHealthData struct
	Description string
}
