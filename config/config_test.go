package config_test

import (
	"os"

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

	Describe("PopulateEnvConfig", func() {
		BeforeEach(func() {
			os.Unsetenv("TEST_ENV_KEY")
			os.Unsetenv("TEST_ENV_BOOL_KEY")
			os.Unsetenv("TEST_ENV_INT_KEY")
		})
		type conf struct {
			Test     string `env:"TEST_ENV_KEY"`
			TestBool bool   `env:"TEST_ENV_BOOL_KEY"`
			TestInt  int    `env:"TEST_ENV_INT_KEY"`
		}
		It("populates env tags of a struct from environment variables", func() {
			os.Setenv("TEST_ENV_KEY", "some")
			c := conf{}
			PopulateEnvConfig(&c)
			Expect(c.Test).To(Equal("some"))
			Expect(c.TestBool).To(Equal(false))
			Expect(c.TestInt).To(Equal(0))

			os.Setenv("TEST_ENV_BOOL_KEY", "true")
			PopulateEnvConfig(&c)
			Expect(c.Test).To(Equal("some"))
			Expect(c.TestBool).To(Equal(true))
			Expect(c.TestInt).To(Equal(0))

			os.Setenv("TEST_ENV_INT_KEY", "100")
			PopulateEnvConfig(&c)
			Expect(c.Test).To(Equal("some"))
			Expect(c.TestBool).To(Equal(true))
			Expect(c.TestInt).To(Equal(100))
		})
	})
})
