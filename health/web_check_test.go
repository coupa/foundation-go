package health_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/coupa/foundation-go/health"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WebService", func() {
	Describe("Check", func() {
		It("returns the data from the service health check", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := Health{
					"status":   WARN,
					"version":  "fakeVer",
					"revision": "fakeRev",
				}
				d, _ := json.Marshal(h)
				fmt.Fprint(w, string(d))
			}))
			defer ts.Close()

			d := WebCheck{Name: "test", URL: ts.URL}.Check()
			Expect(d.Name).To(Equal("test"))
			Expect(d.Version).To(Equal("fakeVer"))
			Expect(d.Revision).To(Equal("fakeRev"))
			Expect(d.State.Status).To(Equal(WARN))
		})

		Context("service returns empty data", func() {
			It("returns empty fields", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					h := Health{}
					d, _ := json.Marshal(h)
					fmt.Fprint(w, string(d))
				}))
				defer ts.Close()

				d := WebCheck{Name: "test", URL: ts.URL}.Check()
				Expect(d.Name).To(Equal("test"))
				Expect(d.Version).To(Equal(""))
				Expect(d.Revision).To(Equal(""))
				Expect(d.State.Status).To(Equal(""))
			})
		})

		Context("service returns unparsable data", func() {
			It("returns error parsing text as details", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "")
				}))
				defer ts.Close()

				d := WebCheck{Name: "test", URL: ts.URL}.Check()
				Expect(d.Name).To(Equal("test"))
				Expect(d.State.Status).To(Equal(WARN))
				Expect(d.State.Details).To(HavePrefix("Error parsing response body: "))
			})
		})

		Context("when error", func() {
			It("returns error connection as details", func() {
				d := WebCheck{Name: "test", URL: "abc"}.Check()
				Expect(d.Name).To(Equal("test"))
				Expect(d.State.Status).To(Equal(CRIT))
				Expect(d.State.Details).To(HavePrefix("Error connecting to `abc`: "))
			})
		})
	})
})
