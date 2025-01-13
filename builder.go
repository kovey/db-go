package ksql

import "context"

type BuilderInterface[T RowInterface] interface {
	Sharding(sharding Sharding) BuilderInterface[T]
	Table(table string) BuilderInterface[T]
	TableBy(query QueryInterface, as string) BuilderInterface[T]
	As(as string) BuilderInterface[T]
	Column(column, as string) BuilderInterface[T]
	Func(fun, column, as string) BuilderInterface[T]
	Columns(columns ...string) BuilderInterface[T]
	ColumnsExpress(expresses ...ExpressInterface) BuilderInterface[T]
	Where(column string, op string, val any) BuilderInterface[T]
	WhereExpress(expresses ...ExpressInterface) BuilderInterface[T]
	OrWhere(callback func(WhereInterface)) BuilderInterface[T]
	WhereIsNull(column string) BuilderInterface[T]
	WhereIsNotNull(column string) BuilderInterface[T]
	WhereIn(column string, data []any) BuilderInterface[T]
	WhereNotIn(column string, data []any) BuilderInterface[T]
	WhereInBy(column string, query QueryInterface) BuilderInterface[T]
	WhereNotInBy(column string, query QueryInterface) BuilderInterface[T]
	AndWhere(call func(WhereInterface)) BuilderInterface[T]
	Between(column string, begin, end any) BuilderInterface[T]
	Having(column string, op string, val any) BuilderInterface[T]
	HavingExpress(...ExpressInterface) BuilderInterface[T]
	HavingIsNull(column string) BuilderInterface[T]
	HavingIsNotNull(column string) BuilderInterface[T]
	HavingIn(column string, data []any) BuilderInterface[T]
	HavingNotIn(column string, data []any) BuilderInterface[T]
	HavingInBy(column string, query QueryInterface) BuilderInterface[T]
	HavingNotInBy(column string, query QueryInterface) BuilderInterface[T]
	HavingBetween(column string, begin, end any) BuilderInterface[T]
	AndHaving(call func(HavingInterface)) BuilderInterface[T]
	OrHaving(func(HavingInterface)) BuilderInterface[T]
	Limit(limit int) BuilderInterface[T]
	Offset(offset int) BuilderInterface[T]
	Order(column string) BuilderInterface[T]
	OrderDesc(column string) BuilderInterface[T]
	Group(column ...string) BuilderInterface[T]
	Join(table string) JoinInterface
	JoinExpress(ExpressInterface) JoinInterface
	LeftJoin(table string) JoinInterface
	LeftJoinExpress(ExpressInterface) JoinInterface
	RightJoin(table string) JoinInterface
	RightJoinExpress(ExpressInterface) JoinInterface
	ForUpdate() BuilderInterface[T]
	All(context.Context) error
	First(context.Context) error
	Max(ctx context.Context, column string) error
	Min(ctx context.Context, column string) error
	Exist(ctx context.Context) (bool, error)
	Count(ctx context.Context) (uint64, error)
	SumInt(ctx context.Context, column string) (uint64, error)
	SumFloat(ctx context.Context, column string) (float64, error)
	Pagination(ctx context.Context, page, pageSize int64) (PaginationInterface[T], error)
	Distinct(column string) BuilderInterface[T]
	FuncDistinct(fun, column, as string) BuilderInterface[T]
	WithConn(conn ConnectionInterface) BuilderInterface[T]
}

type TableInterface interface {
	AddColumnInterface
	DropColumn(column string) TableInterface
	AddIndex(name string, t IndexType, column ...string) TableInterface
	DropIndex(name string) TableInterface
	Table(table string) TableInterface
	ChangeColumn(oldColumn, newColumn, t string, length, scale int, sets ...string) ColumnInterface
	Exec(ctx context.Context) error
	Create() TableInterface
	Alter() TableInterface
	Engine(engine string) TableInterface
	Charset(charset string) TableInterface
	Collate(collate string) TableInterface
	Comment(comment string) TableInterface
	AddPrimary(column string) TableInterface
	AddUnique(name string, columns ...string) TableInterface
	HasColumn(ctx context.Context, column string) (bool, error)
	HasIndex(ctx context.Context, index string) (bool, error)
	WithConn(conn ConnectionInterface) TableInterface
}
