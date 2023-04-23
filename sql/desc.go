package sql

import "fmt"

const (
	descFormat = "DESC %s"
)

type Desc struct {
	table string
}

func NewDesc(table string) *Desc {
	return &Desc{table: table}
}

func (d *Desc) Args() []any {
	return []any{}
}

func (d *Desc) Prepare() string {
	return fmt.Sprintf(descFormat, formatValue(d.table))
}

func (d *Desc) String() string {
	return String(d)
}
