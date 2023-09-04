package desc

import (
	"database/sql"

	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/model"
	"github.com/kovey/db-go/v2/sql/meta"
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

func (t *Desc) Columns() []*meta.Column {
	return []*meta.Column{
		meta.NewColumn("COLUMN_NAME"), meta.NewColumn("COLUMN_TYPE"), meta.NewColumn("IS_NULLABLE"),
		meta.NewColumn("COLUMN_KEY"), meta.NewColumn("COLUMN_DEFAULT"), meta.NewColumn("EXTRA"), meta.NewColumn("COLUMN_COMMENT"),
	}
}

func (t *Desc) Fields() []any {
	return []any{&t.Field, &t.Type, &t.Null, &t.Key, &t.Default, &t.Extra, &t.Comment}
}

func (t *Desc) Values() []any {
	return []any{t.Field, t.Type, t.Null, t.Key, t.Default, t.Extra, t.Comment}
}

func (t *Desc) Clone() itf.RowInterface {
	return &Desc{}
}

type DescTable struct {
	*table.Table[*Desc]
}

func NewDescTable() *DescTable {
	return &DescTable{Table: table.NewTable[*Desc]("information_schema.COLUMNS")}
}
