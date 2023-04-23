package desc

import (
	"database/sql"

	"github.com/kovey/db-go/v2/model"
	"github.com/kovey/db-go/v2/table"
)

type Desc struct {
	*model.Base[*Desc]
	Field   string         `db:"Field"`
	Type    string         `db:"Type"`
	Null    string         `db:"Null"`
	Key     sql.NullString `db:"Key"`
	Default sql.NullString `db:"Default"`
	Extra   sql.NullString `db:"Extra"`
}

func NewDesc(name string) *Desc {
	return &Desc{Base: model.NewBase[*Desc](table.NewTable[*Desc](name), model.NewPrimaryId("Field", model.Str))}
}
