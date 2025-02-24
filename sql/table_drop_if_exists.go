package sql

import ksql "github.com/kovey/db-go/v3"

type DropTableIfExists struct {
	*base
}

func NewDropTableIfExists() *DropTable {
	ta := &DropTable{base: &base{hasPrepared: false}}
	ta.keyword("DROP TABLE IF EXISTS ")
	return ta
}

func (d *DropTableIfExists) Table(table string) ksql.DropTableInterface {
	Column(table, &d.builder)
	return d
}
