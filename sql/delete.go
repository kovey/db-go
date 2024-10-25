package sql

import "github.com/kovey/db-go/v3"

type Delete struct {
	*base
	where ksql.WhereInterface
}

func NewDelete() *Delete {
	d := &Delete{base: &base{hasPrepared: false}}
	d.keyword("DELETE FROM ")
	return d
}

func (d *Delete) Where(where ksql.WhereInterface) ksql.DeleteInterface {
	d.where = where
	return d
}

func (d *Delete) Table(table string) ksql.DeleteInterface {
	Column(table, &d.builder)
	return d
}

func (d *Delete) Prepare() string {
	if d.hasPrepared || d.where == nil || d.where.Empty() {
		return d.base.Prepare()
	}

	d.hasPrepared = true
	d.builder.WriteString(" WHERE ")
	d.builder.WriteString(d.where.Prepare())
	d.binds = append(d.binds, d.where.Binds()...)
	return d.base.Prepare()
}
