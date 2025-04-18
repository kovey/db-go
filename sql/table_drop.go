package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type DropTable struct {
	*drop
	option string
}

func NewDropTable() *DropTable {
	ta := &DropTable{drop: newDrop("TABLE")}
	ta.isMulti = true
	ta.opChain.Append(ta._buildOther)
	return ta
}

func (d *DropTable) _buildOther(builder *strings.Builder) {
	operator.BuildPureString(d.option, builder)
}

func (d *DropTable) Table(table string) ksql.DropTableInterface {
	d.names = append(d.names, table)
	return d
}

func (d *DropTable) IfExists() ksql.DropTableInterface {
	d.ifExists = true
	return d
}

func (d *DropTable) Temporary() ksql.DropTableInterface {
	d.keyword = "TEMPORARY TABLE"
	return d
}

func (d *DropTable) Restrict() ksql.DropTableInterface {
	d.option = "RESTRICT"
	return d
}

func (d *DropTable) Cascade() ksql.DropTableInterface {
	d.option = "CASCADE"
	return d
}
