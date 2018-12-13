package logging_test

import (
	"testing"

	"github.com/sirupsen/logrus/hooks/test"

	. "github.com/coupa/foundation-go/logging"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogging(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logging Suite")
}

var _ = Describe("Logging", func() {
	Describe("RequestLogger middleware", func() {
		It("logs the required fields for a request", func() {
			logger, hook := test.NewNullLogger()
			InitLogger("version1", logger)
			logger.Info("")

			data := hook.LastEntry().Data
			Expect(data["version"]).To(Equal("version1"))
			Expect(data["level"]).To(Equal("info"))
			Expect(data["message"]).To(BeNil())
		})

		Context("with log level above Error", func() {
			It("logs the required fields for a request", func() {
				logger, hook := test.NewNullLogger()
				InitLogger("version2", logger)
				logger.Error("")

				data := hook.LastEntry().Data
				Expect(data["version"]).To(Equal("version2"))
				Expect(data["level"]).To(Equal("error"))
				Expect(data["message"]).To(Equal(""))
			})
		})
	})
})
