package rows

import (
	"database/sql"
	"reflect"
)

const (
	Tag_Db = "db"
)

type Row[T any] struct {
	Model     T
	Columns   []any
	Fields    []string
	v         reflect.Value
	isPointer bool
}

func NewRow[T any](m T) *Row[T] {
	r := &Row[T]{}
	mType := reflect.TypeOf(m)
	if mType.Kind() == reflect.Ptr {
		mType = mType.Elem()
	}

	r.v = reflect.New(mType)
	vType := r.v.Type()
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
		r.v = r.v.Elem()
		r.isPointer = true
	}

	fLen := vType.NumField()
	r.Columns = make([]any, 0)
	r.Fields = make([]string, 0)

	for i := 0; i < fLen; i++ {
		field := r.v.Field(i)
		name := vType.Field(i).Tag.Get(Tag_Db)
		if name == "" {
			continue
		}

		r.Columns = append(r.Columns, field.Addr().Interface())
		r.Fields = append(r.Fields, name)
	}

	return r
}

func (r *Row[T]) Scan(row *sql.Row) error {
	if err := row.Scan(r.Columns...); err != nil {
		return err
	}

	r.setModel()
	return nil
}

func (r *Row[T]) ScanByRows(rows *sql.Rows) error {
	if err := rows.Scan(r.Columns...); err != nil {
		return err
	}

	r.setModel()
	return nil
}

func (r *Row[T]) setModel() {
	if r.isPointer {
		r.Model = r.v.Addr().Interface().(T)
		return
	}

	r.Model = r.v.Interface().(T)
}
