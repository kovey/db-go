package sql

import "github.com/kovey/db-go/v3"

type Update struct {
	*base
	columns []string
	where   ksql.WhereInterface
}

func NewUpdate() *Update {
	u := &Update{base: &base{hasPrepared: false}}
	u.keyword("UPDATE ")
	return u
}

func (u *Update) Table(table string) ksql.UpdateInterface {
	Column(table, &u.builder)
	u.builder.WriteString(" SET ")
	return u
}

func (u *Update) Set(column string, data any) ksql.UpdateInterface {
	u.columns = append(u.columns, column)
	u.binds = append(u.binds, data)
	return u
}

func (u *Update) Where(where ksql.WhereInterface) ksql.UpdateInterface {
	u.where = where
	return u
}

func (u *Update) Prepare() string {
	if u.hasPrepared {
		return u.base.Prepare()
	}

	u.hasPrepared = true
	for index, column := range u.columns {
		if index > 0 {
			u.builder.WriteString(",")
		}

		Column(column, &u.builder)
		u.builder.WriteString(" = ?")
	}

	if u.where != nil {
		u.builder.WriteString(" WHERE ")
		u.builder.WriteString(u.where.Prepare())
		u.binds = append(u.binds, u.where.Binds()...)
	}

	return u.base.Prepare()
}
