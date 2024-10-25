package sql

import (
	"github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/table"
)

type Table struct {
	*base
	table   string
	engine  string
	charset string
	collate string
	comment string
	columns []*table.Column
	indexes []*table.Index
}

func NewTable() *Table {
	ta := &Table{base: &base{hasPrepared: false}, engine: "InnoDB", charset: "utf8mb4", collate: "utf8mb4_general_ci"}
	ta.keyword("CREATE TABLE ")
	return ta
}

func (ta *Table) AddColumn(column, t string, length, scale int, sets ...string) ksql.ColumnInterface {
	ct := table.ParseType(t, length, scale, sets...)
	if ct == nil {
		return nil
	}

	col := table.NewColumn(column, ct)
	ta.columns = append(ta.columns, col)
	return col
}

func (ta *Table) AddIndex(name string, t ksql.IndexType, column ...string) ksql.CreateTableInterface {
	index := &table.Index{Name: name, Type: t}
	index.Columns(column...)
	ta.indexes = append(ta.indexes, index)
	return ta
}

func (ta *Table) AddPrimary(column string) ksql.CreateTableInterface {
	return ta.AddIndex("", ksql.Index_Type_Primary, column)
}

func (ta *Table) Table(table string) ksql.CreateTableInterface {
	ta.table = table
	return ta
}

func (ta *Table) Engine(engine string) ksql.CreateTableInterface {
	ta.engine = engine
	return ta
}

func (ta *Table) Charset(charset string) ksql.CreateTableInterface {
	ta.charset = charset
	return ta
}

func (ta *Table) Collate(collate string) ksql.CreateTableInterface {
	ta.collate = collate
	return ta
}

func (ta *Table) Comment(comment string) ksql.CreateTableInterface {
	ta.comment = comment
	return ta
}

func (ta *Table) Prepare() string {
	if ta.hasPrepared {
		return ta.base.Prepare()
	}

	ta.hasPrepared = true
	Column(ta.table, &ta.builder)
	ta.builder.WriteString(" (")
	for idx, column := range ta.columns {
		if idx > 0 {
			ta.builder.WriteString(",")
		}

		ta.builder.WriteString(column.Express())
	}

	for _, index := range ta.indexes {
		ta.builder.WriteString(",")
		ta.builder.WriteString(index.Express())
	}

	ta.builder.WriteString(") ENGINE=")
	ta.builder.WriteString(ta.engine)
	ta.builder.WriteString(" DEFAULT CHARSET=")
	ta.builder.WriteString(ta.charset)
	ta.builder.WriteString(" COLLATE=")
	ta.builder.WriteString(ta.collate)
	if ta.comment != "" {
		ta.builder.WriteString(" COMMENT=")
		Quote(ta.comment, &ta.builder)
	}

	return ta.base.Prepare()
}
