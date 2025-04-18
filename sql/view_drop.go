package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type DropView struct {
	*drop
	option string
}

func NewDropView() *DropView {
	ta := &DropView{drop: newDrop("VIEW")}
	ta.isMulti = true
	ta.opChain.Append(ta._buildOther)
	return ta
}

func (d *DropView) _buildOther(builder *strings.Builder) {
	operator.BuildPureString(d.option, builder)
}

func (d *DropView) View(table string) ksql.DropViewInterface {
	d.names = append(d.names, table)
	return d
}

func (d *DropView) IfExists() ksql.DropViewInterface {
	d.ifExists = true
	return d
}

func (d *DropView) Restrict() ksql.DropViewInterface {
	d.option = "RESTRICT"
	return d
}

func (d *DropView) Cascade() ksql.DropViewInterface {
	d.option = "CASCADE"
	return d
}
