package health

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	r  *gin.Engine
	ts *httptest.Server
)

var _ = Describe("health", func() {
	Describe("trivialHealthCheck", func() {
		BeforeEach(func() {
			gin.SetMode(gin.TestMode)
			r = gin.New()
		})

		Context("get health check", func() {
			It("can get a health check", func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
				}))
				handler := NewHealthCheckHandler("y.yy", "x.xx")

				r.GET("/health", handler.HealthCheckHandler)
				defer ts.Close()
				request, _ := http.NewRequest("GET", "/health", nil)
				response := httptest.NewRecorder()
				r.ServeHTTP(response, request)
				fmt.Println("response code", response.Code)
				Expect(response.Code).To(Equal(http.StatusOK))
				var result HealthInfo
				json.NewDecoder(response.Body).Decode(&result)
				Expect(response.Body).ShouldNot(BeNil())
				Expect(result.Revision).To(Equal("x.xx"))
				Expect(result.Version).To(Equal("y.yy"))
			})
		})

		Context("get DB health check", func() {
			It("gives DB Status", func() {
				dbBasic := DependencyInfo{
					Name: "mysql",
				}
				dependencies := []DBDependency{{
					BasicInfo: dbBasic,
					Dialect:   "mysql",
					DSN:       "root@tcp(127.0.0.1:3306)/test-database?parseTime=true",
				}}
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
				}))
				var serviceDependencies []ServiceDependencyInfo
				serviceDependencies = append(serviceDependencies, NewServiceDependency("testService1", "testVersion1", "testRevision1", "http://testhost/health"))
				detailedHealthCheckHandler, _ := NewDetailedHealthCheckHandler(dependencies, serviceDependencies)
				r.GET("/detailed-health", detailedHealthCheckHandler.DetailedHealthCheckHandler)
				defer ts.Close()
				request, _ := http.NewRequest("GET", "/detailed-health", nil)
				response := httptest.NewRecorder()
				r.ServeHTTP(response, request)
				fmt.Println("response code", response)
				Expect(response.Code).To(Equal(http.StatusOK))
				var result HealthInfo
				json.NewDecoder(response.Body).Decode(&result)
				Expect(response.Body).ShouldNot(BeNil())
				Expect(result.DBDependencies[0].BasicInfo.State.Status).To(Equal("CRIT"))
				Expect(result.ServiceDependencies[0].State.Status).To(Equal("CRIT"))
				Expect(result.ServiceDependencies[0].Name).To(Equal("testService1"))
				Expect(result.ServiceDependencies[0].Version).To(Equal("testVersion1"))
				Expect(result.ServiceDependencies[0].Revision).To(Equal("testRevision1"))
			})
		})
	})

})

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Health check test Suite")
}
