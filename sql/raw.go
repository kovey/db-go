package sql

import (
	"fmt"
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/express"
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
	default:
		return fmt.Sprintf("%d", tmp)
	}
}

func Quote(name string, builder *strings.Builder) {
	builder.WriteString("'")
	builder.WriteString(name)
	builder.WriteString("'")
}

func Backtick(name string, builder *strings.Builder) {
	if strings.HasPrefix(name, "(") {
		builder.WriteString(name)
		return
	}

	builder.WriteString("`")
	builder.WriteString(name)
	builder.WriteString("`")
}

func Column(column string, builder *strings.Builder) {
	if strings.HasPrefix(column, "(") {
		builder.WriteString(column)
		return
	}

	var table = ""
	if strings.Contains(column, ".") {
		info := strings.Split(column, ".")
		table = info[0]
		column = info[1]
	}

	if table != "" {
		Backtick(table, builder)
		builder.WriteString(".")
	}

	Backtick(column, builder)
}

func _formatSharding(name string, sharding ksql.Sharding) string {
	if strings.HasPrefix(name, "(") {
		return name
	}

	return ksql.FormatSharding(name, sharding)
}
