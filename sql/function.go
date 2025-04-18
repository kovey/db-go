package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type Function struct {
	*base
	definer       string
	ifNotExists   bool
	function      string
	args          []*argType
	returns       string
	comment       string
	language      string
	sqlType       ksql.ProcedureSqlType
	sqlSecurity   ksql.SqlSecurity
	deterministic string
	sql           ksql.ExpressInterface
	isCreate      bool
}

func keywordFunction(builder *strings.Builder) {
	builder.WriteString(" FUNCTION")
}

func NewFunction() *Function {
	f := &Function{base: newBase(), isCreate: true}
	f.opChain.Append(f._keyword, f._definer, keywordFunction, f._function, f._args, f._returns, f._comment, f._language, f._deterministic, f._sqlType, f._sqlSecurity, f._body)
	return f
}

func (f *Function) _keyword(builder *strings.Builder) {
	if f.isCreate {
		keywordCreate(builder)
		return
	}

	keywordAlter(builder)
}

func (f *Function) _definer(builder *strings.Builder) {
	if !f.isCreate || f.definer == "" {
		return
	}

	builder.WriteString(" DEFINER = ")
	builder.WriteString(f.definer)
}

func (f *Function) _function(builder *strings.Builder) {
	if f.isCreate && f.ifNotExists {
		builder.WriteString(" IF NOT EXISTS")
	}
	builder.WriteString(" ")
	Backtick(f.function, builder)
}

func (f *Function) _args(builder *strings.Builder) {
	if !f.isCreate {
		return
	}

	builder.WriteString("(")
	for index, arg := range f.args {
		if index > 0 {
			builder.WriteString(", ")
		}

		Backtick(arg.name, builder)
		builder.WriteString(" ")
		builder.WriteString(arg.typ)
	}

	builder.WriteString(")")
}

func (f *Function) _returns(builder *strings.Builder) {
	if !f.isCreate {
		return
	}

	builder.WriteString(" RETURNS ")
	builder.WriteString(f.returns)
}

func (f *Function) _comment(builder *strings.Builder) {
	if f.comment == "" {
		return
	}

	builder.WriteString(" COMMENT ")
	Quote(f.comment, builder)
}

func (f *Function) _language(builder *strings.Builder) {
	if f.language == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(f.language)
}

func (f *Function) _deterministic(builder *strings.Builder) {
	if f.deterministic == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(f.deterministic)
}

func (f *Function) _sqlType(builder *strings.Builder) {
	if f.sqlType == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(string(f.sqlType))
}

func (f *Function) _sqlSecurity(builder *strings.Builder) {
	if f.sqlSecurity == "" {
		return
	}

	builder.WriteString(" SQL SECURITY ")
	builder.WriteString(string(f.sqlSecurity))
}

func (f *Function) _body(builder *strings.Builder) {
	if f.isCreate && f.sql != nil {
		builder.WriteString(" BEGIN ")
		builder.WriteString(DefaultEngine().formatOriginalRaw(f.sql))
		builder.WriteString(";")
		builder.WriteString(" END")
	}
}

func (f *Function) Definer(name string) ksql.FunctionInterface {
	if !f.isCreate {
		return f
	}

	f.definer = name
	return f
}

func (f *Function) IfNotExists() ksql.FunctionInterface {
	if !f.isCreate {
		return f
	}

	f.ifNotExists = true
	return f
}

func (f *Function) Function(name string) ksql.FunctionInterface {
	f.function = name
	return f
}

func (f *Function) Param(name, typ string) ksql.FunctionInterface {
	if !f.isCreate {
		return f
	}

	f.args = append(f.args, &argType{name: name, typ: typ})
	return f
}

func (f *Function) Returns(typ string) ksql.FunctionInterface {
	if !f.isCreate {
		return f
	}

	f.returns = typ
	return f
}

func (f *Function) Comment(comment string) ksql.FunctionInterface {
	f.comment = comment
	return f
}

func (f *Function) Language() ksql.FunctionInterface {
	f.language = "LANGUAGE SQL"
	return f
}

func (f *Function) Deterministic() ksql.FunctionInterface {
	f.deterministic = "DETERMINISTIC"
	return f
}

func (f *Function) DeterministicNot() ksql.FunctionInterface {
	f.deterministic = "NOT DETERMINISTIC"
	return f
}

func (f *Function) SqlType(sqlType ksql.ProcedureSqlType) ksql.FunctionInterface {
	f.sqlType = sqlType
	return f
}

func (f *Function) SqlSecurity(security ksql.SqlSecurity) ksql.FunctionInterface {
	f.sqlSecurity = security
	return f
}

func (f *Function) RoutineBody(raw ksql.ExpressInterface) ksql.FunctionInterface {
	if !f.isCreate {
		return f
	}

	f.sql = raw
	return f
}

func (f *Function) Alter() ksql.FunctionInterface {
	f.isCreate = false
	return f
}
