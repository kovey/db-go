package schema

import "github.com/kovey/db-go/v3"

type SchemaInfoInterface interface {
	Name() string
	Charset() string
	Collation() string
	HasChanged(other SchemaInfoInterface) bool
	HasTable(table TableInfoInterface) bool
	Tables() []TableInfoInterface
}

type ColumnInfoInterface interface {
	Name() string
	Default() string
	Nullable() bool
	Type() string
	Length() int
	NumLen() int
	Scale() int
	DateTimeLen() int
	Comment() string
	Extra() string
	HasChanged(other ColumnInfoInterface) bool
	HasDefault() bool
	AutoIncrement() bool
	Key() string
}

type TableInfoInterface interface {
	Name() string
	Charset() string
	Collation() string
	Comment() string
	Engine() string
	Fields() []ColumnInfoInterface
	HasChanged(other TableInfoInterface) bool
	HasColumn(column ColumnInfoInterface) bool
	Indexes() []IndexInfoInterface
	HasIndex(index IndexInfoInterface) bool
	GetIndex(index string) IndexInfoInterface
	CheckChanges(other TableInfoInterface) ChangedInterface
	GetColumn(column string) ColumnInfoInterface
}

type IndexInfoInterface interface {
	Name() string
	Metas() []IndexMetaInterface
	Add(index IndexMetaInterface) IndexInfoInterface
	Type() ksql.IndexType
	HasChanged(other IndexInfoInterface) bool
	Columns() []string
}

type IndexMetaInterface interface {
	Name() string
	NonUnique() int
	Seq() int
	Column() string
	Comment() string
	IndexComment() string
	Type() ksql.IndexType
	HasChanged(other IndexMetaInterface) bool
}

type IndexChangedInterface interface {
	Adds() []IndexInfoInterface
	Deletes() []IndexInfoInterface
}

type ColumnMetaChangedInterface interface {
	Old() ColumnInfoInterface
	New() ColumnInfoInterface
}

type ColumnChangedInterface interface {
	Adds() []ColumnInfoInterface
	Changes() []ColumnMetaChangedInterface
	Deletes() []ColumnInfoInterface
}

type ChangedInterface interface {
	Column() ColumnChangedInterface
	Index() IndexChangedInterface
}
