package sql

import (
	"fmt"
	"strings"

	"github.com/kovey/db-go/v2/sql/meta"
)

type Join struct {
	table   string
	alias   string
	columns []string
	on      string
	sub     *Select
}

func NewJoin(table, alias, on string, columns ...string) *Join {
	if alias == emptyStr {
		if strings.Contains(table, dot) {
			alias = strings.ReplaceAll(table, dot, underline)
		} else {
			alias = table
		}
	}
	j := &Join{table: table, alias: alias, on: on, columns: make([]string, len(columns))}
	j.init(columns)
	return j
}

func NewJoinSub(sub *Select, alias, on string, columns ...string) *Join {
	j := &Join{sub: sub, alias: alias, on: on, columns: make([]string, len(columns))}
	j.init(columns)
	return j
}

func (j *Join) tableName() string {
	if j.sub == nil {
		return formatValue(j.table)
	}

	return fmt.Sprintf(subFormat, j.sub.Prepare())
}

func (j *Join) init(columns []string) {
	for index, column := range columns {
		col := meta.NewColumn(column)
		col.SetTable(j.alias)
		j.columns[index] = col.String()
	}
}

func (j *Join) Columns(columns ...string) *Join {
	for _, column := range columns {
		col := meta.NewColumn(column)
		col.SetTable(j.alias)
		j.columns = append(j.columns, col.String())
	}

	return j
}

func (j *Join) ColMeta(columns ...*meta.Column) *Join {
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

func (j *Join) args() []any {
	if j.sub == nil {
		return nil
	}

	return j.sub.Args()
}
