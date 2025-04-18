package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type View struct {
	*base
	view        string
	algorithm   ksql.ViewAlg
	definer     string
	sqlSecurity ksql.SqlSecurity
	columns     []string
	as          ksql.QueryInterface
	with        string
	check       string
	replace     string
	isCreate    bool
}

func NewView() *View {
	v := &View{base: newBase(), isCreate: true}
	v.opChain.Append(v._keyword, v._replace, v._alg, v._name, v._as)
	return v
}

func (v *View) _keyword(builder *strings.Builder) {
	if v.isCreate {
		builder.WriteString("CREATE")
		return
	}

	builder.WriteString("ALTER")
}

func (v *View) _replace(builder *strings.Builder) {
	if !v.isCreate {
		return
	}

	operator.BuildPureString(v.replace, builder)
}

func (v *View) _alg(builder *strings.Builder) {
	if v.algorithm != "" {
		builder.WriteString(" ALGORITHM = ")
		builder.WriteString(string(v.algorithm))
	}

	if v.definer != "" {
		builder.WriteString(" DEFINER = ")
		builder.WriteString(v.definer)
	}

	if v.sqlSecurity != "" {
		builder.WriteString(" SQL SECURITY ")
		builder.WriteString(string(v.sqlSecurity))
	}
}

func (v *View) _name(builder *strings.Builder) {
	builder.WriteString(" VIEW")
	operator.BuildColumnString(v.view, builder)
	if len(v.columns) == 0 {
		return
	}

	builder.WriteString(" (")
	for index, column := range v.columns {
		if index > 0 {
			builder.WriteString(", ")
		}

		operator.Column(column, builder)
	}
	builder.WriteString(")")
}

func (v *View) _as(builder *strings.Builder) {
	if v.as != nil {
		builder.WriteString(" AS ")
		builder.WriteString(v.as.Prepare())
		v.binds = append(v.binds, v.as.Binds()...)
	}

	operator.BuildPureString(v.with, builder)
	operator.BuildPureString(v.check, builder)
}

func (v *View) Replace() ksql.ViewInterface {
	if !v.isCreate {
		return v
	}

	v.replace = "OR REPLACE"
	return v
}

func (v *View) Algorithm(alg ksql.ViewAlg) ksql.ViewInterface {
	v.algorithm = alg
	return v
}

func (v *View) Definer(definer string) ksql.ViewInterface {
	v.definer = definer
	return v
}

func (v *View) SqlSecurity(security ksql.SqlSecurity) ksql.ViewInterface {
	v.sqlSecurity = security
	return v
}

func (v *View) View(name string) ksql.ViewInterface {
	v.view = name
	return v
}

func (v *View) Columns(columns ...string) ksql.ViewInterface {
	v.columns = append(v.columns, columns...)
	return v
}

func (v *View) As(query ksql.QueryInterface) ksql.ViewInterface {
	v.as = query
	return v
}

func (v *View) WithCascaded() ksql.ViewInterface {
	v.with = "WITH CASCADED"
	return v
}

func (v *View) WithLocal() ksql.ViewInterface {
	v.with = "WITH LOCAL"
	return v
}

func (v *View) CheckOption() ksql.ViewInterface {
	v.check = "CHECK OPTION"
	return v
}

func (v *View) Alter() ksql.ViewInterface {
	v.isCreate = false
	v.replace = ""
	return v
}
