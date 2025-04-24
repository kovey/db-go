package ksql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type IndexType byte

const (
	Index_Type_Normal   IndexType = 0x1
	Index_Type_Unique   IndexType = 0x2
	Index_Type_Primary  IndexType = 0x3
	Index_Type_FullText IndexType = 0x4
	Index_Type_Spatial  IndexType = 0x5
	Index_Type_Foreign  IndexType = 0x6
)

func (i IndexType) String() string {
	switch i {
	case Index_Type_FullText:
		return "FULLTEXT"
	case Index_Type_Spatial:
		return "SPATIAL"
	case Index_Type_Unique:
		return "UNIQUE"
	default:
		return ""
	}
}

type IndexSubType string

const (
	Index_Sub_Type_Key   IndexSubType = "KEY"
	Index_Sub_Type_Index IndexSubType = "INDEX"
)

type IndexAlg string

const (
	Index_Alg_BTree IndexAlg = "BTREE"
	Index_Alg_Hash  IndexAlg = "HASH"
)

type IndexAlgOption string

const (
	Index_Alg_Option_Default IndexAlgOption = "DEFAULT"
	Index_Alg_Option_Inplace IndexAlgOption = "INPLACE"
	Index_Alg_Option_Copy    IndexAlgOption = "COPY"
)

// DEFAULT | NONE | SHARED | EXCLUSIVE
type IndexLockOption string

const (
	Index_Lock_Option_Default   IndexLockOption = "DEFAULT"
	Index_Lock_Option_None      IndexLockOption = "NONE"
	Index_Lock_Option_Shared    IndexLockOption = "SHARED"
	Index_Lock_Option_Exclusive IndexLockOption = "EXCLUSIVE"
)

type SqlType string

const (
	Sql_Type_Insert     SqlType = "INSERT"
	Sql_Type_Update     SqlType = "UPDATE"
	Sql_Type_Delete     SqlType = "DELETE"
	Sql_Type_Drop       SqlType = "DROP"
	Sql_Type_Alter      SqlType = "ALTER"
	Sql_Type_Create     SqlType = "CREATE"
	Sql_Type_Query      SqlType = "QUERY"
	Sql_Type_Save_Point SqlType = "SAVEPOINT"
	Sql_Type_Release    SqlType = "RELEASE"
	Sql_Type_Rollback   SqlType = "ROLLBACK"
)

type Op string

const (
	Eq   Op = "="
	Le   Op = "<="
	Lt   Op = "<"
	Ge   Op = ">="
	Gt   Op = ">"
	Like Op = "LIKE"
	Neq  Op = "<>"
)

var ops = map[Op]byte{
	Eq:   1,
	Le:   1,
	Lt:   1,
	Ge:   1,
	Gt:   1,
	Like: 1,
	Neq:  1,
}

func SupportOp(op Op) bool {
	_, ok := ops[Op(strings.ToUpper(string(op)))]
	return ok
}

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
	Scan(ctx context.Context, query QueryInterface, data ...any) error
}

type ExpressInterface interface {
	Statement() string
	Binds() []any
	IsExec() bool
	Type() SqlType
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
	BuildInterface
	Where(column string, op Op, data any) WhereInterface
	In(column string, data []any) WhereInterface
	NotIn(column string, data []any) WhereInterface
	IsNull(column string) WhereInterface
	IsNotNull(column string) WhereInterface
	Express(raw ExpressInterface) WhereInterface
	OrWhere(call func(o WhereInterface)) WhereInterface
	InBy(column string, sub QueryInterface) WhereInterface
	NotInBy(column string, sub QueryInterface) WhereInterface
	Between(column string, begin, end any) WhereInterface
	NotBetween(column string, begin, end any) WhereInterface
	AndWhere(call func(o WhereInterface)) WhereInterface
	Empty() bool
	Binds() []any
	Clone() WhereInterface
}

type HavingInterface interface {
	BuildInterface
	Having(column string, op Op, data any) HavingInterface
	In(column string, data []any) HavingInterface
	NotIn(column string, data []any) HavingInterface
	IsNull(column string) HavingInterface
	IsNotNull(column string) HavingInterface
	Express(raw ExpressInterface) HavingInterface
	OrHaving(call func(o HavingInterface)) HavingInterface
	InBy(column string, sub QueryInterface) HavingInterface
	NotInBy(column string, sub QueryInterface) HavingInterface
	Between(column string, begin, end any) HavingInterface
	NotBetween(column string, begin, end any) HavingInterface
	AndHaving(call func(o HavingInterface)) HavingInterface
	Empty() bool
	Binds() []any
	Clone() HavingInterface
}

type JoinInterface interface {
	BuildInterface
	Table(table string) JoinInterface
	As(as string) JoinInterface
	On(column, op, val string) JoinInterface
	OnOr(call func(JoinOnInterface)) JoinInterface
	Express(express ExpressInterface) JoinInterface
	Left() JoinInterface
	Right() JoinInterface
	Inner() JoinInterface
	Binds() []any
}

type JoinOnInterface interface {
	BuildInterface
	On(column, op, val string) JoinOnInterface
	OnVal(column, op string, val any) JoinOnInterface
}

type InsertInterface interface {
	SqlInterface
	Add(column string, data any) InsertInterface
	Table(table string) InsertInterface
	From(query QueryInterface) InsertInterface
	FromTable(table string) InsertInterface
	Columns(columns ...string) InsertInterface
	Values(datas ...any) InsertInterface
	Set(column, value string) InsertInterface
	SetColumn(column, otherColumn string) InsertInterface
	SetExpress(expr ExpressInterface) InsertInterface
	LowPriority() InsertInterface
	Delayed() InsertInterface
	HighPriority() InsertInterface
	Ignore() InsertInterface
	Partitions(names ...string) InsertInterface
	As(rowAlias string, colAlias ...string) InsertInterface
	OnDuplicateKeyUpdate(column, value string) InsertInterface
	OnDuplicateKeyUpdateColumn(column, otherColumn string) InsertInterface
	OnDuplicateKeyUpdateExpress(expr ExpressInterface) InsertInterface
}

type UpdateInterface interface {
	SqlInterface
	Set(column string, data any) UpdateInterface
	Where(WhereInterface) UpdateInterface
	Table(table string) UpdateInterface
	LowPriority() UpdateInterface
	Ignore() UpdateInterface
	OrderByAsc(columns ...string) UpdateInterface
	OrderByDesc(columns ...string) UpdateInterface
	SetExpress(expre ExpressInterface) UpdateInterface
	SetColumn(column string, otherColumn string) UpdateInterface
	Limit(limit int) UpdateInterface
	IncColumn(column string, data int) UpdateInterface
}

type UpdateMultiInterface interface {
	SqlInterface
	Set(column string, data any) UpdateMultiInterface
	Where(WhereInterface) UpdateMultiInterface
	Table(table string) UpdateMultiInterface
	SetExpress(expre ExpressInterface) UpdateMultiInterface
	SetColumn(column string, otherColumn string) UpdateMultiInterface
	LowPriority() UpdateMultiInterface
	Ignore() UpdateMultiInterface
	Join(table string) JoinInterface
	JoinExpress(express ExpressInterface) JoinInterface
	LeftJoin(table string) JoinInterface
	RightJoin(table string) JoinInterface
}

type ColumnFormat string

const (
	Column_Format_Fixed   ColumnFormat = "FIXED"
	Column_Format_Dynamic ColumnFormat = "DYNAMIC"
	Column_Format_Default ColumnFormat = "DEFAULT"
)

type ColumnStorage string

const (
	Column_Storage_Disk   ColumnStorage = "DISK"
	Column_Storage_Memory ColumnStorage = "MEMORY"
)

type ReferenceMatch string

const (
	Reference_Match_Full   ReferenceMatch = "FULL"
	Reference_Match_Patial ReferenceMatch = "PARTIAL"
	Reference_Match_Simple ReferenceMatch = "SIMPLE"
)

type ReferenceOnOpt string

const (
	Reference_On_Opt_DELETE ReferenceOnOpt = "DELETE"
	Reference_On_Opt_UPDATE ReferenceOnOpt = "UPDATE"
)

type ReferenceOption string

const (
	Reference_Option_Restrict    ReferenceOption = "RESTRICT"
	Reference_Option_Cascade     ReferenceOption = "CASCADE"
	Reference_Option_Set_Null    ReferenceOption = "SET NULL"
	Reference_Option_No_Action   ReferenceOption = "NO ACTION"
	Reference_Option_Set_Default ReferenceOption = "SET DEFAULT"
)

type ColumnReferenceInterface interface {
	Column(name string, length int, order Order) ColumnReferenceInterface
	Express(express string, order Order) ColumnReferenceInterface
	Match(match ReferenceMatch) ColumnReferenceInterface
	On(op ReferenceOnOpt, option ReferenceOption) ColumnReferenceInterface
	Build(builder *strings.Builder)
}

type ColumnCheckConstraintInterface interface {
	Constraint(symbol string) ColumnCheckConstraintInterface
	Check(expr string) ColumnCheckConstraintInterface
	Enforced() ColumnCheckConstraintInterface
	NotEnforced() ColumnCheckConstraintInterface
	Build(builder *strings.Builder)
}

type ColumnInterface interface {
	Build(builder *strings.Builder)
	Nullable() ColumnInterface
	NotNullable() ColumnInterface
	Default(value string) ColumnInterface
	DefaultExpress(expr string) ColumnInterface
	DefaultBit(value string) ColumnInterface
	Visible() ColumnInterface
	Invisible() ColumnInterface
	AutoIncrement() ColumnInterface
	Unique() ColumnInterface
	Primary() ColumnInterface
	Index() ColumnInterface
	Comment(comment string) ColumnInterface
	Collate(collate string) ColumnInterface
	Format(format ColumnFormat) ColumnInterface
	EngineAttribute(engineAttr string) ColumnInterface
	SecondaryEngineAttribute(secEngineAttr string) ColumnInterface
	Storage(storage ColumnStorage) ColumnInterface
	Reference(table string) ColumnReferenceInterface
	CheckConstraint() ColumnCheckConstraintInterface
	Unsigned() ColumnInterface
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

type PartitionOperateInterface interface {
	BuildInterface
	Drop() PartitionOperateInterface
	Discard() PartitionOperateInterface
	Import() PartitionOperateInterface
	Truncate() PartitionOperateInterface
	Reorganize() PartitionDefinitionInterface
	Exchange(table, validation string) PartitionOperateInterface
	Analyze() PartitionOperateInterface
	Check() PartitionOperateInterface
	Optimize() PartitionOperateInterface
	Rebuild() PartitionOperateInterface
	Repair() PartitionOperateInterface
	Remove() PartitionOperateInterface
}

type AlterInterface interface {
	SqlInterface
	AddColumnInterface
	AlterColumn(column string) AlterColumnInterface
	DropColumn(column string) AlterInterface
	DropColumnIfExists(column string) AlterInterface
	AddIndex(name string) TableIndexInterface
	AlterIndex(index string) AlterIndexInterface
	AddCheck(expr string) ColumnCheckConstraintInterface
	DropCheck(symbol string) AlterInterface
	DropConstraint(symbol string) AlterInterface
	AlterCheck(symbol string) ColumnCheckConstraintInterface
	Algorithm(alg AlterOptAlg) AlterInterface
	DropIndex(name string) AlterInterface
	DropKey(name string) AlterInterface
	DropPrimary() AlterInterface
	DropForeign(name string) AlterInterface
	Force() AlterInterface
	Lock(lock IndexLockOption) AlterInterface
	OrderBy(columns ...string) AlterInterface
	Table(table string) AlterInterface
	ChangeColumn(oldColumn string) ChangeColumnInterface
	ModifyColumn(column string) ModifyColumnInterface
	Comment(comment string) AlterInterface
	AddPrimary(columns ...string) AlterInterface
	AddUnique(name string, columns ...string) TableIndexInterface
	AddForeign(name string, columns ...string) TableIndexInterface
	Default(charset, collate string) AlterInterface
	Charset(charset string) AlterInterface
	Collate(collate string) AlterInterface
	Engine(engine string) AlterInterface
	Rename(as string) AlterInterface
	RenameTo(to string) AlterInterface
	RenameColumn(oldName, newName string) AlterInterface
	RenameKey(oldName, newName string) AlterInterface
	RenameIndex(oldName, newName string) AlterInterface
	WithValidation() AlterInterface
	WithoutValidation() AlterInterface
	DisableKeys() AlterInterface
	EnableKeys() AlterInterface
	DiscardTablespace() AlterInterface
	ImportTablespace() AlterInterface
	ConvertToCharset(charset string, collate string) AlterInterface
	Options() TableOptionsInterface
	AddPartition(name string) PartitionDefinitionInterface
	CoalescePartition(number int) AlterInterface
	PartitionOperate(partitionNames ...string) PartitionOperateInterface
	AddColumnBy(column, t string, length, scale int, sets ...string) TableAddColumnInterface
}

type DeleteInterface interface {
	SqlInterface
	LowPriority() DeleteInterface
	Quick() DeleteInterface
	Ignore() DeleteInterface
	Table(table string) DeleteInterface
	As(as string) DeleteInterface
	Partitions(names ...string) DeleteInterface
	Where(WhereInterface) DeleteInterface
	OrderByDesc(columns ...string) DeleteInterface
	OrderByAsc(columns ...string) DeleteInterface
	Limit(limit int) DeleteInterface
}

type DeleteMultiInterface interface {
	SqlInterface
	LowPriority() DeleteMultiInterface
	Quick() DeleteMultiInterface
	Ignore() DeleteMultiInterface
	Table(table string) DeleteMultiInterface
	Join(table string) JoinInterface
	JoinExpress(express ExpressInterface) JoinInterface
	LeftJoin(table string) JoinInterface
	RightJoin(table string) JoinInterface
	Where(WhereInterface) DeleteMultiInterface
}

type ForInterface interface {
	BuildInterface
	Update() ForInterface
	Share() ForInterface
	Of(tables ...string) ForInterface
	NoWait() ForInterface
	SkipLocked() ForInterface
	LockInShareMode() ForInterface
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
	Where(column string, op Op, val any) QueryInterface
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
	NotBetween(column string, begin, end any) QueryInterface
	Having(column string, op Op, val any) QueryInterface
	HavingExpress(expresses ...ExpressInterface) QueryInterface
	OrHaving(callback func(HavingInterface)) QueryInterface
	HavingIsNull(column string) QueryInterface
	HavingIsNotNull(column string) QueryInterface
	HavingIn(column string, data []any) QueryInterface
	HavingNotIn(column string, data []any) QueryInterface
	HavingInBy(column string, query QueryInterface) QueryInterface
	HavingNotInBy(column string, query QueryInterface) QueryInterface
	HavingBetween(column string, begin, end any) QueryInterface
	HavingNotBetween(column string, begin, end any) QueryInterface
	AndHaving(call func(h HavingInterface)) QueryInterface
	Limit(limit int) QueryInterface
	Offset(offset int) QueryInterface
	Order(column ...string) QueryInterface
	OrderDesc(column ...string) QueryInterface
	Group(column ...string) QueryInterface
	Join(table string) JoinInterface
	JoinExpress(express ExpressInterface) JoinInterface
	LeftJoin(table string) JoinInterface
	RightJoin(table string) JoinInterface
	Clone() QueryInterface
	Pagination(page, pageSize int) QueryInterface
	Distinct() QueryInterface
	FuncDistinct(fun, column, as string) QueryInterface
	IntoVar(vars ...string) QueryInterface
	All() QueryInterface
	DistinctRow() QueryInterface
	HighPriority() QueryInterface
	StraightJoin() QueryInterface
	SqlSmallResult() QueryInterface
	SqlBigResult() QueryInterface
	SqlBufferResult() QueryInterface
	SqlNoCache() QueryInterface
	SqlCalcFoundRows() QueryInterface
	Partitions(names ...string) QueryInterface
	GroupWithRollUp() QueryInterface
	Window(window, as string) QueryInterface
	OrderWithRollUp() QueryInterface
	For() ForInterface
	WhereInCall(column string, call func(query QueryInterface)) QueryInterface
	WhereNotInCall(column string, call func(query QueryInterface)) QueryInterface
}

type CreateTableInterface interface {
	SqlInterface
	AddColumnInterface
	AddIndex(name string) TableIndexInterface
	Table(table string) CreateTableInterface
	Engine(engine string) CreateTableInterface
	Charset(charset string) CreateTableInterface
	Collate(collate string) CreateTableInterface
	Comment(comment string) CreateTableInterface
	AddPrimary(column string) TableIndexInterface
	AddUnique(name string, columns ...string) TableIndexInterface
	As(QueryInterface) CreateTableInterface
	Like(table string) CreateTableInterface
	Temporary() CreateTableInterface
	Options() TableOptionsInterface
	PartitionOptions() PartitionOptionsInterface
	IfNotExists() CreateTableInterface
}

type DropTableInterface interface {
	SqlInterface
	Table(table string) DropTableInterface
	IfExists() DropTableInterface
	Temporary() DropTableInterface
	Restrict() DropTableInterface
	Cascade() DropTableInterface
}

type DropTablespaceInterface interface {
	SqlInterface
	Tablespace(table string) DropTablespaceInterface
	Undo() DropTablespaceInterface
	Engine(engine string) DropTablespaceInterface
}

type DropTriggerInterface interface {
	SqlInterface
	Trigger(trigger string) DropTriggerInterface
	IfExists() DropTriggerInterface
	Schema(schema string) DropTriggerInterface
}

type DropViewInterface interface {
	SqlInterface
	View(view string) DropViewInterface
	IfExists() DropViewInterface
	Restrict() DropViewInterface
	Cascade() DropViewInterface
}

type RenameTableInterface interface {
	SqlInterface
	Table(from, to string) RenameTableInterface
}

type TruncateTableInterface interface {
	SqlInterface
	Table(table string) TruncateTableInterface
}

type DropSchemaInterface interface {
	SqlInterface
	Schema(schema string) DropSchemaInterface
	IfExists() DropSchemaInterface
}

type DropEventInterface interface {
	SqlInterface
	Event(event string) DropEventInterface
	IfExists() DropEventInterface
}

type DropFunctionInterface interface {
	SqlInterface
	Function(function string) DropFunctionInterface
	IfExists() DropFunctionInterface
}

type DropProcedureInterface interface {
	SqlInterface
	Procedure(procedure string) DropProcedureInterface
	IfExists() DropProcedureInterface
}

type DropIndexInterface interface {
	SqlInterface
	Index(index string) DropIndexInterface
	Table(table string) DropIndexInterface
	Algorithm(alg IndexAlgOption) DropIndexInterface
	Lock(lock IndexLockOption) DropIndexInterface
}

type DropLogFileGroupInterface interface {
	SqlInterface
	LogFileGroup(logFileGroup string) DropLogFileGroupInterface
	Engine(engine string) DropLogFileGroupInterface
}

type DropServerInterface interface {
	SqlInterface
	Server(server string) DropServerInterface
	IfExists() DropServerInterface
}

type DropSpatialReferenceSystemInterface interface {
	SqlInterface
	Srid(srid string) DropSpatialReferenceSystemInterface
	IfExists() DropSpatialReferenceSystemInterface
}

type SchemaInterface interface {
	SqlInterface
	Create() SchemaInterface
	Alter() SchemaInterface
	Schema(schema string) SchemaInterface
	Charset(charset string) SchemaInterface
	Collate(collate string) SchemaInterface
	Encryption(encryption string) SchemaInterface
	ReadOnly(readOnly string) SchemaInterface
	IfNotExists() SchemaInterface
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
	WithTraceId(traceId string) ContextInterface
}

type IntervalUnit string

const (
	Interval_Year          IntervalUnit = "YEAR"
	Interval_Quarter       IntervalUnit = "QUARTER"
	Interval_Month         IntervalUnit = "MONTH"
	Interval_Day           IntervalUnit = "DAY"
	Interval_Hour          IntervalUnit = "HOUR"
	Interval_Minute        IntervalUnit = "MINUTE"
	Interval_Week          IntervalUnit = "WEEK"
	Interval_Second        IntervalUnit = "SECOND"
	Interval_Year_Month    IntervalUnit = "YEAR_MONTH"
	Interval_Day_Hour      IntervalUnit = "DAY_HOUR"
	Interval_Day_Minute    IntervalUnit = "DAY_MINUTE"
	Interval_Day_Second    IntervalUnit = "DAY_SECOND"
	Interval_Hour_Minute   IntervalUnit = "HOUR_MINUTE"
	Interval_Hour_Second   IntervalUnit = "HOUR_SECOND"
	Interval_Minute_Second IntervalUnit = "MINUTE_SECOND"
)

func (i IntervalUnit) IsNumber() bool {
	switch i {
	case Interval_Year, Interval_Quarter, Interval_Month, Interval_Day, Interval_Hour, Interval_Minute, Interval_Week, Interval_Second:
		return true
	default:
		return false
	}
}

type EventStatus string

const (
	Event_Status_Enable           EventStatus = "ENABLE"
	Event_Status_Disable          EventStatus = "DISABLE"
	Event_Status_Disable_On_Slave EventStatus = "DISABLE ON SLAVE"
)

type EventInterface interface {
	SqlInterface
	Definer(name string) EventInterface
	IfNotExists() EventInterface
	Event(name string) EventInterface
	Comment(comment string) EventInterface
	Do(sql SqlInterface) EventInterface
	DoRaw(sql ExpressInterface) EventInterface
	At(timestamp string) EventInterface
	AtInterval(interval string, unit IntervalUnit) EventInterface
	Every(interval string, unit IntervalUnit) EventInterface
	Starts(timestamp string) EventInterface
	StartsInterval(interval string, unit IntervalUnit) EventInterface
	Ends(timestamp string) EventInterface
	EndsInterval(interval string, unit IntervalUnit) EventInterface
	Status(status EventStatus) EventInterface
	OnCompletion() EventInterface
	OnCompletionNot() EventInterface
	Rename(name string) EventInterface
	Alter() EventInterface
}

type ProcedureSqlType string

const (
	Procedure_Sql_Type_Contains_Sql      ProcedureSqlType = "CONTAINS SQL"
	Procedure_Sql_Type_No_Sql            ProcedureSqlType = "NO SQL"
	Procedure_Sql_Type_Redis_Sql_Data    ProcedureSqlType = "READS SQL DATA"
	Procedure_Sql_Type_Modifies_Sql_Data ProcedureSqlType = "MODIFIES SQL DATA"
)

type SqlSecurity string

const (
	Sql_Security_Definer SqlSecurity = "DEFINER"
	Sql_Security_Invoker SqlSecurity = "INVOKER"
)

type ProcedureInterface interface {
	SqlInterface
	Definer(name string) ProcedureInterface
	IfNotExists() ProcedureInterface
	Procedure(name string) ProcedureInterface
	In(name, typ string) ProcedureInterface
	Out(name, typ string) ProcedureInterface
	InOut(name, typ string) ProcedureInterface
	Comment(comment string) ProcedureInterface
	Language() ProcedureInterface
	Deterministic() ProcedureInterface
	DeterministicNot() ProcedureInterface
	SqlType(sqlType ProcedureSqlType) ProcedureInterface
	SqlSecurity(security SqlSecurity) ProcedureInterface
	RoutineBody(sql ExpressInterface) ProcedureInterface
	Alter() ProcedureInterface
}

type FunctionInterface interface {
	SqlInterface
	Definer(name string) FunctionInterface
	IfNotExists() FunctionInterface
	Function(name string) FunctionInterface
	Param(name, typ string) FunctionInterface
	Returns(typ string) FunctionInterface
	Comment(comment string) FunctionInterface
	Language() FunctionInterface
	Deterministic() FunctionInterface
	DeterministicNot() FunctionInterface
	SqlType(sqlType ProcedureSqlType) FunctionInterface
	SqlSecurity(security SqlSecurity) FunctionInterface
	RoutineBody(raw ExpressInterface) FunctionInterface
	Alter() FunctionInterface
}

type Order string

const (
	Order_None Order = "NONE"
	Order_Asc  Order = "ASC"
	Order_Desc Order = "DESC"
)

type IndexInterface interface {
	SqlInterface
	Type(typ IndexType) IndexInterface
	Index(name string) IndexInterface
	Algorithm(alg IndexAlg) IndexInterface
	On(table string) IndexInterface
	Column(name string, length int, order Order) IndexInterface
	Express(express string, order Order) IndexInterface
	BlockSize(size string) IndexInterface
	WithParser(parserName string) IndexInterface
	Comment(comment string) IndexInterface
	Visible() IndexInterface
	Invisible() IndexInterface
	EngineAttribute(attr string) IndexInterface
	SecondaryEngineAttribute(attr string) IndexInterface
	AlgorithmOption(option IndexAlgOption) IndexInterface
	LockOption(option IndexLockOption) IndexInterface
}

type LogFileGroupInterface interface {
	SqlInterface
	LogFileGroup(name string) LogFileGroupInterface
	UndoFile(file string) LogFileGroupInterface
	InitialSize(size string) LogFileGroupInterface
	UndoBufferSize(size string) LogFileGroupInterface
	RedoBufferSize(size string) LogFileGroupInterface
	NodeGroupId(nodegroupId string) LogFileGroupInterface
	Wait() LogFileGroupInterface
	Comment(comment string) LogFileGroupInterface
	Engine(engine string) LogFileGroupInterface
	Alter() LogFileGroupInterface
}

type ServOptKey string

const (
	Serv_Opt_Key_Host     ServOptKey = "HOST"
	Serv_Opt_Key_Database ServOptKey = "DATABASE"
	Serv_Opt_Key_User     ServOptKey = "USER"
	Serv_Opt_Key_Password ServOptKey = "PASSWORD"
	Serv_Opt_Key_Socket   ServOptKey = "SOCKET"
	Serv_Opt_Key_Owner    ServOptKey = "OWNER"
	Serv_Opt_Key_Port     ServOptKey = "PORT"
)

type ServerInterface interface {
	SqlInterface
	Server(name string) ServerInterface
	WrapperName(wrapperName string) ServerInterface
	Option(key ServOptKey, val string) ServerInterface
	Alter() ServerInterface
}

type SrsAttr string

const (
	Srs_Attr_Name        SrsAttr = "NAME"
	Srs_Attr_Definition  SrsAttr = "DEFINITION"
	Srs_Attr_Description SrsAttr = "DESCRIPTION"
)

type SpatialReferenceSystemInterface interface {
	SqlInterface
	Replace() SpatialReferenceSystemInterface
	IfNotExists() SpatialReferenceSystemInterface
	Srid(srid uint32) SpatialReferenceSystemInterface
	Atrribute(key SrsAttr, value string) SpatialReferenceSystemInterface
	Organization(value string, identified uint32) SpatialReferenceSystemInterface
}
