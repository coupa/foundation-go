package persistence_test

import (
	"fmt"
	"github.com/coupa/foundation-go/persistence"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var persistenceServiceMysql *persistence.PersistenceServiceMySql

//var tx *sql.Tx

var _ = Describe("PersistenceServiceMySql", func() {
	JustBeforeEach(func() {
		var err error
		if persistenceServiceMysql == nil {
			persistenceServiceMysql, err = persistence.NewPersistenceServiceMySql("root:@/foundation_go_test", "test_entities", reflect.TypeOf(TestEntity{}))
			if err != nil {
				panic(err)
			}
		}
		_, err = persistenceServiceMysql.GetConnection().NewSession(nil).Exec("delete from test_entities")
		if err != nil {
			panic(err)
		}
	})

	//JustAfterEach(func() {
	//	err := tx.Rollback();
	//	if err != nil {
	//		panic(err)
	//	}
	//	tx = nil
	//})

	Context("#CreateOne #FindOne", func() {
		It("saves", func() {
			obj1 := TestEntity{StringA: "a", IntA: 1}
			err := persistenceServiceMysql.CreateOne(&obj1)
			Expect(err).To(BeNil())
			Expect(obj1.ID > 0).To(BeTrue())

			obj2, foundIt, err := persistenceServiceMysql.FindOne(fmt.Sprint(obj1.ID))
			Expect(err).To(BeNil())
			Expect(foundIt).To(BeTrue())
			Expect(obj2).To(Equal(obj1))
		})
	})

	Context("FindMany", func() {
		obj1 := TestEntity{StringA: "abc", IntA: 1}
		obj2 := TestEntity{StringA: "def", IntA: 1}

		JustBeforeEach(func() {
			persistenceServiceMysql.CreateOne(&obj1)
			persistenceServiceMysql.CreateOne(&obj2)
		})

		It("finds all", func() {
			values, err := persistenceServiceMysql.FindMany(persistence.QueryParams{})
			Expect(err).To(BeNil())
			Expect(len(values.([]TestEntity))).To(Equal(2))
		})
		It("limits", func() {
			values, err := persistenceServiceMysql.FindMany(persistence.QueryParams{Limit: 1})
			Expect(err).To(BeNil())
			Expect(len(values.([]TestEntity))).To(Equal(1))
		})
		It("orders", func() {
			values, err := persistenceServiceMysql.FindMany(persistence.QueryParams{Order: []persistence.OrderStatement{{Key: "string_a", Direction: persistence.ORDER_DIRECTION_DESC}}})
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(2))
			Expect(entities[0].ID).To(Equal(obj2.ID))
			Expect(entities[0].StringA).To(Equal(obj2.StringA))
		})
		It("offsets", func() {
			values, err := persistenceServiceMysql.FindMany(persistence.QueryParams{Limit: 10000, Offset: 1, Order: []persistence.OrderStatement{{Key: "string_a", Direction: persistence.ORDER_DIRECTION_DESC}}})
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(1))
			Expect(entities[0].ID).To(Equal(obj1.ID))
			Expect(entities[0].StringA).To(Equal(obj1.StringA))
		})
		It("equals based on string", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_EQ, Value: "abc"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(1))
			Expect(entities[0].ID).To(Equal(obj1.ID))
			Expect(entities[0].StringA).To(Equal(obj1.StringA))
		})
		It("equals based on int", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_EQ, Value: "1"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(2))
		})
		It("not equals - partial match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_NEQ, Value: "abc"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(1))
			Expect(entities[0].ID).To(Equal(obj2.ID))
			Expect(entities[0].StringA).To(Equal(obj2.StringA))
		})
		It("not equals - all match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_NEQ, Value: "something that doesnt exist in db"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(2))
		})
		It("contains", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_CONTAINS, Value: "bc"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(1))
			Expect(entities[0].ID).To(Equal(obj1.ID))
			Expect(entities[0].StringA).To(Equal(obj1.StringA))
		})
		It("in - strings - all match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_IN, Value: "abc,def,ghi"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(2))
		})
		It("in - strings - partial match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_IN, Value: "abc,ghi"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(1))
			Expect(entities[0].ID).To(Equal(obj1.ID))
			Expect(entities[0].StringA).To(Equal(obj1.StringA))
		})
		It("in - ints - all match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_IN, Value: "1,2"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(2))
		})
		It("in - ints - no match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_IN, Value: "3"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(0))
		})
		It("gt - match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_GT, Value: "0"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(2))
		})
		It("gt - no match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_GT, Value: "1"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(0))
		})
		It("gte - match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_GTE, Value: "0"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(2))
		})
		It("gte - no match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_GTE, Value: "2"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(0))
		})
		It("lt - match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_LT, Value: "2"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(2))
		})
		It("lt - no match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_LT, Value: "1"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(0))
		})
		It("lte - match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_LTE, Value: "1"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(2))
		})
		It("lte - no match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "int_a", Operator: persistence.QUERY_OPERATOR_LTE, Value: "0"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(0))
		})
		It("starts_with - match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_STARTS_WITH, Value: "ab"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(1))
		})
		It("starts_with - no match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_STARTS_WITH, Value: "bc"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(0))
		})
		It("ends_with - match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_ENDS_WITH, Value: "c"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(1))
		})
		It("ends_with - no match", func() {
			params := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_ENDS_WITH, Value: "ab"},
			}}
			values, err := persistenceServiceMysql.FindMany(params)
			Expect(err).To(BeNil())
			entities := values.([]TestEntity)
			Expect(len(entities)).To(Equal(0))
		})
	})

	Context("UpdateOne", func() {
		It("updates", func() {
			obj1 := TestEntity{StringA: "a"}
			err := persistenceServiceMysql.CreateOne(&obj1)
			Expect(err).To(BeNil())
			Expect(obj1.ID > 0).To(BeTrue())
			paramsA := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_EQ, Value: "a"},
			}}
			values, err := persistenceServiceMysql.FindMany(paramsA)
			Expect(err).To(BeNil())
			Expect(len(values.([]TestEntity))).To(Equal(1))

			obj1.StringA = "b"
			foundIt, err := persistenceServiceMysql.UpdateOne(fmt.Sprint(obj1.ID), &obj1)
			Expect(err).To(BeNil())
			Expect(foundIt).To(BeTrue())

			values, err = persistenceServiceMysql.FindMany(paramsA)
			Expect(err).To(BeNil())
			Expect(len(values.([]TestEntity))).To(Equal(0))

			paramsB := persistence.QueryParams{Operands: []persistence.QueryExpression{
				{Key: "string_a", Operator: persistence.QUERY_OPERATOR_EQ, Value: "b"},
			}}
			values, err = persistenceServiceMysql.FindMany(paramsB)
			Expect(err).To(BeNil())
			Expect(len(values.([]TestEntity))).To(Equal(1))
		})
		It("doesnt update - id not found", func() {
			obj1 := TestEntity{}
			foundIt, err := persistenceServiceMysql.UpdateOne("-1", &obj1)
			Expect(err).To(BeNil())
			Expect(foundIt).To(BeFalse())
		})
		It("returns error if input obj is not a pointer", func() {
			_, err := persistenceServiceMysql.UpdateOne("-1", TestEntity{})
			Expect(err).To(Not(BeNil()))
			Expect(err.Error()).To(Equal("obj must be a pointer"))
		})
	})

	Context("DeleteOne", func() {
		It("deletes if exists", func() {
			obj1 := TestEntity{}
			err := persistenceServiceMysql.CreateOne(&obj1)
			Expect(err).To(BeNil())
			Expect(obj1.ID > 0).To(BeTrue())

			foundIt, err := persistenceServiceMysql.DeleteOne(fmt.Sprint(obj1.ID))
			Expect(err).To(BeNil())
			Expect(foundIt).To(BeTrue())

			_, foundIt, err = persistenceServiceMysql.FindOne(fmt.Sprint(obj1.ID))
			Expect(err).To(BeNil())
			Expect(foundIt).To(BeFalse())

		})
		It("returns foundIt=false if doesnt exists", func() {
			foundIt, err := persistenceServiceMysql.DeleteOne("-1")
			Expect(err).To(BeNil())
			Expect(foundIt).To(BeFalse())
		})
	})
})
