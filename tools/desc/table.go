package desc

import (
	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/model"
	"github.com/kovey/db-go/v2/table"
)

type Table struct {
	*model.Base[*Table]
	Name string
}

func (t *Table) Columns() []string {
	return []string{"tables"}
}

func (t *Table) Fields() []any {
	return []any{&t.Name}
}

func (t *Table) Values() []any {
	return []any{t.Name}
}

func (t *Table) Clone() itf.RowInterface {
	return &Table{}
}

func NewTable() *Table {
	return &Table{Base: model.NewBase[*Table](table.NewTable[*Table](""), model.NewPrimaryId("Field", model.Str))}
}
