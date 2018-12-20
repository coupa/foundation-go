package config_test

import (
	. "github.com/coupa/foundation-go/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("SplitByCommaSpace", func() {
		It("can add to an empty Health struct", func() {
			s := "a,b c ,  d, e"
			Expect(SplitByCommaSpace(s)).To(ConsistOf([]string{"a", "b", "c", "d", "e"}))

			s = "http://abc.com"
			Expect(SplitByCommaSpace(s)).To(ConsistOf([]string{"http://abc.com"}))

			s = "http://abc.com,https://bcd.com"
			Expect(SplitByCommaSpace(s)).To(ConsistOf([]string{"http://abc.com", "https://bcd.com"}))

			s = "http://abc.com,https://bcd.com http://cde.com"
			Expect(SplitByCommaSpace(s)).To(ConsistOf([]string{"http://abc.com", "https://bcd.com", "http://cde.com"}))
		})
	})
})
