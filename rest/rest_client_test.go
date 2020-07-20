package rest_test

import (
	"github.com/coupa/foundation-go/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RestClient", func() {

	It("QueryParams to url translation", func() {
		restService.FindMany(persistence.QueryParams{
			Order: []persistence.OrderStatement{
				{Key: "a", Direction: persistence.ORDER_DIRECTION_DESC},
				{Key: "b", Direction: persistence.ORDER_DIRECTION_ASC},
			},
			Limit:  123,
			Offset: 456,
			Operands: []persistence.QueryExpression{
				{Key: "a", Operator: persistence.QUERY_OPERATOR_EQ, Value: "1"},
				{Key: "b", Operator: persistence.QUERY_OPERATOR_NEQ, Value: "2"},
				{Key: "c", Operator: persistence.QUERY_OPERATOR_CONTAINS, Value: "3"},
				{Key: "d", Operator: persistence.QUERY_OPERATOR_IN, Value: "4"},
				{Key: "e", Operator: persistence.QUERY_OPERATOR_GT, Value: "5"},
				{Key: "f", Operator: persistence.QUERY_OPERATOR_GTE, Value: "6"},
				{Key: "g", Operator: persistence.QUERY_OPERATOR_LT, Value: "7"},
				{Key: "h", Operator: persistence.QUERY_OPERATOR_LTE, Value: "8"},
				{Key: "i", Operator: persistence.QUERY_OPERATOR_STARTS_WITH, Value: "9"},
				{Key: "j", Operator: persistence.QUERY_OPERATOR_ENDS_WITH, Value: "10"},
			},
		})
		Expect(testServer.lastestRawQuery).To(ContainSubstring("order=a,desc"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("order=b,asc"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("limit=123"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("offset=456"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("q[a_equals]=1"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("q[b_not_equals]=2"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("q[c_contains]=3"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("q[d_in]=4"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("q[e_gt]=5"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("q[f_gte]=6"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("q[g_lt]=7"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("q[h_lte]=8"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("q[i_starts_with]=9"))
		Expect(testServer.lastestRawQuery).To(ContainSubstring("q[j_ends_with]=10"))
	})
})
