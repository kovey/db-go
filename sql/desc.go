package sql

import (
	"fmt"

	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

const (
	descFormat = "DESC %s"
	desc_name  = "Desc"
)

func init() {
	pool.DefaultNoCtx(namespace, desc_name, func() any {
		return &Desc{ObjNoCtx: object.NewObjNoCtx(namespace, desc_name)}
	})
}

type Desc struct {
	*object.ObjNoCtx
	table string
}

func NewDesc(table string) *Desc {
	return &Desc{table: table}
}

func NewDescBy(ctx object.CtxInterface, table string) *Desc {
	obj := ctx.GetNoCtx(namespace, desc_name).(*Desc)
	obj.table = table
	return obj
}

func (d *Desc) Reset() {
	d.table = emptyStr
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
