package rows

import (
	"database/sql"
	"reflect"

	"github.com/kovey/db-go/v2/sql/meta"
)

const (
	Tag_Db = "db"
)

type Row[T any] struct {
	Model     T
	columns   map[string]any
	Fields    []*meta.Column
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
	r.columns = make(map[string]any)
	r.Fields = make([]*meta.Column, 0)

	for i := 0; i < fLen; i++ {
		field := r.v.Field(i)
		name := vType.Field(i).Tag.Get(Tag_Db)
		if name == "" {
			continue
		}

		r.columns[name] = field.Addr().Interface()
		r.Fields = append(r.Fields, meta.NewColumn(name, name))
	}

	return r
}

func (r *Row[T]) Columns() []any {
	res := make([]any, len(r.Fields))
	for index, name := range r.Fields {
		res[index] = r.columns[name.Name.Name]
	}

	return res
}

func (r *Row[T]) Scan(row *sql.Row) error {
	if err := row.Scan(r.Columns()...); err != nil {
		return err
	}

	r.setModel()
	return nil
}

func (r *Row[T]) getColumnsBy(fields []string) []any {
	res := make([]any, len(fields))
	for index, name := range fields {
		res[index] = r.columns[name]
	}

	return res
}

func (r *Row[T]) ScanByRows(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	if err := rows.Scan(r.getColumnsBy(cols)...); err != nil {
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
