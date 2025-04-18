package ksql

import "strings"

type TableIndexInterface interface {
	Type(typ IndexType) TableIndexInterface
	SubType(subType IndexSubType) TableIndexInterface
	Algorithm(algorithm IndexAlg) TableIndexInterface
	Column(name string, length int, order Order) TableIndexInterface
	Express(express string, order Order) TableIndexInterface
	Columns(columns ...string) TableIndexInterface
	BlockSize(size string) TableIndexInterface
	WithParser(parserName string) TableIndexInterface
	Comment(comment string) TableIndexInterface
	Visible() TableIndexInterface
	Invisible() TableIndexInterface
	EngineAttribute(attr string) TableIndexInterface
	SecondaryEngineAttribute(attr string) TableIndexInterface
	Constraint(symbol string) TableIndexInterface
	Primary() TableIndexInterface
	Unique() TableIndexInterface
	Foreign() TableIndexInterface
	Reference(table string) ColumnReferenceInterface
	Build(builder *strings.Builder)
}

type TableOptionsInterface interface {
	Append(key TableOptKey, value string) TableOptionsInterface
	AppendWithPrefix(prefix TableOptPrefix, key TableOptKey, value string) TableOptionsInterface
	AppendWith(key TableOptKey, value TableOptVal) TableOptionsInterface
	InsertMethod(value InsertMethod) TableOptionsInterface
	StartTransation() TableOptionsInterface
	RowFormat(value RowFormat) TableOptionsInterface
	Build(builder *strings.Builder)
	Union(tables ...string) TableOptionsInterface
	SpaceOption(name string, storage ColumnStorage) TableOptionsInterface
}

type PartitionOptionsKeyInterface interface {
	Algorithm(alg string) PartitionOptionsKeyInterface
	Linear() PartitionOptionsKeyInterface
	Build(builder *strings.Builder)
}

type PartitionOptionsRangeInterface interface {
	Expr(expr string) PartitionOptionsRangeInterface
	Columns(columns []string) PartitionOptionsRangeInterface
	Build(builder *strings.Builder)
}

type PartitionOptionsSubInterface interface {
	Hash(expr string) PartitionOptionsSubInterface
	Key(columns []string) PartitionOptionsKeyInterface
	Build(builder *strings.Builder)
}

type PartitionDefinitionLessthanInterface interface {
	Expr(expr string) PartitionDefinitionLessthanInterface
	ValueList(valueList []string) PartitionDefinitionLessthanInterface
	Build(builder *strings.Builder)
	MaxValue() PartitionDefinitionLessthanInterface
}

type PartitionDefinitionInInterface interface {
	Build(builder *strings.Builder)
}

type PartitionDefinitionValuesInterface interface {
	LessThan() PartitionDefinitionLessthanInterface
	In(valueList []string) PartitionDefinitionInInterface
	Build(builder *strings.Builder)
}

type PartitionDefinitionSubInterface interface {
	Build(builder *strings.Builder)
	Option(key PartitionDefinitionOptKey, value string) PartitionDefinitionSubInterface
}

type PartitionDefinitionInterface interface {
	Values() PartitionDefinitionValuesInterface
	Option(key PartitionDefinitionOptKey, value string) PartitionDefinitionInterface
	Sub(name string) PartitionDefinitionSubInterface
	Build(builder *strings.Builder)
}

type PartitionOptionsHashInterface interface {
	Linear() PartitionOptionsHashInterface
	Build(builder *strings.Builder)
}

type PartitionOptionsInterface interface {
	Hash(expr string) PartitionOptionsHashInterface
	Key(columns []string) PartitionOptionsKeyInterface
	Range() PartitionOptionsRangeInterface
	List() PartitionOptionsRangeInterface
	Sub() PartitionOptionsSubInterface
	Definition(name string) PartitionDefinitionInterface
	Build(builder *strings.Builder)
}

type AlterOptAlg string

const (
	Alter_Opt_Alg_Default AlterOptAlg = "DEFAULT"
	Alter_Opt_Alg_Instant AlterOptAlg = "INSTANT"
	Alter_Opt_Alg_Inplace AlterOptAlg = "INPLACE"
	Alter_Opt_Alg_Copy    AlterOptAlg = "COPY"
)

type BuildInterface interface {
	Build(builder *strings.Builder)
}

type AlterColumnInterface interface {
	BuildInterface
	Column(column string) AlterColumnInterface
	Default(value string) AlterColumnInterface
	DefaultExpress(expr string) AlterColumnInterface
	Visible() AlterColumnInterface
	Invisible() AlterColumnInterface
	DropDefault() AlterColumnInterface
}

type AlterIndexInterface interface {
	BuildInterface
	Index(index string) AlterIndexInterface
	Visible() AlterIndexInterface
	Invisible() AlterIndexInterface
}

type ChangeColumnInterface interface {
	BuildInterface
	Old(column string) ChangeColumnInterface
	New(column, t string, length, scale int, sets ...string) ColumnInterface
	First() ChangeColumnInterface
	After(column string) ChangeColumnInterface
}

type ModifyColumnInterface interface {
	BuildInterface
	Column(t string, length, scale int, sets ...string) ColumnInterface
	First() ModifyColumnInterface
	After(column string) ModifyColumnInterface
}

type TableAddColumnInterface interface {
	BuildInterface
	Column(column, t string, length, scale int, sets ...string) ColumnInterface
	First() TableAddColumnInterface
	After(column string) TableAddColumnInterface
}

type TableMultiAddColumnInterface interface {
	BuildInterface
	Column(column, t string, length, scale int, sets ...string) ColumnInterface
	First() TableAddColumnInterface
	After(column string) TableAddColumnInterface
}

type RenameColumnInterface interface {
	BuildInterface
	Old(column string) RenameColumnInterface
	New(column string) RenameColumnInterface
}

type RenameIndexInterface interface {
	BuildInterface
	Old(index string) RenameIndexInterface
	New(index string) RenameIndexInterface
	Type(typ IndexSubType) RenameIndexInterface
}

type TablespaceInterface interface {
	SqlInterface
	Tablespace(tablespace string) TablespaceInterface
	Option(key, value string) TablespaceInterface
	OptionStr(key, value string) TablespaceInterface
	OptionWith(prefix, key, value string) TablespaceInterface
	OptionStrWith(prefix, key, value string) TablespaceInterface
	OptionOnlyKey(key string) TablespaceInterface
	Alter() TablespaceInterface
	Undo() TablespaceInterface
}

type TriggerOrderType string

const (
	Trigger_Order_Type_Follows  TriggerOrderType = "FOLLOWS"
	Trigger_Order_Type_Precedes TriggerOrderType = "PRECEDES"
)

type TriggerInterface interface {
	SqlInterface
	Trigger(trigger string) TriggerInterface
	Definer(definer string) TriggerInterface
	IfNotExists() TriggerInterface
	Before() TriggerInterface
	After() TriggerInterface
	Insert() TriggerInterface
	Update() TriggerInterface
	Delete() TriggerInterface
	On(table string) TriggerInterface
	Order(typ TriggerOrderType, otherTrigger string) TriggerInterface
	Body(sql SqlInterface) TriggerInterface
	BodyRaw(sql ExpressInterface) TriggerInterface
}

type ViewAlg string

const (
	View_Alg_Undefined ViewAlg = "UNDEFINED"
	View_Alg_Merge     ViewAlg = "MERGE"
	View_Alg_Tempable  ViewAlg = "TEMPTABLE"
)

type ViewInterface interface {
	SqlInterface
	Replace() ViewInterface
	Algorithm(alg ViewAlg) ViewInterface
	Definer(definer string) ViewInterface
	SqlSecurity(security SqlSecurity) ViewInterface
	View(name string) ViewInterface
	Columns(columns ...string) ViewInterface
	As(query QueryInterface) ViewInterface
	WithCascaded() ViewInterface
	WithLocal() ViewInterface
	CheckOption() ViewInterface
	Alter() ViewInterface
}
