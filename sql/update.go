package sql

import (
	"fmt"
	"strings"
)

const (
	updateFormat   string = "UPDATE %s SET %s %s"
	updateCkFormat string = "ALTER TABLE %s UPDATE %s %s"
)

type Update struct {
	table  string
	data   map[string]interface{}
	args   []interface{}
	where  *Where
	format string
}

func NewUpdate(table string) *Update {
	return &Update{table: table, data: make(map[string]interface{}), where: nil, format: updateFormat}
}

func NewCkUpdate(table string) *Update {
	return &Update{table: table, data: make(map[string]interface{}), where: nil, format: updateCkFormat}
}

func (u *Update) Set(field string, value interface{}) *Update {
	u.data[field] = value
	return u
}

func (u *Update) Args() []interface{} {
	if u.where == nil {
		return u.args
	}

	return append(u.args, u.where.Args()...)
}

func (u *Update) getPlaceholder() []string {
	placeholders := make([]string, len(u.data))
	u.args = make([]interface{}, len(u.data))
	index := 0
	for field, v := range u.data {
		t, ok := v.(string)
		if !ok {
			placeholders[index] = fmt.Sprintf("%s = ?", formatValue(field))
			u.args[index] = v
			index++
			continue
		}

		var value = t
		var op = "="
		if len(value) > 2 {
			prefix := t[0:2]
			if prefix == "+=" {
				value = t[2:]
				op = "+="
			} else if prefix == "-=" {
				value = t[2:]
				op = "-="
			}
		}

		u.args[index] = value
		placeholders[index] = fmt.Sprintf("%s %s ?", formatValue(field), op)
		index++
	}

	return placeholders
}

func (u *Update) Prepare() string {
	if u.where == nil {
		return fmt.Sprintf(u.format, formatValue(u.table), strings.Join(u.getPlaceholder(), ","), "")
	}

	return fmt.Sprintf(u.format, formatValue(u.table), strings.Join(u.getPlaceholder(), ","), u.where.Prepare())
}

func (u *Update) Where(w *Where) *Update {
	u.where = w
	return u
}

func (u *Update) WhereByMap(where map[string]interface{}) *Update {
	if u.where == nil {
		u.where = NewWhere()
	}

	for field, value := range where {
		u.where.Eq(field, value)
	}

	return u
}

func (u *Update) WhereByList(where []string) *Update {
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
