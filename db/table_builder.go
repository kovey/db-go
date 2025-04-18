package db

import (
	"context"

	ksql "github.com/kovey/db-go/v3"
)

type TableBuilder struct {
	createMode bool
	alterMode  bool
	alter      ksql.AlterInterface
	create     ksql.CreateTableInterface
	table      string
	conn       ksql.ConnectionInterface
}

func NewTableBuilder() *TableBuilder {
	return &TableBuilder{}
}

func (ta *TableBuilder) WithConn(conn ksql.ConnectionInterface) ksql.TableInterface {
	ta.conn = conn
	return ta
}

func (ta *TableBuilder) Create() ksql.TableInterface {
	if ta.alterMode {
		return ta
	}

	ta.createMode = true
	ta.create = NewCreateTable()
	if ta.table != "" {
		ta.create.Table(ta.table)
	}
	return ta
}

func (ta *TableBuilder) Alter() ksql.TableInterface {
	if ta.createMode {
		return ta
	}

	ta.alterMode = true
	ta.alter = NewAlterTable()
	if ta.table != "" {
		ta.alter.Table(ta.table)
	}
	return ta
}

func (ta *TableBuilder) AddDecimal(column string, length, scale int) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddDecimal(column, length, scale)
	}

	if ta.createMode {
		return ta.create.AddDecimal(column, length, scale)
	}

	return nil
}

func (ta *TableBuilder) AddDouble(column string, length, scale int) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddDouble(column, length, scale)
	}

	if ta.createMode {
		return ta.create.AddDouble(column, length, scale)
	}

	return nil
}

func (ta *TableBuilder) AddFloat(column string, length, scale int) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddFloat(column, length, scale)
	}

	if ta.createMode {
		return ta.create.AddFloat(column, length, scale)
	}

	return nil
}

func (ta *TableBuilder) AddBinary(column string, length int) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddBinary(column, length)
	}

	if ta.createMode {
		return ta.create.AddBinary(column, length)
	}

	return nil
}

func (ta *TableBuilder) AddGeoMetry(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddGeoMetry(column)
	}

	if ta.createMode {
		return ta.create.AddGeoMetry(column)
	}

	return nil
}

func (ta *TableBuilder) AddPolygon(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddPolygon(column)
	}

	if ta.createMode {
		return ta.create.AddPolygon(column)
	}

	return nil
}

func (ta *TableBuilder) AddPoint(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddPoint(column)
	}

	if ta.createMode {
		return ta.create.AddPoint(column)
	}

	return nil
}

func (ta *TableBuilder) AddLineString(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddLineString(column)
	}

	if ta.createMode {
		return ta.create.AddLineString(column)
	}

	return nil
}

func (ta *TableBuilder) AddBlob(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddBlob(column)
	}

	if ta.createMode {
		return ta.create.AddBlob(column)
	}

	return nil
}

func (ta *TableBuilder) AddText(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddText(column)
	}

	if ta.createMode {
		return ta.create.AddText(column)
	}

	return nil
}

func (ta *TableBuilder) AddSet(column string, sets []string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddSet(column, sets)
	}

	if ta.createMode {
		return ta.create.AddSet(column, sets)
	}

	return nil
}

func (ta *TableBuilder) AddEnum(column string, options []string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddEnum(column, options)
	}

	if ta.createMode {
		return ta.create.AddEnum(column, options)
	}

	return nil
}

func (ta *TableBuilder) AddDate(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddDate(column)
	}

	if ta.createMode {
		return ta.create.AddDate(column)
	}

	return nil
}

func (ta *TableBuilder) AddDateTime(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddDateTime(column)
	}

	if ta.createMode {
		return ta.create.AddDateTime(column)
	}

	return nil
}

func (ta *TableBuilder) AddTimestamp(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddTimestamp(column)
	}

	if ta.createMode {
		return ta.create.AddTimestamp(column)
	}

	return nil
}

func (ta *TableBuilder) AddSmallInt(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddSmallInt(column)
	}

	if ta.createMode {
		return ta.create.AddSmallInt(column)
	}

	return nil
}
func (ta *TableBuilder) AddTinyInt(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddTinyInt(column)
	}

	if ta.createMode {
		return ta.create.AddTinyInt(column)
	}

	return nil
}

func (ta *TableBuilder) AddBigInt(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddBigInt(column)
	}

	if ta.createMode {
		return ta.create.AddBigInt(column)
	}

	return nil
}

func (ta *TableBuilder) AddInt(column string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddInt(column)
	}

	if ta.createMode {
		return ta.create.AddInt(column)
	}

	return nil
}

func (ta *TableBuilder) AddString(column string, length int) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddString(column, length)
	}

	if ta.createMode {
		return ta.create.AddString(column, length)
	}

	return nil
}

func (ta *TableBuilder) AddChar(column string, length int) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddChar(column, length)
	}

	if ta.createMode {
		return ta.create.AddChar(column, length)
	}

	return nil
}

func (ta *TableBuilder) AddColumn(column, t string, length, scale int, sets ...string) ksql.ColumnInterface {
	if ta.alterMode {
		return ta.alter.AddColumn(column, t, length, scale, sets...)
	}

	if ta.createMode {
		return ta.create.AddColumn(column, t, length, scale, sets...)
	}

	return nil
}

func (ta *TableBuilder) DropColumn(column string) ksql.TableInterface {
	if !ta.alterMode {
		return ta
	}

	ta.alter.DropColumn(column)
	return ta
}

func (ta *TableBuilder) DropColumnIfExists(column string) ksql.TableInterface {
	if !ta.alterMode {
		return ta
	}

	ta.alter.DropColumnIfExists(column)
	return ta
}

func (ta *TableBuilder) AddIndex(name string) ksql.TableIndexInterface {
	if ta.alterMode {
		return ta.alter.AddIndex(name)
	}

	return ta.create.AddIndex(name)
}

func (ta *TableBuilder) DropIndex(name string) ksql.TableInterface {
	if !ta.alterMode {
		return ta
	}

	ta.alter.DropIndex(name)
	return ta
}

func (ta *TableBuilder) Table(table string) ksql.TableInterface {
	ta.table = table
	if ta.alterMode {
		ta.alter.Table(table)
		return ta
	}

	if ta.createMode {
		ta.create.Table(table)
	}

	return ta
}

func (ta *TableBuilder) ChangeColumn(oldColumn, newColumn, t string, length, scale int, sets ...string) ksql.ColumnInterface {
	if !ta.alterMode {
		return nil
	}

	return ta.ChangeColumn(oldColumn, newColumn, t, length, scale, sets...)
}

func (t *TableBuilder) Exec(ctx context.Context) error {
	var op ksql.SqlInterface
	if t.alterMode {
		op = t.alter
	}
	if t.createMode {
		op = t.create
	}
	if op == nil {
		return nil
	}

	if t.conn != nil {
		_, err := t.conn.Exec(ctx, op)
		return err
	}

	_, err := Exec(ctx, op)
	return err
}

func (t *TableBuilder) Engine(engine string) ksql.TableInterface {
	if !t.createMode {
		return t
	}

	t.create.Engine(engine)
	return t
}

func (t *TableBuilder) Charset(charset string) ksql.TableInterface {
	if !t.createMode {
		return t
	}

	t.create.Charset(charset)
	return t
}

func (t *TableBuilder) Collate(collate string) ksql.TableInterface {
	if !t.createMode {
		return t
	}

	t.create.Collate(collate)
	return t
}

func (t *TableBuilder) Comment(comment string) ksql.TableInterface {
	if t.createMode {
		t.create.Comment(comment)
		return t
	}

	if t.alterMode {
		t.alter.Comment(comment)
		return t
	}

	return t
}

func (t *TableBuilder) AddUnique(name string, columns ...string) ksql.TableInterface {
	if t.createMode {
		t.create.AddUnique(name, columns...)
		return t
	}

	if t.alterMode {
		t.alter.AddUnique(name, columns...)
		return t
	}

	return t
}

func (t *TableBuilder) AddPrimary(column string) ksql.TableInterface {
	if t.createMode {
		t.create.AddPrimary(column)
		return t
	}

	if t.alterMode {
		t.alter.AddPrimary(column)
		return t
	}

	return t
}

func (t *TableBuilder) HasColumn(ctx context.Context, column string) (bool, error) {
	if t.conn != nil {
		return HasColumnBy(ctx, t.conn, t.table, column)
	}

	return HasColumn(ctx, t.table, column)
}

func (t *TableBuilder) HasIndex(ctx context.Context, index string) (bool, error) {
	if t.conn == nil {
		return HasIndex(ctx, t.table, index)
	}

	return HasIndexBy(ctx, t.conn, t.table, index)
}

func (t *TableBuilder) From(query ksql.QueryInterface) ksql.TableInterface {
	if t.createMode {
		t.create.As(query)
	}

	return t
}

func (t *TableBuilder) Like(table string) ksql.TableInterface {
	if t.createMode {
		t.create.Like(table)
	}

	return t
}
