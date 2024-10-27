package db

import (
	"context"

	ksql "github.com/kovey/db-go/v3"
)

type Builder[T ksql.RowInterface] struct {
	query ksql.QueryInterface
	conn  ksql.ConnectionInterface
}

func NewBuilder[T ksql.RowInterface](model T) *Builder[T] {
	return &Builder[T]{query: NewQuery(), conn: model.Conn()}
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

func (b *Builder[T]) Where(column string, op string, val any) ksql.BuilderInterface[T] {
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

func (b *Builder[T]) Between(column string, begin, end any) ksql.BuilderInterface[T] {
	b.query.Between(column, begin, end)
	return b
}

func (b *Builder[T]) Having(column string, op string, val any) ksql.BuilderInterface[T] {
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

func (b *Builder[T]) LeftJoinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	return b.query.LeftJoinExpress(express)
}

func (b *Builder[T]) RightJoin(table string) ksql.JoinInterface {
	return b.query.RightJoin(table)
}

func (b *Builder[T]) RightJoinExpress(express ksql.ExpressInterface) ksql.JoinInterface {
	return b.query.RightJoinExpress(express)
}

func (b *Builder[T]) ForUpdate() ksql.BuilderInterface[T] {
	b.query.ForUpdate()
	return b
}

func (b *Builder[T]) All(ctx context.Context, models *[]T) error {
	return QueryBy(ctx, b._conn(), b.query, models)
}

func (b *Builder[T]) First(ctx context.Context, model T) error {
	if b.conn == nil {
		return QueryRow(ctx, b.query, model)
	}

	return b.conn.QueryRow(ctx, b.query, model)
}

func (b *Builder[T]) _conn() ksql.ConnectionInterface {
	if b.conn != nil {
		return b.conn
	}

	return database
}

func (b *Builder[T]) SumFloat(ctx context.Context, column string) (float64, error) {
	b.Func("SUM", column, column)
	stmt, err := b._conn().Prepare(ctx, b.query)
	if err != nil {
		return 0, err
	}

	row := stmt.QueryRowContext(ctx, b.query.Binds()...)
	if row.Err() != nil {
		return 0, err
	}

	var sum *float64
	if err := row.Scan(&sum); err != nil {
		return 0, err
	}

	if sum == nil {
		return 0, nil
	}

	return *sum, nil
}

func (b *Builder[T]) SumInt(ctx context.Context, column string) (uint64, error) {
	b.Func("SUM", column, column)
	stmt, err := b._conn().Prepare(ctx, b.query)
	if err != nil {
		return 0, err
	}

	row := stmt.QueryRowContext(ctx, b.query.Binds()...)
	if row.Err() != nil {
		return 0, err
	}

	var sum *uint64
	if err := row.Scan(&sum); err != nil {
		return 0, err
	}

	if sum == nil {
		return 0, nil
	}

	return *sum, nil
}

func (b *Builder[T]) Count(ctx context.Context) (uint64, error) {
	b.ColumnsExpress(Raw("COUNT(1) as count"))
	stmt, err := b._conn().Prepare(ctx, b.query)
	if err != nil {
		return 0, err
	}

	row := stmt.QueryRowContext(ctx, b.query.Binds()...)
	if row.Err() != nil {
		return 0, err
	}

	var count *uint64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	if count == nil {
		return 0, nil
	}

	return *count, nil
}

func (b *Builder[T]) Exist(ctx context.Context) (bool, error) {
	stmt, err := b._conn().Prepare(ctx, b.query)
	if err != nil {
		return false, err
	}

	rows, err := stmt.QueryContext(ctx, b.query.Binds()...)
	if err != nil {
		return false, _err(err, b.query)
	}
	if rows.Err() != nil {
		return false, _err(rows.Err(), b.query)
	}

	defer rows.Close()
	if rows.Next() {
		return true, nil
	}

	return false, _err(rows.Err(), b.query)
}

func offset(page, pageSize int64) int {
	return int((page - 1) * pageSize)
}

func (b *Builder[T]) Pagination(ctx context.Context, page, pageSize int64) (ksql.PaginationInterface[T], error) {
	var models []T
	b.query.Limit(int(pageSize)).Offset(offset(page, pageSize))
	if err := b.All(ctx, &models); err != nil {
		return nil, err
	}

	total := &Builder[T]{query: b.query.Clone(), conn: b.conn}
	count, err := total.Count(ctx)
	if err != nil {
		return nil, err
	}
	pageInfo := NewPageInfo(models)
	pageInfo.Set(count, uint64(pageSize))
	return pageInfo, nil
}

func Build[T ksql.RowInterface](row T) ksql.BuilderInterface[T] {
	return NewBuilder(row)
}
