package row

import (
	"database/sql"
	"reflect"

	"github.com/kovey/logger-go/logger"
)

type Row struct {
	columns []interface{}
	value   reflect.Value
	fields  []string
}

func New(t reflect.Type) *Row {
	value := reflect.New(t).Elem()
	vType := value.Type()
	columns := make([]interface{}, 0)
	fields := make([]string, 0)

	for i := 0; i < vType.NumField(); i++ {
		field := value.Field(i)
		name := vType.Field(i).Tag.Get("db")
		if len(name) == 0 {
			continue
		}

		if field.Kind() == reflect.Ptr {
			logger.Debug("field[%s] is ptr", name)
			field = field.Elem()
			logger.Debug("field: %v", field)
		}

		columns = append(columns, field.Addr().Interface())
		fields = append(fields, name)
	}

	logger.Debug("fields: %v", fields)
	return &Row{columns: columns, value: value, fields: fields}
}

func (r *Row) Convert(rows *sql.Rows) {
	rows.Scan(r.columns...)
}

func (r *Row) ConvertByRow(row *sql.Row) {
	row.Scan(r.columns...)
}

func (r *Row) Value() interface{} {
	return r.value.Interface()
}

func (r *Row) Fields() []string {
	return r.fields
}

func (r *Row) Addr() reflect.Value {
	return r.value
}
