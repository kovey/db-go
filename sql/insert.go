package sql

import (
	"github.com/kovey/db-go/v3"
)

type Insert struct {
	*base
	columns []string
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
