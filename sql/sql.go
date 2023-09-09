package sql

import (
	"fmt"
	"strings"
)

const (
	valueStringFormat string = "'%s'"
	valueFormat       string = "`%s`"
)

type SqlInterface interface {
	Args() []any
	Prepare() string
	String() string
}

type WhereInterface interface {
	SqlInterface
	Eq(field string, value any)
	Neq(field string, value any)
	Like(field string, value any)
	Between(field string, from, to any)
	Gt(field string, value any)
	Ge(field string, value any)
	Lt(field string, value any)
	Le(field string, value any)
	In(field string, value []any)
	NotIn(field string, value []any)
	IsNull(field string)
	IsNotNull(field string)
	Statement(statement string)
	OrPrepare() string
}

func FormatValue(field string) string {
	return formatValue(field)
}

func formatString(value string) string {
	return fmt.Sprintf(valueStringFormat, value)
}

func formatValue(field string) string {
	info := strings.Split(field, ".")
	if len(info) != 2 {
		return fmt.Sprintf(valueFormat, field)
	}

	info[0] = fmt.Sprintf(valueFormat, info[0])
	info[1] = fmt.Sprintf(valueFormat, info[1])

	return strings.Join(info, ".")
}

func String(s SqlInterface) string {
	sql := s.Prepare()
	for _, arg := range s.Args() {
		switch a := arg.(type) {
		case string:
			sql = strings.Replace(sql, "?", formatString(a), 1)
		case float32, float64:
			sql = strings.Replace(sql, "?", formatString(fmt.Sprintf("%f", a)), 1)
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64:
			sql = strings.Replace(sql, "?", formatString(fmt.Sprintf("%d", a)), 1)
		default:
			sql = strings.Replace(sql, "?", formatString(fmt.Sprintf("%s", a)), 1)
		}
	}

	return sql
}

func formatOrder(order string) string {
	info := strings.Split(order, " ")
	if len(info) != 2 {
		return order
	}

	info[0] = formatValue(info[0])

	return strings.Join(info, " ")
}
