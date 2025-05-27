package criteria

import (
	"fmt"
	"strings"
)

type Operator string

var (
	GT Operator = ">"
	LT Operator = "<"
	EQ Operator = "="
)

type Condition interface {
	ToSQL() string
}

type Criteria struct {
	Condition Condition
}

type SimpleCondition struct {
	Field    string
	Operator Operator
	Value    any
}

func (sc *SimpleCondition) ToSQL() string {
	return fmt.Sprintf("WHERE %s %s %v", sc.Field, sc.Operator, sc.Value)
}

func (cr Criteria) Build(query string, startFrom uint64, valueOperator string) string {
	sql := cr.Condition.ToSQL()

	fields := strings.Fields(sql)

	for i, f := range fields {
		if f == "?" {
			fields[i] = fmt.Sprintf(valueOperator, startFrom)
			startFrom++
		}
	}

	return fmt.Sprintf("%s\n%s", query, strings.Join(fields, " "))
}