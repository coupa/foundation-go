package config

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWSSecretsManager", func() {
	if os.Getenv("TEST_SECRETS_MANAGER") == "true" {
		Describe("GetSecrets", func() {
			It("can add to an empty Health struct", func() {
				m, _, err := GetSecrets("dev/application/for_testing")
				if err != nil {
					panic("To test AWS Secrets Manager, you need to provide AWS credentials that have access to the `dev/application/for_testing` secrets")
				}
				Expect(m).NotTo(BeNil())
				Expect(m["TEST_DESCRIPTION"]).To(Equal("For testing"))
				Expect(m["not_set_to_env"]).To(Equal("Will not be set to your ENV variables"))
			})
		})

		Describe("WriteSecretsToENV", func() {
			It("sets secrets to ENV variables", func() {
				err := WriteSecretsToENV("dev/application/for_testing")
				if err != nil {
					panic("To test AWS Secrets Manager, you need to provide AWS credentials that have access to the `dev/application/for_testing` secrets")
				}
				Expect(os.Getenv("TEST_DESCRIPTION")).To(Equal("For testing"))
				Expect(os.Getenv("not_set_to_env")).To(BeEmpty())
			})
		})
	}

	Describe("regexENVName", func() {
		It("matches only standard environment variable names", func() {
			tests := []string{"A", "AB", "A1", "A_B", "A_1", "A_B_1"}
			for _, s := range tests {
				Expect(regexENVName.MatchString(s)).To(BeTrue())
			}
			tests = []string{"1", "_", "a", "1A", "_A", "aA", "A_", "AB_", "Aa", "A_a", "A__", "A-A"}
			for _, s := range tests {
				Expect(regexENVName.MatchString(s)).To(BeFalse())
			}
		})
	})
})
