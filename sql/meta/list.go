package meta

type List []string

func NewList() List {
	return make(List, 0)
}
