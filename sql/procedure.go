package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

func keywordProcedure(builder *strings.Builder) {
	builder.WriteString(" PROCEDURE")
}

type argType struct {
	name string
	typ  string
}

type Procedure struct {
	*base
	definer       string
	ifNotExists   bool
	procedure     string
	ins           []*argType
	outs          []*argType
	inouts        []*argType
	comment       string
	language      string
	sqlType       ksql.ProcedureSqlType
	sqlSecurity   ksql.SqlSecurity
	deterministic string
	sql           ksql.ExpressInterface
	isCreate      bool
}

func NewProcedure() *Procedure {
	p := &Procedure{base: newBase(), isCreate: true}
	p.opChain.Append(p._keyword, p._definer, keywordProcedure, p._procedure, p._args, p._comment, p._language, p._deterministic, p._sqlType, p._sqlSecurity, p._body)
	return p
}

func (p *Procedure) _keyword(builder *strings.Builder) {
	if p.isCreate {
		keywordCreate(builder)
		return
	}

	keywordAlter(builder)
}

func (p *Procedure) _definer(builder *strings.Builder) {
	if !p.isCreate || p.definer == "" {
		return
	}

	builder.WriteString(" DEFINER = ")
	builder.WriteString(p.definer)
}

func (p *Procedure) _procedure(builder *strings.Builder) {
	if p.isCreate && p.ifNotExists {
		builder.WriteString(" IF NOT EXISTS")
	}
	builder.WriteString(" ")
	Backtick(p.procedure, builder)
}

func (p *Procedure) _arg(arg *argType, builder *strings.Builder, index *int, typ string) {
	if *index > 0 {
		builder.WriteString(", ")
	}

	builder.WriteString(typ)
	builder.WriteString(" ")
	Backtick(arg.name, builder)
	builder.WriteString(" ")
	builder.WriteString(arg.typ)
	*index++
}

func (p *Procedure) _args(builder *strings.Builder) {
	if !p.isCreate {
		return
	}

	builder.WriteString("(")
	if len(p.ins) == 0 && len(p.outs) == 0 && len(p.inouts) == 0 {
		builder.WriteString(")")
		return
	}

	index := 0
	for _, arg := range p.ins {
		p._arg(arg, builder, &index, "IN")
	}

	for _, arg := range p.outs {
		p._arg(arg, builder, &index, "OUT")
	}

	for _, arg := range p.inouts {
		p._arg(arg, builder, &index, "INOUT")
	}

	builder.WriteString(")")
}

func (p *Procedure) _comment(builder *strings.Builder) {
	if p.comment == "" {
		return
	}

	builder.WriteString(" COMMENT ")
	Quote(p.comment, builder)
}

func (p *Procedure) _language(builder *strings.Builder) {
	if p.language == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(p.language)
}

func (p *Procedure) _deterministic(builder *strings.Builder) {
	if !p.isCreate || p.deterministic == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(p.deterministic)
}

func (p *Procedure) _sqlType(builder *strings.Builder) {
	if p.sqlType == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(string(p.sqlType))
}

func (p *Procedure) _sqlSecurity(builder *strings.Builder) {
	if p.sqlSecurity == "" {
		return
	}

	builder.WriteString(" SQL SECURITY ")
	builder.WriteString(string(p.sqlSecurity))
}

func (p *Procedure) _body(builder *strings.Builder) {
	if p.isCreate && p.sql != nil {
		builder.WriteString(" BEGIN ")
		builder.WriteString(DefaultEngine().formatOriginalRaw(p.sql))
		builder.WriteString(";")
		builder.WriteString(" END")
	}
}

func (p *Procedure) Definer(name string) ksql.ProcedureInterface {
	if !p.isCreate {
		return p
	}

	p.definer = name
	return p
}

func (p *Procedure) IfNotExists() ksql.ProcedureInterface {
	if !p.isCreate {
		return p
	}

	p.ifNotExists = true
	return p
}

func (p *Procedure) Procedure(name string) ksql.ProcedureInterface {
	p.procedure = name
	return p
}

func (p *Procedure) In(name, typ string) ksql.ProcedureInterface {
	if !p.isCreate {
		return p
	}

	p.ins = append(p.ins, &argType{name: name, typ: typ})
	return p
}

func (p *Procedure) Out(name, typ string) ksql.ProcedureInterface {
	if !p.isCreate {
		return p
	}

	p.outs = append(p.outs, &argType{name: name, typ: typ})
	return p
}

func (p *Procedure) InOut(name, typ string) ksql.ProcedureInterface {
	if !p.isCreate {
		return p
	}

	p.inouts = append(p.inouts, &argType{name: name, typ: typ})
	return p
}

func (p *Procedure) Comment(comment string) ksql.ProcedureInterface {
	p.comment = comment
	return p
}

func (p *Procedure) Language() ksql.ProcedureInterface {
	p.language = "LANGUAGE SQL"
	return p
}

func (p *Procedure) Deterministic() ksql.ProcedureInterface {
	if !p.isCreate {
		return p
	}

	p.deterministic = "DETERMINISTIC"
	return p
}

func (p *Procedure) DeterministicNot() ksql.ProcedureInterface {
	if !p.isCreate {
		return p
	}

	p.deterministic = "NOT DETERMINISTIC"
	return p
}

func (p *Procedure) SqlType(sqlType ksql.ProcedureSqlType) ksql.ProcedureInterface {
	p.sqlType = sqlType
	return p
}

func (p *Procedure) SqlSecurity(security ksql.SqlSecurity) ksql.ProcedureInterface {
	p.sqlSecurity = security
	return p
}

func (p *Procedure) RoutineBody(sql ksql.ExpressInterface) ksql.ProcedureInterface {
	if !p.isCreate || sql.IsExec() {
		return p
	}

	p.sql = sql
	return p
}

func (p *Procedure) Alter() ksql.ProcedureInterface {
	p.isCreate = false
	return p
}
