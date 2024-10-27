package db

import ksql "github.com/kovey/db-go/v3"

type Data struct {
	data map[string]any
	keys []string
}

func NewData() *Data {
	return &Data{data: make(map[string]any)}
}

func (d *Data) Set(key string, val any) *Data {
	if _, ok := d.data[key]; !ok {
		d.keys = append(d.keys, key)
	}

	d.data[key] = val
	return d
}

func (d *Data) Keys() []string {
	return d.keys
}

func (d *Data) Get(key string) any {
	return d.data[key]
}

func (d *Data) Range(call func(key string, val any)) {
	for _, key := range d.keys {
		call(key, d.data[key])
	}
}

type Map[K comparable, T any] struct {
	data map[K]T
	keys []K
}

func NewMap[K comparable, T any]() *Map[K, T] {
	return &Map[K, T]{data: make(map[K]T)}
}

func (d *Map[K, T]) Set(key K, val T) *Map[K, T] {
	if _, ok := d.data[key]; !ok {
		d.keys = append(d.keys, key)
	}

	d.data[key] = val
	return d
}

func (d *Map[K, T]) Keys() []K {
	return d.keys
}

func (d *Map[K, T]) Get(key K) T {
	return d.data[key]
}

func (d *Map[K, T]) Range(call func(key K, val T)) {
	for _, key := range d.keys {
		call(key, d.data[key])
	}
}

func (d *Map[K, T]) Values() []T {
	values := make([]T, len(d.keys))
	for i, key := range d.keys {
		values[i] = d.Get(key)
	}

	return values
}

func (m *Map[K, T]) Has(key K) bool {
	_, ok := m.data[key]
	return ok
}

func (m *Map[K, T]) GetBy(seq int) T {
	if seq >= len(m.keys) {
		var val T
		return val
	}

	return m.Get(m.keys[seq])
}

type PageInfo[T ksql.RowInterface] struct {
	list       []T
	totalCount uint64
	totalPage  uint64
}

func NewPageInfo[T ksql.RowInterface](list []T) *PageInfo[T] {
	return &PageInfo[T]{list: list}
}

func (p *PageInfo[T]) List() []T {
	return p.list
}

func (p *PageInfo[T]) TotalPage() uint64 {
	return p.totalPage
}

func (p *PageInfo[T]) TotalCount() uint64 {
	return p.totalCount
}

func (p *PageInfo[T]) Set(total, pageSize uint64) {
	p.totalCount = total
	p.totalPage = p.totalCount / uint64(pageSize)
	if p.totalCount%uint64(pageSize) != 0 {
		p.totalPage++
	}
}
