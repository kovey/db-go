package db

import (
	"context"

	ksql "github.com/kovey/db-go/v3"
)

type Builder[T ksql.RowInterface] struct {
	query  ksql.QueryInterface
	conn   ksql.ConnectionInterface
	model  T
	models *[]T
}

func NewBuilder[T ksql.RowInterface](model T) *Builder[T] {
	return &Builder[T]{query: NewQuery(), conn: model.Conn(), model: model}
}

func (b *Builder[T]) Sharding(sharding ksql.Sharding) ksql.BuilderInterface[T] {
	b.query.Sharding(sharding)
	return b
}

func (b *Builder[T]) Table(table string) ksql.BuilderInterface[T] {
	b.query.Table(table)
	return b
}

func (b *Builder[T]) TableBy(op ksql.QueryInterface, as string) ksql.BuilderInterface[T] {
	b.query.TableBy(op, as)
	return b
}

func (b *Builder[T]) As(as string) ksql.BuilderInterface[T] {
	b.query.As(as)
	return b
}

func (b *Builder[T]) Column(column, as string) ksql.BuilderInterface[T] {
	b.query.Column(column, as)
	return b
}

func (b *Builder[T]) Func(fun, column, as string) ksql.BuilderInterface[T] {
	b.query.Func(fun, column, as)
	return b
}

func (b *Builder[T]) Columns(columns ...string) ksql.BuilderInterface[T] {
	b.query.Columns(columns...)
	return b
}

func (b *Builder[T]) ColumnsExpress(expresses ...ksql.ExpressInterface) ksql.BuilderInterface[T] {
	b.query.ColumnsExpress(expresses...)
	return b
}

func (b *Builder[T]) Where(column string, op ksql.Op, val any) ksql.BuilderInterface[T] {
	b.query.Where(column, op, val)
	return b
}

func (b *Builder[T]) WhereExpress(express ...ksql.ExpressInterface) ksql.BuilderInterface[T] {
	b.query.WhereExpress(express...)
	return b
}

func (b *Builder[T]) OrWhere(call func(ksql.WhereInterface)) ksql.BuilderInterface[T] {
	b.query.OrWhere(call)
	return b
}

func (b *Builder[T]) WhereIsNull(column string) ksql.BuilderInterface[T] {
	b.query.WhereIsNull(column)
	return b
}

func (b *Builder[T]) WhereIsNotNull(column string) ksql.BuilderInterface[T] {
	b.query.WhereIsNull(column)
	return b
}

func (b *Builder[T]) WhereIn(column string, data []any) ksql.BuilderInterface[T] {
	b.query.WhereIn(column, data)
	return b
}

func (b *Builder[T]) WhereNotIn(column string, data []any) ksql.BuilderInterface[T] {
	b.query.WhereNotIn(column, data)
	return b
}

func (b *Builder[T]) WhereInBy(column string, sub ksql.QueryInterface) ksql.BuilderInterface[T] {
	b.query.WhereInBy(column, sub)
	return b
}

func (b *Builder[T]) WhereNotInBy(column string, sub ksql.QueryInterface) ksql.BuilderInterface[T] {
	b.query.WhereNotInBy(column, sub)
	return b
}

func (b *Builder[T]) AndWhere(call func(w ksql.WhereInterface)) ksql.BuilderInterface[T] {
	b.query.AndWhere(call)
	return b
}

func (b *Builder[T]) Between(column string, begin, end any) ksql.BuilderInterface[T] {
	b.query.Between(column, begin, end)
	return b
}

func (b *Builder[T]) NotBetween(column string, begin, end any) ksql.BuilderInterface[T] {
	b.query.NotBetween(column, begin, end)
	return b
}

func (b *Builder[T]) Having(column string, op ksql.Op, val any) ksql.BuilderInterface[T] {
	b.query.Having(column, op, val)
	return b
}

func (b *Builder[T]) HavingExpress(expresses ...ksql.ExpressInterface) ksql.BuilderInterface[T] {
	b.query.HavingExpress(expresses...)
	return b
}

func (b *Builder[T]) OrHaving(call func(ksql.HavingInterface)) ksql.BuilderInterface[T] {
	b.query.OrHaving(call)
	return b
}

func (b *Builder[T]) HavingIsNull(column string) ksql.BuilderInterface[T] {
	b.query.HavingIsNull(column)
	return b
}

func (b *Builder[T]) HavingIsNotNull(column string) ksql.BuilderInterface[T] {
	b.query.HavingIsNull(column)
	return b
}

func (b *Builder[T]) HavingIn(column string, data []any) ksql.BuilderInterface[T] {
	b.query.HavingIn(column, data)
	return b
}

func (b *Builder[T]) HavingNotIn(column string, data []any) ksql.BuilderInterface[T] {
	b.query.HavingNotIn(column, data)
	return b
}

func (b *Builder[T]) HavingInBy(column string, sub ksql.QueryInterface) ksql.BuilderInterface[T] {
	b.query.HavingInBy(column, sub)
	return b
}

func (b *Builder[T]) HavingNotInBy(column string, sub ksql.QueryInterface) ksql.BuilderInterface[T] {
	b.query.HavingNotInBy(column, sub)
	return b
}

func (b *Builder[T]) HavingBetween(column string, begin, end any) ksql.BuilderInterface[T] {
	b.query.HavingBetween(column, begin, end)
	return b
}

func (b *Builder[T]) HavingNotBetween(column string, begin, end any) ksql.BuilderInterface[T] {
	b.query.HavingNotBetween(column, begin, end)
	return b
}

func (b *Builder[T]) AndHaving(call func(w ksql.HavingInterface)) ksql.BuilderInterface[T] {
	b.query.AndHaving(call)
	return b
}

func (b *Builder[T]) Distinct() ksql.BuilderInterface[T] {
	b.query.Distinct()
	return b
}

func (b *Builder[T]) FuncDistinct(fun, column, as string) ksql.BuilderInterface[T] {
	b.query.FuncDistinct(fun, column, as)
	return b
}

func (b *Builder[T]) Limit(limit int) ksql.BuilderInterface[T] {
	b.query.Limit(limit)
	return b
}

func (b *Builder[T]) Offset(offset int) ksql.BuilderInterface[T] {
	b.query.Offset(offset)
	return b
}

func (b *Builder[T]) Order(column string) ksql.BuilderInterface[T] {
	b.query.Order(column)
	return b
}

func (b *Builder[T]) OrderDesc(column string) ksql.BuilderInterface[T] {
	b.query.OrderDesc(column)
	return b
}

func (b *Builder[T]) Group(column ...string) ksql.BuilderInterface[T] {
	b.query.Group(column...)
	return b
}

func (b *Builder[T]) Join(table string) ksql.JoinInterface {
	return b.query.Join(table)
}

func (b *Builder[T]) JoinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	return b.query.JoinExpress(express)
}

func (b *Builder[T]) LeftJoin(table string) ksql.JoinInterface {
	return b.query.LeftJoin(table)
}

func (b *Builder[T]) RightJoin(table string) ksql.JoinInterface {
	return b.query.RightJoin(table)
}

func (b *Builder[T]) ForUpdate() ksql.BuilderInterface[T] {
	b.query.For().Update()
	return b
}

func (b *Builder[T]) For() ksql.ForInterface {
	return b.query.For()
}

func (b *Builder[T]) All(ctx context.Context) error {
	return QueryBy(ctx, b._conn(), b.query, b.models)
}

func (b *Builder[T]) First(ctx context.Context) error {
	if b.conn == nil {
		return QueryRow(ctx, b.query, b.model)
	}

	return b.conn.QueryRow(ctx, b.query, b.model)
}

func (b *Builder[T]) _conn() ksql.ConnectionInterface {
	if b.conn != nil {
		return b.conn
	}

	return database
}

func _scanNum[T uint64 | float64](ctx context.Context, conn ksql.ConnectionInterface, query ksql.SqlInterface) (T, error) {
	cc := NewContext(ctx)
	cc.SqlLogStart(query)
	defer cc.SqlLogEnd()

	stmt, err := conn.Prepare(ctx, query)
	if err != nil {
		return 0, _err(err, query)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, query.Binds()...)
	if row.Err() != nil {
		return 0, _err(row.Err(), query)
	}

	var num *T
	if err := row.Scan(&num); err != nil {
		return 0, _err(err, query)
	}

	if num == nil {
		return 0, nil
	}

	return *num, nil
}

func (b *Builder[T]) SumFloat(ctx context.Context, column string) (float64, error) {
	b.Func("SUM", column, column)
	return _scanNum[float64](ctx, b._conn(), b.query)
}

func (b *Builder[T]) SumInt(ctx context.Context, column string) (uint64, error) {
	b.Func("SUM", column, column)
	return _scanNum[uint64](ctx, b._conn(), b.query)
}

func (b *Builder[T]) Count(ctx context.Context) (uint64, error) {
	b.ColumnsExpress(Raw("COUNT(1) as count"))
	return _scanNum[uint64](ctx, b._conn(), b.query)
}

func (b *Builder[T]) Exist(ctx context.Context) (bool, error) {
	cc := NewContext(ctx)
	cc.SqlLogStart(b.query)
	defer cc.SqlLogEnd()

	stmt, err := b._conn().Prepare(ctx, b.query)
	if err != nil {
		return false, _err(err, b.query)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, b.query.Binds()...)
	if err != nil {
		return false, _err(err, b.query)
	}
	defer rows.Close()
	if rows.Err() != nil {
		return false, _err(rows.Err(), b.query)
	}

	if rows.Next() {
		return true, nil
	}

	return false, _err(rows.Err(), b.query)
}

func offset(page, pageSize int64) int {
	return int((page - 1) * pageSize)
}

func (b *Builder[T]) Pagination(ctx context.Context, page, pageSize int64) (ksql.PaginationInterface[T], error) {
	ctx = NewContext(ctx)
	b.query.Limit(int(pageSize)).Offset(offset(page, pageSize))
	if err := b.All(ctx); err != nil {
		return nil, err
	}

	total := &Builder[T]{query: b.query.Clone(), conn: b.conn}
	count, err := total.Count(ctx)
	if err != nil {
		return nil, err
	}
	pageInfo := NewPageInfo(*b.models)
	pageInfo.Set(count, uint64(pageSize))
	return pageInfo, nil
}

func (b *Builder[T]) WithConn(conn ksql.ConnectionInterface) ksql.BuilderInterface[T] {
	b.conn = conn
	return b
}

func (b *Builder[T]) Max(ctx context.Context, column string) error {
	b.query.Func("MAX", column, column)
	return b.First(ctx)
}

func (b *Builder[T]) Min(ctx context.Context, column string) error {
	b.query.Func("MIN", column, column)
	return b.First(ctx)
}

func Build[T ksql.RowInterface](row T) ksql.BuilderInterface[T] {
	return NewBuilder(row)
}

func Rows[T ksql.RowInterface](rows *[]T) ksql.BuilderInterface[T] {
	builder := &Builder[T]{query: NewQuery(), models: rows}
	return builder
}

func Model[T ksql.ModelInterface](model T) ksql.BuilderInterface[T] {
	return NewBuilder(model).Table(model.Table()).Columns(model.Columns()...)
}

func Models[T ksql.ModelInterface](models *[]T) ksql.BuilderInterface[T] {
	var m T
	tmp := m.Clone().(T)
	builder := &Builder[T]{query: NewQuery(), models: models}
	builder.Table(tmp.Table()).Columns(tmp.Columns()...)
	return builder
}
