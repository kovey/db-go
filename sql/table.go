package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/table"
)

type Table struct {
	*base
	table            string
	ifNotExists      bool
	columns          []*table.Column
	indexes          []*table.Index
	options          *table.Options
	partitionOptions *table.PartitionOptions
	asOption         string
	as               ksql.QueryInterface
	likeTable        string
	temporary        string
}

func NewTable() *Table {
	ta := &Table{base: newBase(), options: table.NewOptions()}
	ta.opChain.Append(keywordCreate, ta._table, ta._like, ta._columns, ta._options, ta._partOptions, ta._as)
	return ta
}

func (ta *Table) _table(builder *strings.Builder) {
	if ta.temporary != "" {
		builder.WriteString(" ")
		builder.WriteString(ta.temporary)
	}

	builder.WriteString(" TABLE")
	if ta.ifNotExists {
		builder.WriteString(" IF NOT EXISTS")
	}

	builder.WriteString(" ")
	Column(ta.table, builder)
}

func (ta *Table) _like(builder *strings.Builder) {
	if ta.likeTable == "" {
		return
	}

	builder.WriteString(" (LIKE ")
	Column(ta.likeTable, builder)
	builder.WriteString(")")
}

func (ta *Table) _columns(builder *strings.Builder) {
	if ta.likeTable != "" || len(ta.columns) == 0 {
		return
	}

	builder.WriteString(" (")
	index := 0
	for _, column := range ta.columns {
		if index > 0 {
			builder.WriteString(", ")
		}

		column.Build(builder)
		index++
	}

	for _, i := range ta.indexes {
		if index > 0 {
			builder.WriteString(",")
		}

		i.Build(builder)
	}

	builder.WriteString(")")
}

func (ta *Table) _options(builder *strings.Builder) {
	if ta.likeTable != "" || ta.options.Empty() {
		return
	}

	builder.WriteString(" ")
	ta.options.Build(builder)
}

func (ta *Table) _partOptions(builder *strings.Builder) {
	if ta.likeTable != "" || ta.partitionOptions == nil {
		return
	}

	ta.partitionOptions.Build(builder)
}

func (ta *Table) _as(builder *strings.Builder) {
	if ta.likeTable != "" || ta.as == nil {
		return
	}

	if ta.asOption != "" {
		builder.WriteString(" ")
		builder.WriteString(ta.asOption)
	}

	builder.WriteString(" AS ")
	builder.WriteString(ta.as.Prepare())
	ta.binds = append(ta.binds, ta.as.Binds()...)
}

func (ta *Table) IfNotExists() ksql.CreateTableInterface {
	ta.ifNotExists = true
	return ta
}

func (ta *Table) Temporary() ksql.CreateTableInterface {
	ta.temporary = "TEMPORARY"
	return ta
}

func (ta *Table) Like(table string) ksql.CreateTableInterface {
	ta.likeTable = table
	return ta
}

func (ta *Table) As(query ksql.QueryInterface) ksql.CreateTableInterface {
	ta.as = query
	return ta
}

func (ta *Table) AddDecimal(column string, length, scale int) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Decimal, length, scale)
}

func (ta *Table) AddDouble(column string, length, scale int) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Double, length, scale)
}

func (ta *Table) AddFloat(column string, length, scale int) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Float, length, scale)
}

func (ta *Table) AddBinary(column string, length int) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Binary, length, 0)
}

func (ta *Table) AddGeoMetry(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_GeoMetry, 0, 0)
}

func (ta *Table) AddPolygon(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Polygon, 0, 0)
}

func (ta *Table) AddPoint(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Point, 0, 0)
}

func (ta *Table) AddLineString(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_LineString, 0, 0)
}

func (ta *Table) AddBlob(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Blob, 0, 0)
}

func (ta *Table) AddText(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Text, 0, 0)
}

func (ta *Table) AddSet(column string, sets []string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Set, 0, 0, sets...)
}

func (ta *Table) AddEnum(column string, options []string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Enum, 0, 0, options...)
}

func (ta *Table) AddDate(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Date, 0, 0)
}

func (ta *Table) AddDateTime(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_DateTime, 19, 0)
}

func (ta *Table) AddTimestamp(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Timestamp, 19, 0)
}

func (ta *Table) AddSmallInt(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_SmallInt, 3, 0)
}
func (ta *Table) AddTinyInt(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_TinyInt, 1, 0)
}

func (ta *Table) AddBigInt(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_BigInt, 20, 0)
}

func (ta *Table) AddInt(column string) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Int, 11, 0)
}

func (ta *Table) AddString(column string, length int) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_VarChar, length, 0)
}

func (ta *Table) AddChar(column string, length int) ksql.ColumnInterface {
	return ta.AddColumn(column, table.Type_Char, length, 0)
}

func (ta *Table) AddColumn(column, t string, length, scale int, sets ...string) ksql.ColumnInterface {
	ct := table.ParseType(t, length, scale, sets...)
	if ct == nil {
		return nil
	}

	col := table.NewColumn(column, ct)
	ta.columns = append(ta.columns, col)
	return col
}

func (ta *Table) AddIndex(name string) ksql.TableIndexInterface {
	index := table.NewIndex(name)
	ta.indexes = append(ta.indexes, index)
	return index
}

func (ta *Table) AddPrimary(column string) ksql.TableIndexInterface {
	return ta.AddIndex("").Primary().Columns(column)
}

func (ta *Table) AddUnique(name string, columns ...string) ksql.TableIndexInterface {
	return ta.AddIndex(name).Unique().Columns(columns...)
}

func (ta *Table) Table(table string) ksql.CreateTableInterface {
	ta.table = table
	return ta
}

func (ta *Table) Options() ksql.TableOptionsInterface {
	return ta.options
}

func (ta *Table) PartitionOptions() ksql.PartitionOptionsInterface {
	if ta.partitionOptions == nil {
		ta.partitionOptions = &table.PartitionOptions{}
	}

	return ta.partitionOptions
}

func (ta *Table) Engine(engine string) ksql.CreateTableInterface {
	ta.options.Append(ksql.Table_Opt_Key_Engine, engine)
	return ta
}

func (ta *Table) Charset(charset string) ksql.CreateTableInterface {
	ta.options.Append(ksql.Table_Opt_Key_Character_Set, charset)
	return ta
}

func (ta *Table) Collate(collate string) ksql.CreateTableInterface {
	ta.options.Append(ksql.Table_Opt_Key_Collate, collate)
	return ta
}

func (ta *Table) Comment(comment string) ksql.CreateTableInterface {
	ta.options.Append(ksql.Table_Opt_Key_Comment, comment)
	return ta
}
