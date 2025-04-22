package db

import (
	"database/sql"
	"time"

	ksql "github.com/kovey/db-go/v3"
)

type Data struct {
	data map[string]any
	keys []string
}

func NewData() *Data {
	return &Data{data: make(map[string]any)}
}

func (d *Data) From(o *Data) {
	o.Range(func(key string, val any) {
		if _, ok := d.data[key]; !ok {
			d.keys = append(d.keys, key)
		}
		d.data[key] = val
	})
}

func (d *Data) Set(key string, val any) *Data {
	if _, ok := d.data[key]; !ok {
		d.keys = append(d.keys, key)
	}

	switch tmp := val.(type) {
	case *string:
		d.data[key] = *tmp
	case **string:
		d.data[key] = *tmp
	case *int:
		d.data[key] = *tmp
	case **int:
		d.data[key] = *tmp
	case *int8:
		d.data[key] = *tmp
	case **int8:
		d.data[key] = *tmp
	case *int16:
		d.data[key] = *tmp
	case **int16:
		d.data[key] = *tmp
	case *int32:
		d.data[key] = *tmp
	case **int32:
		d.data[key] = *tmp
	case *int64:
		d.data[key] = *tmp
	case **int64:
		d.data[key] = *tmp
	case *uint:
		d.data[key] = *tmp
	case **uint:
		d.data[key] = *tmp
	case *uint8:
		d.data[key] = *tmp
	case **uint8:
		d.data[key] = *tmp
	case *uint16:
		d.data[key] = *tmp
	case **uint16:
		d.data[key] = *tmp
	case *uint32:
		d.data[key] = *tmp
	case **uint32:
		d.data[key] = *tmp
	case *uint64:
		d.data[key] = *tmp
	case **uint64:
		d.data[key] = *tmp
	case *bool:
		d.data[key] = *tmp
	case **bool:
		d.data[key] = *tmp
	case *float32:
		d.data[key] = *tmp
	case **float32:
		d.data[key] = *tmp
	case *float64:
		d.data[key] = *tmp
	case **float64:
		d.data[key] = *tmp
	case *time.Time:
		d.data[key] = *tmp
	case **time.Time:
		d.data[key] = *tmp
	default:
		d.data[key] = val
	}

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

func (d *Data) Changed(key string, val any) bool {
	old, ok := d.data[key]
	if !ok {
		return true
	}

	switch tmp := val.(type) {
	case **string:
		v, ok := old.(*string)
		if !ok {
			return true
		}
		return v != *tmp
	case **int:
		v, ok := old.(*int)
		if !ok {
			return true
		}
		return v != *tmp
	case **int8:
		v, ok := old.(*int8)
		if !ok {
			return true
		}
		return v != *tmp
	case **int16:
		v, ok := old.(*int16)
		if !ok {
			return true
		}
		return v != *tmp
	case **int32:
		v, ok := old.(*int32)
		if !ok {
			return true
		}
		return v != *tmp
	case **int64:
		v, ok := old.(*int64)
		if !ok {
			return true
		}
		return v != *tmp
	case **uint:
		v, ok := old.(*uint)
		if !ok {
			return true
		}
		return v != *tmp
	case **uint8:
		v, ok := old.(*uint8)
		if !ok {
			return true
		}
		return v != *tmp
	case **uint16:
		v, ok := old.(*uint16)
		if !ok {
			return true
		}
		return v != *tmp
	case **uint32:
		v, ok := old.(*uint32)
		if !ok {
			return true
		}
		return v != *tmp
	case **uint64:
		v, ok := old.(*uint64)
		if !ok {
			return true
		}
		return v != *tmp
	case **bool:
		v, ok := old.(*bool)
		if !ok {
			return true
		}
		return v != *tmp
	case **float32:
		v, ok := old.(*float32)
		if !ok {
			return true
		}
		return v != *tmp
	case **float64:
		v, ok := old.(*float64)
		if !ok {
			return true
		}
		return v != *tmp
	case **time.Time:
		v, ok := old.(*time.Time)
		if !ok {
			return true
		}
		return !v.Equal(**tmp)
	case *string:
		v, ok := old.(string)
		if !ok {
			return true
		}
		return v != *tmp
	case *int:
		v, ok := old.(int)
		if !ok {
			return true
		}
		return v != *tmp
	case *int8:
		v, ok := old.(int8)
		if !ok {
			return true
		}
		return v != *tmp
	case *int16:
		v, ok := old.(int16)
		if !ok {
			return true
		}
		return v != *tmp
	case *int32:
		v, ok := old.(int32)
		if !ok {
			return true
		}
		return v != *tmp
	case *int64:
		v, ok := old.(int64)
		if !ok {
			return true
		}
		return v != *tmp
	case *uint:
		v, ok := old.(uint)
		if !ok {
			return true
		}
		return v != *tmp
	case *uint8:
		v, ok := old.(uint8)
		if !ok {
			return true
		}
		return v != *tmp
	case *uint16:
		v, ok := old.(uint16)
		if !ok {
			return true
		}
		return v != *tmp
	case *uint32:
		v, ok := old.(uint32)
		if !ok {
			return true
		}
		return v != *tmp
	case *uint64:
		v, ok := old.(uint64)
		if !ok {
			return true
		}
		return v != *tmp
	case *bool:
		v, ok := old.(bool)
		if !ok {
			return true
		}
		return v != *tmp
	case *float32:
		v, ok := old.(float32)
		if !ok {
			return true
		}
		return v != *tmp
	case *float64:
		v, ok := old.(float64)
		if !ok {
			return true
		}
		return v != *tmp
	case *time.Time:
		v, ok := old.(time.Time)
		if !ok {
			return true
		}
		return !v.Equal(*tmp)
	case *sql.NullBool:
		v, ok := old.(sql.NullBool)
		if !ok {
			return true
		}
		return v.Bool != tmp.Bool
	case *sql.NullByte:
		v, ok := old.(sql.NullByte)
		if !ok {
			return true
		}
		return v.Byte != tmp.Byte
	case *sql.NullFloat64:
		v, ok := old.(sql.NullFloat64)
		if !ok {
			return true
		}
		return v.Float64 != tmp.Float64
	case *sql.NullInt16:
		v, ok := old.(sql.NullInt16)
		if !ok {
			return true
		}
		return v.Int16 != tmp.Int16
	case *sql.NullInt32:
		v, ok := old.(sql.NullInt32)
		if !ok {
			return true
		}
		return v.Int32 != tmp.Int32
	case *sql.NullInt64:
		v, ok := old.(sql.NullInt64)
		if !ok {
			return true
		}
		return v.Int64 != tmp.Int64
	case *sql.NullTime:
		v, ok := old.(sql.NullTime)
		if !ok {
			return true
		}
		return v.Time.Equal(tmp.Time)
	default:
		return old != val
	}
}

func (d *Data) Empty() bool {
	return len(d.keys) == 0
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
