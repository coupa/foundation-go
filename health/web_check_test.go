package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

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

			d := WebCheck{Name: "test", URL: ts.URL, Type: TypeService}.Check()
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

				d := WebCheck{Name: "test", URL: ts.URL, Type: TypeService}.Check()
				Expect(d.Name).To(Equal("test"))
				Expect(d.Version).To(Equal(""))
				Expect(d.Revision).To(Equal(""))
				Expect(d.State.Status).To(Equal(""))
			})
		})

		It("sets WARN for redirect response", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := Health{
					"status":   WARN,
					"version":  "fakeVer",
					"revision": "fakeRev",
				}
				w.WriteHeader(399)
				d, _ := json.Marshal(h)
				fmt.Fprint(w, string(d))
			}))
			defer ts.Close()
			d := WebCheck{Name: "test", URL: ts.URL, Type: TypeService}.Check()
			Expect(d.Name).To(Equal("test"))
			Expect(d.Version).To(Equal(""))
			Expect(d.Revision).To(Equal(""))
			Expect(d.State.Status).To(Equal(WARN))
			Expect(d.State.Details).To(HaveSuffix("redirected"))
		})

		It("sets CRIT for error response", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := Health{
					"status":   WARN,
					"version":  "fakeVer",
					"revision": "fakeRev",
				}
				w.WriteHeader(404)
				d, _ := json.Marshal(h)
				fmt.Fprint(w, string(d))
			}))
			defer ts.Close()

			d := WebCheck{Name: "test", URL: ts.URL, Type: TypeService}.Check()
			Expect(d.Name).To(Equal("test"))
			Expect(d.Version).To(Equal(""))
			Expect(d.Revision).To(Equal(""))
			Expect(d.State.Status).To(Equal(CRIT))
			Expect(d.State.Details).To(ContainSubstring("error checking"))
		})

		Context("service returns unparsable data", func() {
			It("returns error parsing text as details and WARN as status", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, "")
				}))
				defer ts.Close()

				d := WebCheck{Name: "test", URL: ts.URL, Type: TypeService}.Check()
				Expect(d.Name).To(Equal("test"))
				Expect(d.State.Status).To(Equal(WARN))
				Expect(d.State.Details).To(HavePrefix("Response body is not a key-value JSON object: "))
			})

			Context("with third-party type", func() {
				It("returns error parsing text as details and OK as status", func() {
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, "")
					}))
					defer ts.Close()

					d := WebCheck{Name: "test", URL: ts.URL, Type: TypeThirdParty}.Check()
					Expect(d.Name).To(Equal("test"))
					Expect(d.State.Status).To(Equal(OK))
					Expect(d.State.Details).To(HavePrefix("Response body is not a key-value JSON object: "))
				})
			})
		})

		Context("with ExpectedStatusCode set", func() {
			Context("and the status code matches", func() {
				It("sets OK regardless of the status in the response body", func() {
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						h := Health{
							"status":   WARN,
							"version":  "fakeVer",
							"revision": "fakeRev",
						}
						w.WriteHeader(200)
						d, _ := json.Marshal(h)
						fmt.Fprint(w, string(d))
					}))
					defer ts.Close()

					d := WebCheck{Name: "test", URL: ts.URL, Type: TypeService, ExpectedStatusCode: 200}.Check()
					Expect(d.Name).To(Equal("test"))
					Expect(d.Version).To(Equal("fakeVer"))
					Expect(d.Revision).To(Equal("fakeRev"))
					Expect(d.State.Status).To(Equal(OK))
				})

				It("sets OK as long as the status code matches", func() {
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						h := Health{
							"status":   WARN,
							"version":  "fakeVer",
							"revision": "fakeRev",
						}
						w.WriteHeader(400)
						d, _ := json.Marshal(h)
						fmt.Fprint(w, string(d))
					}))
					defer ts.Close()

					d := WebCheck{Name: "test", URL: ts.URL, Type: TypeService, ExpectedStatusCode: 400}.Check()
					Expect(d.Name).To(Equal("test"))
					Expect(d.Version).To(Equal(""))
					Expect(d.Revision).To(Equal(""))
					Expect(d.State.Status).To(Equal(OK))
				})

				Context("service returns unparsable data", func() {
					It("returns error parsing text as details and OK as status", func() {
						ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(300)
							fmt.Fprint(w, "")
						}))
						defer ts.Close()

						d := WebCheck{Name: "test", URL: ts.URL, Type: TypeService, ExpectedStatusCode: 300}.Check()
						Expect(d.Name).To(Equal("test"))
						Expect(d.State.Status).To(Equal(OK))
						Expect(d.State.Details).To(HaveSuffix("redirected"))

						d = WebCheck{Name: "test", URL: ts.URL, Type: TypeThirdParty, ExpectedStatusCode: 300}.Check()
						Expect(d.Name).To(Equal("test"))
						Expect(d.State.Status).To(Equal(OK))
						Expect(d.State.Details).To(HaveSuffix("redirected"))
					})
				})
			})

			Context("and the status code does not match", func() {
				It("sets CRIT regardless of the status in the response body", func() {
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						h := Health{
							"status":   WARN,
							"version":  "fakeVer",
							"revision": "fakeRev",
						}
						w.WriteHeader(200)
						d, _ := json.Marshal(h)
						fmt.Fprint(w, string(d))
					}))
					defer ts.Close()

					d := WebCheck{Name: "test", URL: ts.URL, Type: TypeService, ExpectedStatusCode: 201}.Check()
					Expect(d.Name).To(Equal("test"))
					Expect(d.Version).To(Equal("fakeVer"))
					Expect(d.Revision).To(Equal("fakeRev"))
					Expect(d.State.Status).To(Equal(CRIT))
				})

				Context("service returns unparsable data", func() {
					It("returns error parsing text as details and CRIT as status", func() {
						ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(300)
							fmt.Fprint(w, "")
						}))
						defer ts.Close()

						d := WebCheck{Name: "test", URL: ts.URL, Type: TypeService, ExpectedStatusCode: 301}.Check()
						Expect(d.Name).To(Equal("test"))
						Expect(d.State.Status).To(Equal(CRIT))
						Expect(d.State.Details).To(HavePrefix("Expected status code"))
						Expect(d.State.Details).To(HaveSuffix("redirected"))

						d = WebCheck{Name: "test", URL: ts.URL, Type: TypeThirdParty, ExpectedStatusCode: 301}.Check()
						Expect(d.Name).To(Equal("test"))
						Expect(d.State.Status).To(Equal(CRIT))
						Expect(d.State.Details).To(HavePrefix("Expected status code"))
						Expect(d.State.Details).To(HaveSuffix("redirected"))
					})
				})
			})
		})

		Context("when error", func() {
			It("returns error connection as details", func() {
				d := WebCheck{Name: "test", URL: "abc", Type: TypeService}.Check()
				Expect(d.Name).To(Equal("test"))
				Expect(d.State.Status).To(Equal(CRIT))
				Expect(d.State.Details).To(HavePrefix("Error connecting to `abc`: "))

				d = WebCheck{Name: "test", URL: "abc", Type: TypeService, ExpectedStatusCode: 200}.Check()
				Expect(d.Name).To(Equal("test"))
				Expect(d.State.Status).To(Equal(CRIT))
				Expect(d.State.Details).To(HavePrefix("Error connecting to `abc`: "))
			})
		})
	})
})
