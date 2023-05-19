package desc

import (
	"database/sql"

	"github.com/kovey/db-go/v2/model"
	"github.com/kovey/db-go/v2/table"
)

type Desc struct {
	*model.Base[*Desc]
	Field   string         `db:"COLUMN_NAME"`
	Type    string         `db:"COLUMN_TYPE"`
	Null    string         `db:"IS_NULLABLE"`
	Key     sql.NullString `db:"COLUMN_KEY"`
	Default sql.NullString `db:"COLUMN_DEFAULT"`
	Extra   sql.NullString `db:"EXTRA"`
	Comment sql.NullString `db:"COLUMN_COMMENT"`
}

func NewDesc() *Desc {
	return &Desc{Base: model.NewBase[*Desc](NewDescTable(), model.NewPrimaryId("Field", model.Str))}
}

type DescTable struct {
	*table.Table[*Desc]
}

func NewDescTable() *DescTable {
	return &DescTable{Table: table.NewTable[*Desc]("information_schema.COLUMNS")}
}
