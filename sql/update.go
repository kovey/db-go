package sql

import (
	"fmt"
	"strings"

	"github.com/kovey/db-go/v2/sql/meta"
)

const (
	updateFormat      = "UPDATE %s SET %s %s"
	updateCkFormat    = "ALTER TABLE %s UPDATE %s %s"
	updatePlaceFormat = "%s = ?"
	addEq             = "+="
	updateAddFormat   = "= %s +"
	subEq             = "-="
	updateSubFormat   = "= %s -"
)

type Update struct {
	table  string
	data   meta.Data
	args   []any
	where  WhereInterface
	format string
}

func NewUpdate(table string) *Update {
	return &Update{table: table, data: meta.NewData(), where: nil, format: updateFormat}
}

func NewCkUpdate(table string) *Update {
	return &Update{table: table, data: meta.NewData(), where: nil, format: updateCkFormat}
}

func (u *Update) Set(field string, value any) *Update {
	u.data[field] = value
	return u
}

func (u *Update) Args() []any {
	if u.where == nil {
		return u.args
	}

	return append(u.args, u.where.Args()...)
}

func (u *Update) getPlaceholder() []string {
	placeholders := make([]string, len(u.data))
	u.args = make([]any, len(u.data))
	index := 0
	for field, v := range u.data {
		t, ok := v.(string)
		if !ok {
			placeholders[index] = fmt.Sprintf(updatePlaceFormat, formatValue(field))
			u.args[index] = v
			index++
			continue
		}

		var value = t
		var op = eq
		if len(value) > 2 {
			prefix := t[0:2]
			if prefix == addEq {
				value = t[2:]
				op = fmt.Sprintf(updateAddFormat, field)
			} else if prefix == subEq {
				value = t[2:]
				op = fmt.Sprintf(updateSubFormat, field)
			}
		}

		u.args[index] = value
		placeholders[index] = fmt.Sprintf(whereFields, formatValue(field), op)
		index++
	}

	return placeholders
}

func (u *Update) Prepare() string {
	if u.where == nil {
		return fmt.Sprintf(u.format, formatValue(u.table), strings.Join(u.getPlaceholder(), comma), emptyStr)
	}

	return fmt.Sprintf(u.format, formatValue(u.table), strings.Join(u.getPlaceholder(), comma), u.where.Prepare())
}

func (u *Update) Where(w WhereInterface) *Update {
	u.where = w
	return u
}

func (u *Update) WhereByMap(where meta.Where) *Update {
	if u.where == nil {
		u.where = NewWhere()
	}

	for field, value := range where {
		u.where.Eq(field, value)
	}

	return u
}

func (u *Update) WhereByList(where meta.List) *Update {
	if u.where == nil {
		u.where = NewWhere()
	}

	for _, value := range where {
		u.where.Statement(value)
	}

	return u
}

func (u *Update) String() string {
	return String(u)
}
