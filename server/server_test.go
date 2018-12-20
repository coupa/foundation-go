package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coupa/foundation-go/health"
	"github.com/coupa/foundation-go/middleware"
	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var _ = Describe("Server", func() {
	handler := func(c *gin.Context) {
		c.JSON(http.StatusNoContent, "")
	}

	Describe("Server", func() {
		It("makes a server with middleware and simple health", func() {
			svr := Server{
				Engine:               gin.New(),
				AppInfo:              &health.AppInfo{},
				ProjectInfo:          &health.ProjectInfo{},
				AdditionalHealthData: map[string]*health.AdditionalHealthData{},
			}
			svr.UseMiddleware(middleware.Correlation())

			svr.Engine.GET("/test", handler)

			svr.RegisterSimpleHealth()

			req, _ := http.NewRequest("GET", "/test", nil)
			resp := httptest.NewRecorder()
			svr.Engine.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusNoContent))
			Expect(resp.Header().Get("X-Correlation-Id")).NotTo(BeEmpty())

			req, _ = http.NewRequest("GET", "/health", nil)
			resp = httptest.NewRecorder()
			svr.Engine.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Header().Get("X-Correlation-Id")).To(BeEmpty())
		})

		It("makes a server with middleware and detailed health", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := health.Health{
					"status":   health.WARN,
					"version":  "fakeVer",
					"revision": "fakeRev",
				}
				d, _ := json.Marshal(h)
				fmt.Fprint(w, string(d))
			}))
			defer ts.Close()

			svr := Server{
				Engine:               gin.New(),
				AppInfo:              &health.AppInfo{},
				ProjectInfo:          &health.ProjectInfo{},
				AdditionalHealthData: map[string]*health.AdditionalHealthData{},
			}
			svr.UseMiddleware(middleware.Correlation())

			dbCheck := health.SQLCheck{
				Name: "mysql",
				Type: "internal",
			}

			serviceCheck1 := health.WebCheck{
				Name: "some web",
				Type: "service",
				URL:  ts.URL,
			}

			custom := health.AdditionalHealthData{
				DependencyChecks: []health.HealthChecker{dbCheck, serviceCheck1},
			}

			svr.RegisterSimpleHealth()
			svr.RegisterDetailedHealth("/v1", "This is v1 detailed health", &custom)

			svr.Engine.GET("/test", handler)

			req, _ := http.NewRequest("GET", "/test", nil)
			resp := httptest.NewRecorder()
			svr.Engine.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusNoContent))
			Expect(resp.Header().Get("X-Correlation-Id")).NotTo(BeEmpty())

			req, _ = http.NewRequest("GET", "/health", nil)
			resp = httptest.NewRecorder()
			svr.Engine.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Header().Get("X-Correlation-Id")).To(BeEmpty())

			req, _ = http.NewRequest("GET", "/v1/health/detailed", nil)
			resp = httptest.NewRecorder()
			svr.Engine.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Header().Get("X-Correlation-Id")).NotTo(BeEmpty())
			d, _ := ioutil.ReadAll(resp.Body)
			h := health.Health{}
			json.Unmarshal(d, &h)
			Expect(len(h["dependencies"].([]interface{}))).To(Equal(2))
		})

		Describe("Timeout", func() {
			It("times out after HealthTimeout time if some checks have not finished", func() {
				svr := Server{
					Engine:               gin.New(),
					AppInfo:              &health.AppInfo{},
					ProjectInfo:          &health.ProjectInfo{},
					AdditionalHealthData: map[string]*health.AdditionalHealthData{},
				}

				dbCheck := health.SQLCheck{
					Name: "mysql",
					Type: "internal",
				}

				serviceCheck1 := ServiceCheck{}
				custom := health.AdditionalHealthData{
					DependencyChecks: []health.HealthChecker{dbCheck, serviceCheck1},
				}

				HealthTimeout = 0 * time.Second
				svr.RegisterSimpleHealth()
				svr.RegisterDetailedHealth("/v1", "This is v1 detailed health", &custom)

				svr.Engine.GET("/test", handler)

				req, _ := http.NewRequest("GET", "/test", nil)
				resp := httptest.NewRecorder()
				svr.Engine.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusNoContent))

				req, _ = http.NewRequest("GET", "/health", nil)
				resp = httptest.NewRecorder()
				svr.Engine.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusOK))

				req, _ = http.NewRequest("GET", "/v1/health/detailed", nil)
				resp = httptest.NewRecorder()
				svr.Engine.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusOK))
				d, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(string(d))
			})
		})

		Describe("RegisterDetailedHealth", func() {
			It("sets the description in the detailed health", func() {
				svr := Server{
					Engine:               gin.New(),
					AppInfo:              &health.AppInfo{},
					ProjectInfo:          &health.ProjectInfo{},
					AdditionalHealthData: map[string]*health.AdditionalHealthData{},
				}

				svr.RegisterDetailedHealth("/v1", "hello", nil)
				svr.RegisterDetailedHealth("/v2", "pizza", &health.AdditionalHealthData{
					Description: "This should be overwritten",
				})

				req, _ := http.NewRequest("GET", "/v1/health/detailed", nil)
				resp := httptest.NewRecorder()
				svr.Engine.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusOK))

				d, _ := ioutil.ReadAll(resp.Body)
				var h health.Health
				json.Unmarshal(d, &h)
				Expect(h["description"].(string)).To(Equal("hello"))

				req, _ = http.NewRequest("GET", "/v2/health/detailed", nil)
				resp = httptest.NewRecorder()
				svr.Engine.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusOK))

				d, _ = ioutil.ReadAll(resp.Body)
				var h2 health.Health
				json.Unmarshal(d, &h2)
				Expect(h2["description"].(string)).To(Equal("pizza"))
			})
		})
	})

	Describe("extractVersionKey", func() {
		It("", func() {
			//Versioned paths return the version
			Expect(extractVersionKey("/v1/some")).To(Equal("/v1"))
			Expect(extractVersionKey("/v10/some")).To(Equal("/v10"))

			Expect(extractVersionKey("/v1")).To(Equal("/v1"))
			Expect(extractVersionKey("v1")).To(Equal("/v1"))
			Expect(extractVersionKey("v1/some")).To(Equal("/v1"))

			//Non-versioned path returns a slash
			Expect(extractVersionKey("health")).To(Equal("/"))
			Expect(extractVersionKey("/health")).To(Equal("/"))
			Expect(extractVersionKey("/health/some")).To(Equal("/"))
			Expect(extractVersionKey("")).To(Equal("/"))
		})
	})
})

type ServiceCheck struct {
}

func (sc ServiceCheck) Check() *health.DependencyInfo {
	time.Sleep(1 * time.Second)
	return &health.DependencyInfo{}
}

func (sc ServiceCheck) GetName() string {
	return "fake"
}

func (sc ServiceCheck) GetType() string {
	return "service"
}
