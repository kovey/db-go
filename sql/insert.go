package sql

import (
	ksql "github.com/kovey/db-go/v3"
)

type Insert struct {
	*base
	columns []string
	from    ksql.QueryInterface
}

func NewInsert() *Insert {
	i := &Insert{base: &base{hasPrepared: false}}
	i.keyword("INSERT INTO ")
	return i
}

func (i *Insert) Table(table string) ksql.InsertInterface {
	Column(table, &i.builder)
	i.builder.WriteString(" (")
	return i
}

func (i *Insert) Add(column string, data any) ksql.InsertInterface {
	i.columns = append(i.columns, column)
	i.binds = append(i.binds, data)
	return i
}

func (i *Insert) Columns(columns ...string) ksql.InsertInterface {
	i.columns = columns
	return i
}

func (i *Insert) buildFrom() string {
	i.builder.WriteString(") ")
	i.builder.WriteString(i.from.Prepare())
	i.binds = append(i.binds, i.from.Binds()...)
	return i.base.Prepare()
}

func (i *Insert) buildNormal() string {
	i.builder.WriteString(") VALUES (")
	for j := 0; j < len(i.binds); j++ {
		if j > 0 {
			i.builder.WriteString(",")
		}
		i.builder.WriteString("?")
	}
	i.builder.WriteString(")")
	return i.base.Prepare()
}

func (i *Insert) Prepare() string {
	if i.hasPrepared {
		return i.base.Prepare()
	}

	i.hasPrepared = true
	for index, column := range i.columns {
		if index > 0 {
			i.builder.WriteString(",")
		}

		Column(column, &i.builder)
	}

	if i.from != nil {
		return i.buildFrom()
	}

	return i.buildNormal()
}

func (i *Insert) From(query ksql.QueryInterface) ksql.InsertInterface {
	i.from = query
	return i
}
