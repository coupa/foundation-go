package persistence

type QueryParams struct {
	Operands []QueryExpression
	Limit    uint64
	Offset   uint64
	Order    []OrderStatement
}

type QueryOperatorType string

const (
	QUERY_OPERATOR_EQ          QueryOperatorType = "equals"
	QUERY_OPERATOR_NEQ                           = "not_equals"
	QUERY_OPERATOR_CONTAINS                      = "contains"
	QUERY_OPERATOR_IN                            = "in"
	QUERY_OPERATOR_GT                            = "gt"
	QUERY_OPERATOR_GTE                           = "gte"
	QUERY_OPERATOR_LT                            = "lt"
	QUERY_OPERATOR_LTE                           = "lte"
	QUERY_OPERATOR_STARTS_WITH                   = "starts_with"
	QUERY_OPERATOR_ENDS_WITH                     = "ends_with"
)

type QueryExpression struct {
	Key      string
	Operator QueryOperatorType
	Value    string
}

type OrderStatement struct {
	Key       string
	Direction OrderDirection
}

type OrderDirection string

const (
	ORDER_DIRECTION_ASC  OrderDirection = "asc"
	ORDER_DIRECTION_DESC                = "desc"
)
