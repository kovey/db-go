package itf

type RowInterface interface {
	Columns() []string
	Fields() []any
	Clone() RowInterface
}

type ModelInterface interface {
	RowInterface
	Values() []any
	SetEmpty()
}

type Row struct {
}

func (r *Row) Columns() []string {
	return nil
}

func (r *Row) Fields() []any {
	return nil
}

func (r *Row) Clone() RowInterface {
	return &Row{}
}

func (r *Row) Values() []any {
	return nil
}

func (r *Row) SetEmpty() {
}
