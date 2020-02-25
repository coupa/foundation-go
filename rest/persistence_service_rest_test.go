package rest_test

import (
	"fmt"
	"github.com/coupa/foundation-go/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CrudController and RestClient", func() {

	JustBeforeEach(func() {
		testServer.resetMockDb()
	})

	Context("Index", func() {
		It("200 when records exists", func() {
			obj1 := TestModel{StringA: "1"}
			testServer.mockDb.CreateOne(&obj1)
			Expect(obj1.ID).To(Not(Equal(0)))

			values, err := restService.FindMany(persistence.QueryParams{})
			Expect(err).To(BeNil())
			Expect(testServer.latestResponseStatusCode).To(Equal(200))
			Expect(values).To(Equal([]TestModel{obj1}))
		})
		It("200 when no records", func() {
			values, err := restService.FindMany(persistence.QueryParams{})
			Expect(err).To(BeNil())
			Expect(testServer.latestResponseStatusCode).To(Equal(200))
			Expect(values).To(Equal([]TestModel{}))
		})
	})

	Context("Create", func() {
		It("200", func() {
			obj1 := TestModel{StringA: "1", IntA: 2}
			err := restService.CreateOne(&obj1)
			Expect(err).To(BeNil())
			Expect(obj1.ID).To(Not(Equal(0)))
			_, foundIt := testServer.mockDb.db[obj1.ID]
			Expect(foundIt).To(BeTrue())
			Expect(testServer.latestResponseStatusCode).To(Equal(200))
		})
	})

	Context("Show", func() {
		It("200", func() {
			obj1 := TestModel{StringA: "1"}
			testServer.mockDb.CreateOne(&obj1)
			Expect(obj1.ID).To(Not(Equal(0)))
			obj2, foundIt, err := restService.FindOne(fmt.Sprint(obj1.ID))
			Expect(err).To(BeNil())
			Expect(testServer.latestResponseStatusCode).To(Equal(200))
			Expect(foundIt).To(BeTrue())
			Expect(obj2).To(Equal(obj1))
		})
		It("404", func() {
			_, foundIt, err := restService.FindOne("-1")
			Expect(err).To(BeNil())
			Expect(testServer.latestResponseStatusCode).To(Equal(404))
			Expect(foundIt).To(BeFalse())
		})
	})

	Context("Update", func() {
		It("200", func() {
			obj1 := TestModel{StringA: "1"}
			testServer.mockDb.CreateOne(&obj1)
			Expect(obj1.ID).To(Not(Equal(0)))
			obj1.StringA = "2"
			foundIt, err := restService.UpdateOne(fmt.Sprint(obj1.ID), &obj1)

			Expect(err).To(BeNil())
			Expect(testServer.latestResponseStatusCode).To(Equal(200))
			Expect(foundIt).To(BeTrue())
			Expect(testServer.mockDb.db[obj1.ID].StringA).To(Equal("2"))
		})
		It("404", func() {
			obj1 := TestModel{StringA: "1"}
			foundIt, err := restService.UpdateOne("-1", &obj1)
			Expect(err).To(BeNil())
			Expect(testServer.latestResponseStatusCode).To(Equal(404))
			Expect(foundIt).To(BeFalse())
		})
	})

	Context("Delete", func() {
		var obj1 TestModel

		JustBeforeEach(func() {
			obj1 = TestModel{StringA: "1"}
			testServer.mockDb.CreateOne(&obj1)
			Expect(obj1.ID).To(Not(Equal(0)))
		})

		It("204", func() {
			Expect(testServer.mockDb.db[obj1.ID]).To(Not(BeNil()))
			foundIt, err := restService.DeleteOne(fmt.Sprint(obj1.ID))

			Expect(err).To(BeNil())
			Expect(testServer.latestResponseStatusCode).To(Equal(204))
			Expect(foundIt).To(BeTrue())
			Expect(testServer.mockDb.db[obj1.ID]).To(BeNil())
		})
		It("404", func() {
			foundIt, err := restService.DeleteOne("-1")

			Expect(err).To(BeNil())
			Expect(testServer.latestResponseStatusCode).To(Equal(404))
			Expect(foundIt).To(BeFalse())
		})
	})

})
