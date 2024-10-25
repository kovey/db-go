package db

import (
	"context"

	"github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql"
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

func (ta *TableBuilder) SetConn(conn ksql.ConnectionInterface) ksql.TableInterface {
	ta.conn = conn
	return ta
}

func (ta *TableBuilder) Create() ksql.TableInterface {
	if ta.alterMode {
		return ta
	}

	ta.createMode = true
	ta.create = sql.NewTable()
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
	ta.alter = sql.NewAlter()
	if ta.table != "" {
		ta.alter.Table(ta.table)
	}
	return ta
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

func (ta *TableBuilder) AddIndex(name string, t ksql.IndexType, column ...string) ksql.TableInterface {
	if ta.alterMode {
		ta.alter.AddIndex(name, t, column...)
		return ta
	}

	ta.create.AddIndex(name, t, column...)
	return ta
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
