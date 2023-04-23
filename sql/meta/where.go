package meta

type Where map[string]any

func NewWhere() Where {
	return make(Where)
}

func (w Where) Add(name string, data any) {
	w[name] = data
}
