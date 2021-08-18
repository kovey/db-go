package sql

import "fmt"

const (
	deleteFormat string = "DELETE FROM %s %s"
)

type Delete struct {
	table string
	where *Where
}

func NewDelete(table string) *Delete {
	return &Delete{table: table, where: nil}
}

func (d *Delete) Where(w *Where) *Delete {
	d.where = w
	return d
}

func (d *Delete) Args() []interface{} {
	if d.where == nil {
		return []interface{}{}
	}

	return d.where.Args()
}

func (d *Delete) Prepare() string {
	if d.where == nil {
		return fmt.Sprintf(deleteFormat, formatValue(d.table), "")
	}
	return fmt.Sprintf(deleteFormat, formatValue(d.table), d.where.Prepare())
}

func (d *Delete) String() string {
	return String(d)
}

func (d *Delete) WhereByMap(where map[string]interface{}) *Delete {
	if d.where == nil {
		d.where = NewWhere()
	}

	for field, value := range where {
		d.where.Eq(field, value)
	}

	return d
}

func (d *Delete) WhereByList(where []string) *Delete {
	if d.where == nil {
		d.where = NewWhere()
	}

	for _, value := range where {
		d.where.Statement(value)
	}

	return d
}
