package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
	"github.com/kovey/db-go/v3/sql/table"
)

type Alter struct {
	*base
	options           *table.AlterOptions
	tableOptions      *table.Options
	partitionOperates *table.PartitionOperates
	table             string
	orders            []string
}

func NewAlter() *Alter {
	u := &Alter{base: newBase(), options: &table.AlterOptions{}, tableOptions: table.NewOptions(), partitionOperates: &table.PartitionOperates{}}
	u.opChain.Append(u._keyword, u._tableOptions, u._options, u._partitionOptions)
	return u
}

func (u *Alter) _keyword(builder *strings.Builder) {
	builder.WriteString("ALTER TABLE")
	operator.BuildColumnString(u.table, builder)
}

func (u *Alter) _options(builder *strings.Builder) {
	if u.options.Empty() {
		return
	}

	if u.tableOptions.Empty() {
		builder.WriteString(" ")
	} else {
		builder.WriteString(", ")
	}

	u.options.Build(builder)
}

func (u *Alter) _tableOptions(builder *strings.Builder) {
	if u.tableOptions.Empty() {
		return
	}

	builder.WriteString(" ")
	u.tableOptions.Build(builder)
}

func (u *Alter) _partitionOptions(builder *strings.Builder) {
	if u.partitionOperates.Empty() {
		return
	}

	u.partitionOperates.Build(builder)
}

func (u *Alter) AlterColumn(column string) ksql.AlterColumnInterface {
	ac := table.NewAlterColumn()
	u.options.Append(ac)
	return ac
}

func (u *Alter) Table(table string) ksql.AlterInterface {
	u.table = table
	return u
}

func (a *Alter) AddCheck(expr string) ksql.ColumnCheckConstraintInterface {
	cc := table.NewColumnCheckConstraint().Check(expr)
	a.options.Append(table.NewAddAlterOption(cc))
	return cc
}

func (a *Alter) DropCheck(symbol string) ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("DROP", "CHECK", symbol))
	return a
}

func (a *Alter) DropConstraint(symbol string) ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("DROP", "CONSTRAINT", symbol))
	return a
}

func (a *Alter) AlterCheck(symbol string) ksql.ColumnCheckConstraintInterface {
	cc := table.NewColumnCheckConstraint().Constraint(symbol)
	a.options.Append(table.NewEditAlterOption(cc))
	return cc
}

func (a *Alter) Algorithm(alg ksql.AlterOptAlg) ksql.AlterInterface {
	a.options.Append(table.NewEqAlterOption("ALGORITHM", string(alg)))
	return a
}

func (a *Alter) ModifyColumn(column string) ksql.ModifyColumnInterface {
	m := table.NewModifyColumn(column)
	a.options.Append(m)
	return m
}

func (a *Alter) ChangeColumn(oldColumn string) ksql.ChangeColumnInterface {
	c := table.NewChangeColumn().Old(oldColumn)
	a.options.Append(c)
	return c
}

func (a *Alter) AddIndex(name string) ksql.TableIndexInterface {
	add := table.NewAddIndex(name)
	a.options.Append(add)
	return add.Index()
}

func (a *Alter) AlterIndex(index string) ksql.AlterIndexInterface {
	al := table.NewAlterIndex().Index(index)
	a.options.Append(al)
	return al
}

func (a *Alter) DropIndex(name string) ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("DROP", "INDEX", name).IsField())
	return a
}

func (a *Alter) DropKey(name string) ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("DROP", "KEY", name).IsField())
	return a
}

func (a *Alter) DropPrimary() ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("DROP", "PRIMARY", "KEY"))
	return a
}

func (a *Alter) DropForeign(name string) ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("DROP", "FOREIGN KEY", name).IsField())
	return a
}

func (a *Alter) Force() ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("FORCE", "", ""))
	return a
}

func (a *Alter) Lock(lock ksql.IndexLockOption) ksql.AlterInterface {
	a.options.Append(table.NewEqAlterOption("LOCK", string(lock)))
	return a
}

func (a *Alter) OrderBy(columns ...string) ksql.AlterInterface {
	a.orders = append(a.orders, columns...)
	return a
}

func (a *Alter) AddDecimal(column string, length, scale int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Decimal, length, scale)
}

func (a *Alter) AddDouble(column string, length, scale int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Double, length, scale)
}

func (a *Alter) AddFloat(column string, length, scale int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Float, length, scale)
}

func (a *Alter) AddBinary(column string, length int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Binary, length, 0)
}

func (a *Alter) AddGeoMetry(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_GeoMetry, 0, 0)
}

func (a *Alter) AddPolygon(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Polygon, 0, 0)
}

func (a *Alter) AddPoint(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Point, 0, 0)
}

func (a *Alter) AddLineString(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_LineString, 0, 0)
}

func (a *Alter) AddBlob(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Blob, 0, 0)
}

func (a *Alter) AddText(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Text, 0, 0)
}

func (a *Alter) AddSet(column string, sets []string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Set, 0, 0, sets...)
}

func (a *Alter) AddEnum(column string, options []string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Enum, 0, 0, options...)
}

func (a *Alter) AddDate(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Date, 0, 0)
}

func (a *Alter) AddDateTime(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_DateTime, 0, 0)
}

func (a *Alter) AddTimestamp(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Timestamp, 0, 0)
}

func (a *Alter) AddSmallInt(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_SmallInt, 3, 0)
}
func (a *Alter) AddTinyInt(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_TinyInt, 1, 0)
}

func (a *Alter) AddBigInt(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_BigInt, 20, 0)
}

func (a *Alter) AddInt(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Int, 11, 0)
}

func (a *Alter) AddString(column string, length int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_VarChar, length, 0)
}

func (a *Alter) AddChar(column string, length int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Char, length, 0)
}

func (u *Alter) AddColumn(column, t string, length, scale int, sets ...string) ksql.ColumnInterface {
	add := table.NewAddColumn()
	col := add.Column(column, t, length, scale, sets...)
	if col != nil {
		u.options.Append(add)
	}
	return col
}

func (u *Alter) AddColumnBy(column, t string, length, scale int, sets ...string) ksql.TableAddColumnInterface {
	add := table.NewAddColumn()
	if add.Column(column, t, length, scale, sets...) != nil {
		u.options.Append(add)
	}
	return add
}

func (u *Alter) DropColumn(column string) ksql.AlterInterface {
	d := table.NewAlterOption("DROP", "COLUMN", column).IsField()
	u.options.Append(d)
	return u
}

func (u *Alter) DropColumnIfExists(column string) ksql.AlterInterface {
	return u.DropColumn(column)
}

func (u *Alter) Comment(comment string) ksql.AlterInterface {
	u.tableOptions.Append(ksql.Table_Opt_Key_Comment, comment)
	return u
}

func (u *Alter) Charset(charset string) ksql.AlterInterface {
	u.tableOptions.Append(ksql.Table_Opt_Key_Character_Set, charset)
	return u
}

func (u *Alter) Collate(collate string) ksql.AlterInterface {
	u.tableOptions.Append(ksql.Table_Opt_Key_Collate, collate)
	return u
}

func (u *Alter) Engine(engine string) ksql.AlterInterface {
	u.tableOptions.Append(ksql.Table_Opt_Key_Engine, engine)
	return u
}

func (a *Alter) AddForeign(name string, columns ...string) ksql.TableIndexInterface {
	add := table.NewAddIndex(name)
	add.Index().Foreign().Columns(columns...)
	a.options.Append(add)
	return add.Index()
}

func (a *Alter) AddUnique(name string, columns ...string) ksql.TableIndexInterface {
	add := table.NewAddIndex(name)
	add.Index().Unique().Columns(columns...)
	a.options.Append(add)
	return add.Index()
}

func (a *Alter) AddPrimary(columns ...string) ksql.AlterInterface {
	add := table.NewAddIndex("")
	add.Index().Primary().Columns(columns...)
	a.options.Append(add)
	return a
}

func (a *Alter) Default(charset, collate string) ksql.AlterInterface {
	d := &table.DefaultAlterOption{}
	if charset != "" {
		d.Option("CHARACTER SET", charset)
	}

	if collate != "" {
		d.Option("COLLATE", collate)
	}

	if !d.Empty() {
		a.options.Append(d)
	}
	return a
}

func (a *Alter) Rename(as string) ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("RENAME", "AS", as).IsField())
	return a
}

func (a *Alter) RenameTo(to string) ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("RENAME", "TO", to).IsField())
	return a
}

func (a *Alter) RenameColumn(oldName, newName string) ksql.AlterInterface {
	rc := &table.RenameColumn{}
	rc.Old(oldName).New(newName)
	a.options.Append(rc)
	return a
}

func (a *Alter) RenameKey(oldName, newName string) ksql.AlterInterface {
	rc := &table.RenameIndex{}
	rc.Old(oldName).New(newName).Type(ksql.Index_Sub_Type_Key)
	a.options.Append(rc)
	return a
}

func (a *Alter) RenameIndex(oldName, newName string) ksql.AlterInterface {
	rc := &table.RenameIndex{}
	rc.Old(oldName).New(newName).Type(ksql.Index_Sub_Type_Index)
	a.options.Append(rc)
	return a
}

func (a *Alter) WithValidation() ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("WITH", "VALIDATION", ""))
	return a
}

func (a *Alter) WithoutValidation() ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("WITHOUT", "VALIDATION", ""))
	return a
}

func (a *Alter) DisableKeys() ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("DISABLE", "KEYS", ""))
	return a
}

func (a *Alter) EnableKeys() ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("ENABLE", "KEYS", ""))
	return a
}

func (a *Alter) DiscardTablespace() ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("DISCARD", "TABLESPACE", ""))
	return a
}

func (a *Alter) ImportTablespace() ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("IMPORT", "TABLESPACE", ""))
	return a
}

func (a *Alter) ConvertToCharset(charset string, collate string) ksql.AlterInterface {
	a.options.Append(table.NewAlterOption("CONVERT TO", "CHARACTER SET", charset))
	if collate != "" {
		a.options.Append(table.NewAlterOption("COLLATE", collate, ""))
	}
	return a
}

func (a *Alter) Options() ksql.TableOptionsInterface {
	return a.tableOptions
}

func (a *Alter) AddPartition(name string) ksql.PartitionDefinitionInterface {
	op := table.NewPartitionOperate(name)
	a.partitionOperates.Append(op)
	return op.Add()
}

func (a *Alter) CoalescePartition(number int) ksql.AlterInterface {
	op := table.NewPartitionOperate()
	op.Coalesce(number)
	a.partitionOperates.Append(op)
	return a
}

func (a *Alter) PartitionOperate(partitionNames ...string) ksql.PartitionOperateInterface {
	op := table.NewPartitionOperate(partitionNames...)
	a.partitionOperates.Append(op)
	return op
}
