package rest

import (
	"fmt"
	"github.com/coupa/foundation-go/persistence"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var activeAdminRegex = regexp.MustCompile(`q\[(.+?)_(.+?)\]`)

type HttpQueryParser interface {
	Parse(values url.Values) (persistence.QueryParams, error)
}

var _ HttpQueryParser = (*HttpQueryParserRailsActiveAdmin)(nil)

type HttpQueryParserRailsActiveAdmin struct{}

func (self *HttpQueryParserRailsActiveAdmin) Parse(query url.Values) (persistence.QueryParams, error) {
	ret := persistence.QueryParams{
		Operands: []persistence.QueryExpression{},
	}
	for key, values := range query {
		if len(values) == 0 {
			continue
		}
		if strings.EqualFold("limit", key) {
			limit, _ := strconv.Atoi(values[0])
			ret.Limit = uint64(limit)
		} else if strings.EqualFold("offset", key) {
			offset, _ := strconv.Atoi(values[0])
			ret.Offset = uint64(offset)
		} else if strings.EqualFold("order", key) {
			ret.Order = []persistence.OrderStatement{}
			for _, value := range values {
				splitValue := strings.Split(value+",", ",")
				direction := persistence.ORDER_DIRECTION_ASC
				if strings.EqualFold("desc", strings.ToLower(splitValue[1])) {
					direction = persistence.ORDER_DIRECTION_DESC
				}
				ret.Order = append(ret.Order, persistence.OrderStatement{Key: splitValue[0], Direction: direction})
			}
		} else {
			matched := activeAdminRegex.FindAllStringSubmatch(key, -1)
			if matched == nil {
				continue
			}
			if len(matched[0]) == 1 {
				queryExpression := persistence.QueryExpression{
					Key:      matched[0][1],
					Operator: persistence.QUERY_OPERATOR_EQ,
					Value:    values[0],
				}
				ret.Operands = append(ret.Operands, queryExpression)
			} else if len(matched[0]) == 3 {
				var operator persistence.QueryOperatorType
				switch matched[0][2] {
				case "equals":
					operator = persistence.QUERY_OPERATOR_EQ
				case "not_equals":
					operator = persistence.QUERY_OPERATOR_NEQ
				case "gt":
					operator = persistence.QUERY_OPERATOR_GT
				case "gte":
					operator = persistence.QUERY_OPERATOR_GTE
				case "lt":
					operator = persistence.QUERY_OPERATOR_LT
				case "lte":
					operator = persistence.QUERY_OPERATOR_LTE
				case "in":
					operator = persistence.QUERY_OPERATOR_IN
				case "contains":
					operator = persistence.QUERY_OPERATOR_CONTAINS
				case "starts_with":
					operator = persistence.QUERY_OPERATOR_STARTS_WITH
				case "ends_with":
					operator = persistence.QUERY_OPERATOR_ENDS_WITH
				default:
					return ret, fmt.Errorf("unknown operator '%v'", matched[0][2])
				}
				queryExpression := persistence.QueryExpression{
					Key:      matched[0][1],
					Operator: operator,
					Value:    values[0],
				}
				ret.Operands = append(ret.Operands, queryExpression)
			}

		}
	}
	return ret, nil
}
