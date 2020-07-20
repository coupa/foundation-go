package rest_test

import (
	"github.com/coupa/foundation-go/persistence"
	"github.com/coupa/foundation-go/rest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("QueryParser", func() {
	Context("HttpQueryParserRailsActiveAdmin", func() {
		It("#Parse", func() {
			httpQueryParser := &rest.HttpQueryParserRailsActiveAdmin{}
			values, err := url.ParseQuery(
				"&q[a_equals]=3" +
					"&q[b_not_equals]=4" +
					"&q[c_contains]=5" +
					"&q[d_in]=6" +
					"&q[e_gt]=7" +
					"&q[f_gte]=8" +
					"&q[g_lt]=9" +
					"&q[h_lte]=10" +
					"&q[i_starts_with]=11" +
					"&q[j_ends_with]=12" +

					"&limit=1" +
					"&offset=2" +
					"&order=a,desc" +
					"&order=b,asc" +
					"&order=c")

			qp, err := httpQueryParser.Parse(values)
			Expect(err).To(BeNil())

			Expect(len(qp.Operands)).To(Equal(10))
			qpm := map[string]persistence.QueryExpression{}
			for _, queryExpression := range qp.Operands {
				qpm[queryExpression.Key] = queryExpression
			}
			Expect(qpm["a"]).To(Equal(persistence.QueryExpression{Key: "a", Operator: persistence.QUERY_OPERATOR_EQ, Value: "3"}))
			Expect(qpm["b"]).To(Equal(persistence.QueryExpression{Key: "b", Operator: persistence.QUERY_OPERATOR_NEQ, Value: "4"}))
			Expect(qpm["c"]).To(Equal(persistence.QueryExpression{Key: "c", Operator: persistence.QUERY_OPERATOR_CONTAINS, Value: "5"}))
			Expect(qpm["d"]).To(Equal(persistence.QueryExpression{Key: "d", Operator: persistence.QUERY_OPERATOR_IN, Value: "6"}))
			Expect(qpm["e"]).To(Equal(persistence.QueryExpression{Key: "e", Operator: persistence.QUERY_OPERATOR_GT, Value: "7"}))
			Expect(qpm["f"]).To(Equal(persistence.QueryExpression{Key: "f", Operator: persistence.QUERY_OPERATOR_GTE, Value: "8"}))
			Expect(qpm["g"]).To(Equal(persistence.QueryExpression{Key: "g", Operator: persistence.QUERY_OPERATOR_LT, Value: "9"}))
			Expect(qpm["h"]).To(Equal(persistence.QueryExpression{Key: "h", Operator: persistence.QUERY_OPERATOR_LTE, Value: "10"}))
			Expect(qpm["i"]).To(Equal(persistence.QueryExpression{Key: "i", Operator: persistence.QUERY_OPERATOR_STARTS_WITH, Value: "11"}))
			Expect(qpm["j"]).To(Equal(persistence.QueryExpression{Key: "j", Operator: persistence.QUERY_OPERATOR_ENDS_WITH, Value: "12"}))
			Expect(qp.Limit).To(Equal(uint64(1)))
			Expect(qp.Offset).To(Equal(uint64(2)))
			Expect(len(qp.Order)).To(Equal(3))
			Expect(qp.Order[0]).To(Equal(persistence.OrderStatement{Key: "a", Direction: persistence.ORDER_DIRECTION_DESC}))
			Expect(qp.Order[1]).To(Equal(persistence.OrderStatement{Key: "b", Direction: persistence.ORDER_DIRECTION_ASC}))
			Expect(qp.Order[2]).To(Equal(persistence.OrderStatement{Key: "c", Direction: persistence.ORDER_DIRECTION_ASC}))
		})
	})
})
