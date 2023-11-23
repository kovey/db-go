package desc

import (
	"database/sql"
	"fmt"

	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/model"
	"github.com/kovey/db-go/v2/table"
	"github.com/kovey/pool/object"
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

func (t *Desc) GetComment() string {
	if t.Comment.String == "" {
		return ""
	}

	return fmt.Sprintf("// %s", t.Comment.String)
}

func (t *Desc) Columns() []string {
	return []string{"COLUMN_NAME", "COLUMN_TYPE", "IS_NULLABLE", "COLUMN_KEY", "COLUMN_DEFAULT", "EXTRA", "COLUMN_COMMENT"}
}

func (t *Desc) Fields() []any {
	return []any{&t.Field, &t.Type, &t.Null, &t.Key, &t.Default, &t.Extra, &t.Comment}
}

func (t *Desc) Values() []any {
	return []any{t.Field, t.Type, t.Null, t.Key, t.Default, t.Extra, t.Comment}
}

func (t *Desc) Clone(object.CtxInterface) itf.RowInterface {
	return NewDesc()
}

type DescTable struct {
	*table.Table[*Desc]
}

func NewDescTable() *DescTable {
	return &DescTable{Table: table.NewTable[*Desc]("information_schema.COLUMNS")}
}
