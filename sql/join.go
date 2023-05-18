package sql

import "github.com/kovey/db-go/v2/sql/meta"

type Join struct {
	table   string
	alias   string
	columns []string
	on      string
}

func NewJoin(table, alias, on string) *Join {
	if alias == "" {
		alias = table
	}
	return &Join{table: table, alias: alias, columns: make([]string, 0), on: on}
}

func (j *Join) Columns(columns ...*meta.Column) *Join {
	for _, column := range columns {
		column.SetTable(j.alias)
		j.columns = append(j.columns, column.String())
	}

	return j
}

func (j *Join) CaseWhen(caseWhens ...*meta.CaseWhen) *Join {
	for _, caseWhen := range caseWhens {
		j.columns = append(j.columns, caseWhen.String())
	}

	return j
}
