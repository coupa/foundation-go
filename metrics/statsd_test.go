package metrics

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMetrics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metrics Suite")
}

var _ = Describe("Metrics", func() {
	Describe("flatten", func() {
		It("flattens string hashes into a slice", func() {
			t1 := map[string]string{
				"t1": "v1",
				"t2": "v2",
			}
			t2 := map[string]string{
				"t3": "v3",
			}
			r := flatten(t1, t2)
			Expect(r).To(ConsistOf([]string{"t1", "v1", "t2", "v2", "t3", "v3"}))

			r = flatten(t1)
			Expect(r).To(ConsistOf([]string{"t1", "v1", "t2", "v2"}))

			r = flatten()
			Expect(r).To(BeEmpty())
		})
	})

	Describe("addNameTag", func() {
		It("adds name as name tag to the tags", func() {
			r := addNameTag("name1")
			Expect(r).To(Equal([]map[string]string{
				map[string]string{"name": "name1"},
			}))

			r = addNameTag("name2", map[string]string{"a": "a"})
			Expect(r).To(ConsistOf([]map[string]string{
				map[string]string{"name": "name2"},
				map[string]string{"a": "a"},
			}))

			r = addNameTag("name3", map[string]string{"a": "a", "b": "b"})
			Expect(r).To(ConsistOf([]map[string]string{
				map[string]string{"name": "name3"},
				map[string]string{"a": "a", "b": "b"},
			}))
		})
	})
})
