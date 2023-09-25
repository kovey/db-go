package meta

type Page[T any] struct {
	List       []T
	TotalCount int64
	TotalPage  int64
}

func NewPage[T any](list []T) *Page[T] {
	return &Page[T]{List: list}
}
