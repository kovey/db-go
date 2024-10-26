package ksql

import "context"

type BuilderInterface[T RowInterface] interface {
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
	Between(column string, begin, end any) BuilderInterface[T]
	Having(column string, op string, val any) BuilderInterface[T]
	HavingExpress(...ExpressInterface) BuilderInterface[T]
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
	All(context.Context, *[]T) error
	First(context.Context, T) error
	Exist(ctx context.Context) (bool, error)
}

type TableInterface interface {
	AddColumn(column, t string, length, scale int, sets ...string) ColumnInterface
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
	SetConn(conn ConnectionInterface) TableInterface
	AddDecimal(column string, length, scale int) ColumnInterface
	AddDouble(column string, length, scale int) ColumnInterface
	AddFloat(column string, length, scale int) ColumnInterface
	AddBinary(column string, length int) ColumnInterface
	AddGeoMetry(column string) ColumnInterface
	AddPolygon(column string) ColumnInterface
	AddPoint(column string) ColumnInterface
	AddLineString(column string) ColumnInterface
	AddBlob(column string) ColumnInterface
	AddText(column string) ColumnInterface
	AddSet(column string, sets []string) ColumnInterface
	AddEnum(column string, options []string) ColumnInterface
	AddDate(column string) ColumnInterface
	AddDateTime(column string) ColumnInterface
	AddTimestamp(column string) ColumnInterface
	AddSmallInt(column string) ColumnInterface
	AddTinyInt(column string) ColumnInterface
	AddBigInt(column string) ColumnInterface
	AddInt(column string) ColumnInterface
	AddString(column string, length int) ColumnInterface
	AddChar(column string, length int) ColumnInterface
}
