package health_test

import (
	. "github.com/coupa/foundation-go/health"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health", func() {
	Describe("AddDependency", func() {
		It("can add to an empty Health struct", func() {
			h := Health{}
			h.AddDependency(DependencyInfo{Name: "test"})
			ds := h["dependencies"].([]interface{})
			Expect(ds).To(HaveLen(1))
			d := ds[0].(DependencyInfo)
			Expect(d.Name).To(Equal("test"))

			h.AddDependency(DependencyInfo{Name: "test2"})
			ds = h["dependencies"].([]interface{})
			Expect(ds).To(HaveLen(2))
			d = ds[1].(DependencyInfo)
			Expect(d.Name).To(Equal("test2"))
		})
	})
})
