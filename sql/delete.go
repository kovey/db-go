package sql

import (
	"fmt"

	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

const (
	deleteFormat   string = "DELETE FROM %s %s"
	deleteCkFormat string = "ALTER TABLE %s DELETE %s"
	del_name              = "Delete"
)

func init() {
	pool.DefaultNoCtx(namespace, del_name, func() any {
		return &Delete{ObjNoCtx: object.NewObjNoCtx(namespace, del_name)}
	})
}

type Delete struct {
	*object.ObjNoCtx
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

func NewDeleteBy(ctx object.CtxInterface, table string) *Delete {
	obj := ctx.GetNoCtx(namespace, del_name).(*Delete)
	obj.table = table
	obj.format = deleteFormat
	return obj
}

func NewCkDeleteBy(ctx object.CtxInterface, table string) *Delete {
	obj := ctx.GetNoCtx(namespace, del_name).(*Delete)
	obj.table = table
	obj.format = deleteCkFormat
	return obj
}

func (d *Delete) Reset() {
	d.table = emptyStr
	d.where = nil
	d.format = emptyStr
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
		return fmt.Sprintf(d.format, formatValue(d.table), emptyStr)
	}
	return fmt.Sprintf(d.format, formatValue(d.table), d.where.Prepare())
}

func (d *Delete) String() string {
	return String(d)
}

func (d *Delete) WhereByMap(where meta.Where) *Delete {
	if d.where == nil {
		d.where = NewWhere()
	}

	for field, value := range where {
		d.where.Eq(field, value)
	}

	return d
}

func (d *Delete) WhereByList(where meta.List) *Delete {
	if d.where == nil {
		d.where = NewWhere()
	}

	for _, value := range where {
		d.where.Statement(value)
	}

	return d
}
