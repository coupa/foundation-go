package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/coupa/foundation-go/health"
	"github.com/gin-gonic/gin"
)

var (
	HealthTimeout = 5 * time.Second
)

//Server is based on Gin https://github.com/gin-gonic/gin
type Server struct {
	Engine      *gin.Engine
	AppInfo     *health.AppInfo
	ProjectInfo *health.ProjectInfo
	//The key is like "/v1" or "/v2" with a leading slash
	AdditionalHealthData map[string]*health.AdditionalHealthData
}

func (s *Server) UseMiddleware(mw gin.HandlerFunc) {
	s.Engine.Use(mw)
}

//RegisterSimpleHealth registers /health
//Simple health is used by the load balancer's health checks and dependent services'
//detailed health.
func (s *Server) RegisterSimpleHealth() {
	s.Engine.GET("/health", s.simpleHealth)
}

//RegisterDetailedHealth registers detail health at /<versionGroup>/health/detailed
//versionGroup must be like "/v1", then the endpoint is "/v1/health/detailed"
//There should be a leading slash for the versionGroup.
//If versionGroup is empty string "", then this detailed health endpoint is not
//versioned, like "/health/detailed"
//A detailed health should only check for other service's simple health. Never
//check the detailed health of a depending service.
func (s *Server) RegisterDetailedHealth(versionGroup, description string, h *health.AdditionalHealthData) {
	if h == nil {
		h = new(health.AdditionalHealthData)
	}
	h.Description = description
	if s.AdditionalHealthData == nil {
		s.AdditionalHealthData = map[string]*health.AdditionalHealthData{}
	}
	s.AdditionalHealthData[versionGroup] = h
	s.Engine.GET(versionGroup+"/health/detailed", s.detailedHealth)
}

func (s *Server) simpleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, health.NewSimpleHealth(s.AppInfo, health.OK))
}

func (s *Server) detailedHealth(c *gin.Context) {
	ver := extractVersionKey(c.Request.URL.Path)
	ahd := new(health.AdditionalHealthData)
	if s.AdditionalHealthData != nil {
		ahd = s.AdditionalHealthData[ver]
	}
	h := health.NewDetailedHealth(s.AppInfo, s.ProjectInfo, ahd.Description)

	if ahd.DataProvider != nil {
		if extraData := ahd.DataProvider(c); extraData != nil {
			for k, v := range extraData {
				h[k] = v
			}
		}
	}

	checks := ahd.DependencyChecks
	if num := len(checks); num > 0 {
		buffer := make(chan *health.DependencyInfo, num)
		hChecks := make(chan health.HealthChecker, num)

		for _, hc := range checks {
			hChecks <- hc
			go func(healthCheck health.HealthChecker) {
				buffer <- healthCheck.Check()
			}(hc)
		}

		for i := 0; i < num; i++ {
			select {
			case info := <-buffer:
				<-hChecks
				h.AddDependency(*info)
			case <-time.After(HealthTimeout):
			}
		}
		close(hChecks)

		//If there is any item in hChecks, it is timed out
		for hc := range hChecks {
			h.AddDependency(health.DependencyInfo{
				Name:         hc.GetName(),
				Type:         hc.GetType(),
				ResponseTime: HealthTimeout.Seconds(),
				State: health.DependencyState{
					Status:  health.CRIT,
					Details: fmt.Sprintf("Health check timed out after %f seconds", HealthTimeout.Seconds()),
				},
			})
		}
	}
	c.JSON(http.StatusOK, h)
}

func extractVersionKey(path string) string {
	regex, err := regexp.Compile(`^/*v\d+`)
	if err != nil {
		return ""
	}
	v := regex.FindString(path)
	if !strings.HasPrefix(v, "/") {
		v = "/" + v
	}
	return v
}
