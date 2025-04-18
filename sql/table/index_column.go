package table

import (
	"strconv"
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type IndexColumnType byte

const (
	Index_Column_Type_Name      IndexColumnType = 0x0
	Index_Column_Type_Expr      IndexColumnType = 0x1
	Index_Column_Type_Pure_Name IndexColumnType = 0x2
)

type IndexColumn struct {
	Type   IndexColumnType
	Name   string
	Length int
	Order  ksql.Order
}

type IndexColumns struct {
	columns []*IndexColumn
}

func (i *IndexColumns) Append(column *IndexColumn) *IndexColumns {
	i.columns = append(i.columns, column)
	return i
}

func (i *IndexColumns) Build(builder *strings.Builder) {
	builder.WriteString(" (")
	for index, column := range i.columns {
		if index > 0 {
			builder.WriteString(", ")
		}

		switch column.Type {
		case Index_Column_Type_Name:
			operator.Backtick(column.Name, builder)
			if column.Length > 0 {
				builder.WriteString("(")
				builder.WriteString(strconv.Itoa(column.Length))
				builder.WriteString(")")
			}
		case Index_Column_Type_Expr:
			builder.WriteString("(")
			builder.WriteString(column.Name)
			builder.WriteString(")")
		case Index_Column_Type_Pure_Name:
			operator.Backtick(column.Name, builder)
		}

		if column.Order != ksql.Order_None {
			builder.WriteString(" ")
			builder.WriteString(string(column.Order))
		}
	}
	builder.WriteString(")")
}
