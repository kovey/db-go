package ksql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type IndexType byte

const (
	Index_Type_Normal   IndexType = 1
	Index_Type_Unique   IndexType = 2
	Index_Type_Primary  IndexType = 3
	Index_Type_FullText IndexType = 4
	Index_Type_Spatial  IndexType = 5
)

type Sharding byte

const (
	Sharding_None  Sharding = 0
	Sharding_Day   Sharding = 1
	Sharding_Month Sharding = 2
)

func FormatSharding(table string, sharding Sharding) string {
	switch sharding {
	case Sharding_Day:
		return fmt.Sprintf("%s_%s", table, time.Now().Format(Day_Format))
	case Sharding_Month:
		return fmt.Sprintf("%s_%s", table, time.Now().Format(Month_Format))
	default:
		return table
	}
}

const (
	Day_Format   = "20060102"
	Month_Format = "200601"
)

type ShardingInterface interface {
	Sharding(Sharding)
}

type TxError interface {
	error
	Begin() error
	Call() error
	Rollback() error
	Commit() error
}

type ConnectionInterface interface {
	Exec(ctx context.Context, op SqlInterface) (int64, error)
	QueryRow(ctx context.Context, op QueryInterface, model RowInterface) error
	Insert(ctx context.Context, op InsertInterface) (int64, error)
	Update(ctx context.Context, op UpdateInterface) (int64, error)
	Delete(ctx context.Context, op DeleteInterface) (int64, error)
	Database() *sql.DB
	Prepare(ctx context.Context, op SqlInterface) (*sql.Stmt, error)
	ExecRaw(ctx context.Context, raw ExpressInterface) (sql.Result, error)
	PrepareRaw(ctx context.Context, raw ExpressInterface) (*sql.Stmt, error)
	QueryRowRaw(ctx context.Context, raw ExpressInterface, model RowInterface) error
	DriverName() string
	InTransaction() bool
	Clone() ConnectionInterface
	Begin(ctx context.Context, options *sql.TxOptions) error
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
	Transaction(ctx context.Context, call func(ctx context.Context, conn ConnectionInterface) error) TxError
	TransactionBy(ctx context.Context, options *sql.TxOptions, call func(ctx context.Context, conn ConnectionInterface) error) TxError
	BeginTo(ctx context.Context, point string) error
	RollbackTo(ctx context.Context, point string) error
	CommitTo(ctx context.Context, point string) error
	ScanRaw(ctx context.Context, raw ExpressInterface, data ...any) error
}

type ExpressInterface interface {
	Statement() string
	Binds() []any
}

type RowInterface interface {
	Values() []any
	Clone() RowInterface
	WithConn(ConnectionInterface)
	Scan(s ScanInterface, r RowInterface) error
	Conn() ConnectionInterface
	Sharding(Sharding)
}

type ModelInterface interface {
	RowInterface
	Table() string
	Columns() []string
	PrimaryId() string
	Save(ctx context.Context) error
	Delete(ctx context.Context) error
	OnUpdateBefore(conn ConnectionInterface) error
	OnUpdateAfter(conn ConnectionInterface) error
	OnCreateBefore(conn ConnectionInterface) error
	OnCreateAfter(conn ConnectionInterface) error
	OnDeleteBefore(conn ConnectionInterface) error
	OnDeleteAfter(conn ConnectionInterface) error
	Empty() bool
}

type SqlInterface interface {
	Prepare() string
	Binds() []any
}

type WhereInterface interface {
	SqlInterface
	Where(column string, op string, data any) WhereInterface
	In(column string, data []any) WhereInterface
	NotIn(column string, data []any) WhereInterface
	IsNull(column string) WhereInterface
	IsNotNull(column string) WhereInterface
	Express(raw ExpressInterface) WhereInterface
	OrWhere(call func(o WhereInterface)) WhereInterface
	InBy(column string, sub QueryInterface) WhereInterface
	NotInBy(column string, sub QueryInterface) WhereInterface
	Between(column string, begin, end any) WhereInterface
	AndWhere(call func(o WhereInterface)) WhereInterface
	Empty() bool
}

type HavingInterface interface {
	SqlInterface
	Having(column string, op string, data any) HavingInterface
	In(column string, data []any) HavingInterface
	NotIn(column string, data []any) HavingInterface
	IsNull(column string) HavingInterface
	IsNotNull(column string) HavingInterface
	Express(raw ExpressInterface) HavingInterface
	OrHaving(call func(o HavingInterface)) HavingInterface
	InBy(column string, sub QueryInterface) HavingInterface
	NotInBy(column string, sub QueryInterface) HavingInterface
	Between(column string, begin, end any) HavingInterface
	AndHaving(call func(o HavingInterface)) HavingInterface
	Empty() bool
}

type JoinInterface interface {
	SqlInterface
	Table(table string) JoinInterface
	As(as string) JoinInterface
	On(column, op, val string) JoinInterface
	OnOr(call func(JoinOnInterface)) JoinInterface
	Express(express ExpressInterface) JoinInterface
	Type() string
}

type JoinOnInterface interface {
	SqlInterface
	On(column, op, val string) JoinInterface
}

type InsertInterface interface {
	SqlInterface
	Add(column string, data any) InsertInterface
	Table(table string) InsertInterface
}

type UpdateInterface interface {
	SqlInterface
	Set(column string, data any) UpdateInterface
	Where(WhereInterface) UpdateInterface
	Table(table string) UpdateInterface
}

type ColumnInterface interface {
	Express() string
	Nullable() ColumnInterface
	AutoIncrement() ColumnInterface
	Unsigned() ColumnInterface
	Default(value string) ColumnInterface
	Comment(comment string) ColumnInterface
	UseCurrent() ColumnInterface
	UseCurrentOnUpdate() ColumnInterface
}

type AddColumnInterface interface {
	AddColumn(column, t string, length, scale int, sets ...string) ColumnInterface
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

type AlterInterface interface {
	SqlInterface
	AddColumnInterface
	DropColumn(column string) AlterInterface
	AddIndex(name string, t IndexType, column ...string) AlterInterface
	DropIndex(name string) AlterInterface
	Table(table string) AlterInterface
	ChangeColumn(oldColumn, newColumn, t string, length, scale int, sets ...string) ColumnInterface
	Comment(comment string) AlterInterface
	AddPrimary(column string) AlterInterface
	AddUnique(name string, columns ...string) AlterInterface
	Charset(charset string) AlterInterface
	Collate(collate string) AlterInterface
	Engine(engine string) AlterInterface
}

type DeleteInterface interface {
	SqlInterface
	Where(WhereInterface) DeleteInterface
	Table(table string) DeleteInterface
}

type QueryInterface interface {
	SqlInterface
	Sharding(sharding Sharding)
	GetSharding() Sharding
	Table(table string) QueryInterface
	TableBy(query QueryInterface, as string) QueryInterface
	As(as string) QueryInterface
	Column(column, as string) QueryInterface
	Func(fun, column, as string) QueryInterface
	Columns(columns ...string) QueryInterface
	ColumnsExpress(expresses ...ExpressInterface) QueryInterface
	Where(column string, op string, val any) QueryInterface
	WhereExpress(expresses ...ExpressInterface) QueryInterface
	OrWhere(callback func(WhereInterface)) QueryInterface
	WhereIsNull(column string) QueryInterface
	WhereIsNotNull(column string) QueryInterface
	WhereIn(column string, data []any) QueryInterface
	WhereNotIn(column string, data []any) QueryInterface
	WhereInBy(column string, query QueryInterface) QueryInterface
	WhereNotInBy(column string, query QueryInterface) QueryInterface
	AndWhere(call func(w WhereInterface)) QueryInterface
	Between(column string, begin, end any) QueryInterface
	Having(column string, op string, val any) QueryInterface
	HavingExpress(expresses ...ExpressInterface) QueryInterface
	OrHaving(callback func(HavingInterface)) QueryInterface
	HavingIsNull(column string) QueryInterface
	HavingIsNotNull(column string) QueryInterface
	HavingIn(column string, data []any) QueryInterface
	HavingNotIn(column string, data []any) QueryInterface
	HavingInBy(column string, query QueryInterface) QueryInterface
	HavingNotInBy(column string, query QueryInterface) QueryInterface
	HavingBetween(column string, begin, end any) QueryInterface
	AndHaving(call func(h HavingInterface)) QueryInterface
	Limit(limit int) QueryInterface
	Offset(offset int) QueryInterface
	Order(column string) QueryInterface
	OrderDesc(column string) QueryInterface
	Group(column ...string) QueryInterface
	Join(table string) JoinInterface
	JoinExpress(express ExpressInterface) JoinInterface
	LeftJoin(table string) JoinInterface
	LeftJoinExpress(express ExpressInterface) JoinInterface
	RightJoin(table string) JoinInterface
	RightJoinExpress(express ExpressInterface) JoinInterface
	ForUpdate() QueryInterface
	Clone() QueryInterface
	Pagination(page, pageSize int) QueryInterface
	Distinct(column string) QueryInterface
	FuncDistinct(fun, column, as string) QueryInterface
}

type CreateTableInterface interface {
	SqlInterface
	AddColumnInterface
	AddIndex(name string, t IndexType, column ...string) CreateTableInterface
	Table(table string) CreateTableInterface
	Engine(engine string) CreateTableInterface
	Charset(charset string) CreateTableInterface
	Collate(collate string) CreateTableInterface
	Comment(comment string) CreateTableInterface
	AddPrimary(column string) CreateTableInterface
	AddUnique(name string, columns ...string) CreateTableInterface
}

type DropTableInterface interface {
	SqlInterface
	Table(table string) DropTableInterface
}

type SchemaInterface interface {
	SqlInterface
	Create() SchemaInterface
	Alter() SchemaInterface
	Schema(schema string) SchemaInterface
	Charset(charset string) SchemaInterface
	Collate(collate string) SchemaInterface
}

type PaginationInterface[T RowInterface] interface {
	List() []T
	TotalPage() uint64
	TotalCount() uint64
	Set(totalCount, pageSize uint64)
}

type ScanInterface interface {
	Scan(...any) error
}

type EngineInterface interface {
	Format(SqlInterface) string
	FormatRaw(ExpressInterface) string
}

type TraceInterface interface {
	TraceId() string
}

type ContextInterface interface {
	context.Context
	SqlLogStart(sql SqlInterface)
	RawSqlLogStart(sql ExpressInterface)
	SqlLogEnd()
}
