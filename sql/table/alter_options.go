package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type AlterOptions struct {
	options []ksql.BuildInterface
}

func (a *AlterOptions) Empty() bool {
	return len(a.options) == 0
}

func (a *AlterOptions) Append(option ksql.BuildInterface) *AlterOptions {
	a.options = append(a.options, option)
	return a
}

func (a *AlterOptions) Build(builder *strings.Builder) {
	for index, option := range a.options {
		if index > 0 {
			builder.WriteString(", ")
		}

		option.Build(builder)
	}
}

type EditAlterOption struct {
	option ksql.BuildInterface
}

func NewEditAlterOption(option ksql.BuildInterface) *EditAlterOption {
	return &EditAlterOption{option: option}
}

func (a *EditAlterOption) Build(builder *strings.Builder) {
	builder.WriteString("ALTER ")
	a.option.Build(builder)
}

type AddAlterOption struct {
	option ksql.BuildInterface
}

func NewAddAlterOption(option ksql.BuildInterface) *AddAlterOption {
	return &AddAlterOption{option: option}
}

func (a *AddAlterOption) Build(builder *strings.Builder) {
	builder.WriteString("ADD ")
	a.option.Build(builder)
}

type DefaultAlterOption struct {
	eqs []*EqAlterOption
}

func (d *DefaultAlterOption) Empty() bool {
	return len(d.eqs) == 0
}

func (d *DefaultAlterOption) Option(key, value string) *DefaultAlterOption {
	d.eqs = append(d.eqs, NewEqAlterOption(key, value))
	return d
}

func (a *DefaultAlterOption) Build(builder *strings.Builder) {
	builder.WriteString("DEFAULT")
	for _, eq := range a.eqs {
		builder.WriteString(" ")
		eq.Build(builder)
	}
}

type EqAlterOption struct {
	key   string
	value string
}

func NewEqAlterOption(key, value string) *EqAlterOption {
	return &EqAlterOption{key: key, value: value}
}

func (a *EqAlterOption) Build(builder *strings.Builder) {
	builder.WriteString(a.key)
	builder.WriteString(" = ")
	builder.WriteString(a.value)
}

type AlterOption struct {
	method  string
	prefix  string
	value   string
	isStr   bool
	isField bool
}

func NewAlterOption(method, prefix, value string) *AlterOption {
	return &AlterOption{method: method, prefix: prefix, value: value}
}

func (a *AlterOption) IsStr() *AlterOption {
	a.isStr = true
	a.isField = false
	return a
}

func (a *AlterOption) IsField() *AlterOption {
	a.isStr = false
	a.isField = true
	return a
}

func (a *AlterOption) Build(builder *strings.Builder) {
	canAddSpace := false
	if a.method != "" {
		builder.WriteString(a.method)
		canAddSpace = true
	}

	if a.prefix != "" {
		if canAddSpace {
			builder.WriteString(" ")
		}

		builder.WriteString(a.prefix)
		canAddSpace = true
	}

	if a.value != "" {
		if canAddSpace {
			builder.WriteString(" ")
		}

		if a.isStr {
			operator.Quote(a.value, builder)
		} else if a.isField {
			operator.Backtick(a.value, builder)
		} else {
			builder.WriteString(a.value)
		}
	}
}
