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
