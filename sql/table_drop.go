package sql

import "github.com/kovey/db-go/v3"

type DropTable struct {
	*base
}

func NewDropTable() *DropTable {
	ta := &DropTable{base: &base{hasPrepared: false}}
	ta.keyword("DROP TABLE ")
	return ta
}

func (d *DropTable) Table(table string) ksql.DropTableInterface {
	Column(table, &d.builder)
	return d
}
