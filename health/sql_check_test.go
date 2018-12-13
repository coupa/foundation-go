package health_test

import (
	. "github.com/coupa/foundation-go/health"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SQL", func() {
	Describe("Check", func() {
		It("shows unknown type of db connection", func() {
			c := SQLCheck{Name: "test"}
			d := c.Check()
			Expect(d.Name).To(Equal("test"))
			Expect(d.State.Status).To(Equal(CRIT))
			Expect(d.State.Details).To(Equal("Unknown type of DB connection"))
		})
	})
})
