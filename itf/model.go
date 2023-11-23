package itf

import "github.com/kovey/pool/object"

type RowInterface interface {
	Columns() []string
	Fields() []any
	Clone(object.CtxInterface) RowInterface
}

type ModelInterface interface {
	RowInterface
	Values() []any
	SetEmpty()
	SetFetch()
}

type Row struct {
}

func (r *Row) Columns() []string {
	return nil
}

func (r *Row) Fields() []any {
	return nil
}

func (r *Row) Clone(object.CtxInterface) RowInterface {
	return &Row{}
}

func (r *Row) Values() []any {
	return nil
}

func (r *Row) SetEmpty() {
}
