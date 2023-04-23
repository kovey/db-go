package meta

type Data map[string]any

func NewData() Data {
	return make(Data)
}

func (d Data) Add(name string, data any) {
	d[name] = data
}
