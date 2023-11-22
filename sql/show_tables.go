package sql

import (
	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

const (
	showTablesFormat = "SHOW TABLES"
	show_name        = "Show"
)

func init() {
	pool.DefaultNoCtx(namespace, show_name, func() any {
		return &ShowTables{object.NewObjNoCtx(namespace, show_name)}
	})
}

type ShowTables struct {
	*object.ObjNoCtx
}

func NewShowTables() *ShowTables {
	return &ShowTables{}
}

func NewShowTablesBy(ctx object.CtxInterface) *ShowTables {
	return ctx.GetNoCtx(namespace, show_name).(*ShowTables)
}

func (d *ShowTables) Args() []any {
	return []any{}
}

func (d *ShowTables) Prepare() string {
	return showTablesFormat
}

func (d *ShowTables) String() string {
	return String(d)
}
