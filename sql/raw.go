package sql

import (
	"fmt"
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/express"
	"github.com/kovey/db-go/v3/sql/operator"
)

func Raw(raw string, binds ...any) ksql.ExpressInterface {
	raw = strings.Trim(raw, " \n\r\t\v")
	return express.NewStatement(raw, binds)
}

func RawValue(val any) string {
	switch tmp := val.(type) {
	case string:
		return fmt.Sprintf("'%s'", tmp)
	case float32, float64:
		return fmt.Sprintf("%f", tmp)
	case bool:
		return fmt.Sprintf("%t", tmp)
	default:
		if t, ok := val.(StringInterface); ok {
			return fmt.Sprintf("'%s'", t.String())
		}

		return fmt.Sprintf("%d", tmp)
	}
}

func Quote(name string, builder *strings.Builder) {
	operator.Quote(name, builder)
}

func Backtick(name string, builder *strings.Builder) {
	operator.Backtick(name, builder)
}

func Column(column string, builder *strings.Builder) {
	operator.Column(column, builder)
}

func _formatSharding(name string, sharding ksql.Sharding) string {
	if strings.HasPrefix(name, "(") {
		return name
	}

	return ksql.FormatSharding(name, sharding)
}
