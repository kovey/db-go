package operator

import "strings"

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

func BuildPureString(value string, builder *strings.Builder) {
	if value == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(value)
}

func BuildQuoteString(value string, builder *strings.Builder) {
	if value == "" {
		return
	}

	builder.WriteString(" ")
	Quote(value, builder)
}

func BuildBacktickString(value string, builder *strings.Builder) {
	if value == "" {
		return
	}

	builder.WriteString(" ")
	Backtick(value, builder)
}

func BuildColumnString(value string, builder *strings.Builder) {
	if value == "" {
		return
	}

	builder.WriteString(" ")
	Column(value, builder)
}
