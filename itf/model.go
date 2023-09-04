package itf

import "github.com/kovey/db-go/v2/sql/meta"

type RowInterface interface {
	Columns() []*meta.Column
	Fields() []any
	Clone() RowInterface
}

type ModelInterface interface {
	RowInterface
	Values() []any
	SetEmpty()
}
