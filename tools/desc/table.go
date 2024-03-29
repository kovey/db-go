package desc

import (
	"database/sql"
	"fmt"

	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/model"
	"github.com/kovey/db-go/v2/table"
	"github.com/kovey/pool/object"
)

type Table struct {
	*model.Base[*Table]
	Name    string         `db:"TABLE_NAME"`
	Comment sql.NullString `db:"TABLE_COMMENT"`
}

func (t *Table) GetComment() string {
	if t.Comment.String == "" {
		return ""
	}

	return fmt.Sprintf("// %s", t.Comment.String)
}

func (t *Table) Columns() []string {
	return []string{"TABLE_NAME", "TABLE_COMMENT"}
}

func (t *Table) Fields() []any {
	return []any{&t.Name, &t.Comment}
}

func (t *Table) Values() []any {
	return []any{t.Name, t.Comment}
}

func (t *Table) Clone(object.CtxInterface) itf.RowInterface {
	return NewTable()
}

func NewTable() *Table {
	return &Table{Base: model.NewBase[*Table](NewTableTable(), model.NewPrimaryId("Field", model.Str))}
}

type TableTable struct {
	*table.Table[*Table]
}

func NewTableTable() *TableTable {
	return &TableTable{Table: table.NewTable[*Table]("information_schema.TABLES")}
}
