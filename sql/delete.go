package sql

import "fmt"

const (
	deleteFormat   string = "DELETE FROM %s %s"
	deleteCkFormat string = "ALTER TABLE %s DELETE %s"
)

type Delete struct {
	table  string
	where  WhereInterface
	format string
}

func NewDelete(table string) *Delete {
	return &Delete{table: table, where: nil, format: deleteFormat}
}

func NewCkDelete(table string) *Delete {
	return &Delete{table: table, where: nil, format: deleteCkFormat}
}

func (d *Delete) Where(w WhereInterface) *Delete {
	d.where = w
	return d
}

func (d *Delete) Args() []any {
	if d.where == nil {
		return []any{}
	}

	return d.where.Args()
}

func (d *Delete) Prepare() string {
	if d.where == nil {
		return fmt.Sprintf(d.format, formatValue(d.table), "")
	}
	return fmt.Sprintf(d.format, formatValue(d.table), d.where.Prepare())
}

func (d *Delete) String() string {
	return String(d)
}

func (d *Delete) WhereByMap(where map[string]any) *Delete {
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
