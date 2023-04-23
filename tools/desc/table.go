package desc

import (
	"github.com/kovey/db-go/v2/model"
	"github.com/kovey/db-go/v2/table"
)

type Table struct {
	*model.Base[*Table]
	Name string `db:"tables"`
}

func NewTable() *Table {
	return &Table{Base: model.NewBase[*Table](table.NewTable[*Table](""), model.NewPrimaryId("Field", model.Str))}
}
