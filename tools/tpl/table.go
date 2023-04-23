package tpl

const (
	Table = `package {package_name}

import(
	"github.com/kovey/db-go/v2/table"
	"github.com/kovey/db-go/v2/model"
	{imports}
)

type {name}Table struct {
	*table.Table[*{name}Row]
}

func New{name}Table() *{name}Table {
	return &{name}Table{Table: table.NewTable[*{name}Row]("{table_name}")}
}

type {name}Row struct {
	*model.Base[*{name}Row]
{row_fields}
}

func New{name}Row() *{name}Row {
	return &{name}Row{Base: model.NewBase[*{name}Row](New{name}Table(), model.NewPrimaryId("{primary_id}", model.{primary_id_type}))}
}
	`
	Field = "	%s %s `db:\"%s\"`"

	Decimal = `"github.com/shopspring/decimal"`
	Sql     = `"database/sql"`
)
