package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/coupa/foundation-go/health"
	"github.com/coupa/foundation-go/logging"
	"github.com/coupa/foundation-go/server"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus/hooks/test"

	. "github.com/coupa/foundation-go/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestLogger", func() {
	handler := func(c *gin.Context) {
		c.JSON(http.StatusNoContent, "")
	}

	Describe("RequestLogger middleware", func() {
		It("logs the required fields for a request", func() {
			logging.InitStandardLogger("VeRsIoN")
			hook := test.NewGlobal()

			svr := server.Server{
				Engine:      gin.New(),
				AppInfo:     &health.AppInfo{},
				ProjectInfo: &health.ProjectInfo{},
			}
			svr.UseMiddleware(RequestLogger(false))

			svr.Engine.GET("/test", handler)

			req, _ := http.NewRequest("GET", "/test", nil)
			resp := httptest.NewRecorder()
			svr.Engine.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusNoContent))

			data := hook.LastEntry().Data
			Expect(data["version"]).To(Equal("VeRsIoN"))
			Expect(data["level"]).To(Equal("info"))
			Expect(data["path"]).To(Equal("/test"))
			Expect(data["status"]).To(Equal(http.StatusNoContent))
			Expect(data["method"]).To(Equal("GET"))
			Expect(data["parameters"]).To(BeEmpty())
			Expect(data["message"]).To(BeNil())
		})

		Context("with correlation ID", func() {
			It("logs the correlation ID", func() {
				logging.InitStandardLogger("VeRsIoN")
				hook := test.NewGlobal()

				svr := server.Server{
					Engine:      gin.New(),
					AppInfo:     &health.AppInfo{},
					ProjectInfo: &health.ProjectInfo{},
				}
				//Order of adding the middlewares should not matter
				svr.UseMiddleware(RequestLogger(false))
				svr.UseMiddleware(Correlation())

				svr.Engine.GET("/test", handler)

				req, _ := http.NewRequest("GET", "/test", nil)
				resp := httptest.NewRecorder()
				svr.Engine.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusNoContent))

				data := hook.LastEntry().Data
				Expect(data["correlation_id"]).NotTo(BeEmpty())
			})
		})
	})
})
