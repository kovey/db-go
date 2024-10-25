package mysql

import (
	"context"
	"fmt"

	"github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/db-go/v3/migrate/schema"
	"github.com/kovey/db-go/v3/sql"
)

func DiffSchema(ctx context.Context, from, to ksql.ConnectionInterface, fromDbname, toDbname string) ksql.SqlInterface {
	var schema = &Schema{base: &base{conn: from, empty: true}}
	if err := db.Build(schema).Table("information_schema.SCHEMATA").Columns(schema.Columns()...).Where("SCHEMA_NAME", "=", fromDbname).First(ctx, schema); err != nil {
		return nil
	}

	if schema.empty {
		return nil
	}

	var toSchema = &Schema{base: &base{conn: to, empty: true}}
	if err := db.Build(toSchema).Table("information_schema.SCHEMATA").Columns(toSchema.Columns()...).Where("SCHEMA_NAME", "=", toDbname).First(ctx, toSchema); err != nil {
		return nil
	}

	var op = sql.NewSchema()
	if toSchema.empty {
		op.Create().Charset(schema.Charset()).Collate(schema.Collation()).Schema(schema.Name())
		return op
	}

	if !schema.HasChanged(toSchema) {
		return nil
	}

	op.Alter().Schema(schema.Name())
	if schema.Charset() != toSchema.Charset() {
		op.Charset(schema.Charset())
	}
	if schema.Collation() != toSchema.Collation() {
		op.Collate(schema.Collation())
	}

	return op
}

func DiffTable(ctx context.Context, fromTable, toTable schema.TableInfoInterface) ksql.SqlInterface {
	if fromTable == nil {
		if toTable == nil {
			return nil
		}

		return sql.NewDropTable().Table(toTable.Name())
	}

	if toTable == nil {
		op := sql.NewTable().Table(fromTable.Name()).Engine(fromTable.Engine()).Charset(fromTable.Charset()).Collate(fromTable.Collation()).Comment(fromTable.Comment())
		for _, column := range fromTable.Fields() {
			c := op.AddColumn(column.Name(), column.Type(), column.Length(), column.Scale()).Comment(column.Comment())
			if column.AutoIncrement() {
				c.AutoIncrement()
			}
			if column.HasDefault() {
				c.Default(column.Default(), ksql.IsDefaultKeyword(column.Default()))
			}
			if column.Nullable() {
				c.Nullable()
			}
		}

		for _, index := range fromTable.Indexes() {
			if index.Type() == ksql.Index_Type_Primary {
				op.AddPrimary(index.Columns()[0])
				continue
			}

			op.AddIndex(index.Name(), index.Type(), index.Columns()...)
		}

		return op
	}

	if !fromTable.HasChanged(toTable) {
		return nil
	}

	op := sql.NewAlter().Table(toTable.Name())
	if fromTable.Engine() != toTable.Engine() {
		op.Engine(fromTable.Engine())
	}
	if fromTable.Charset() != toTable.Charset() {
		op.Charset(fromTable.Charset())
	}
	if fromTable.Collation() != toTable.Collation() {
		op.Collate(fromTable.Collation())
	}
	if fromTable.Comment() != toTable.Comment() {
		op.Comment(fromTable.Comment())
	}

	changes := fromTable.CheckChanges(toTable)
	if changes == nil {
		return op
	}

	if changes.Index() != nil {
		for _, index := range changes.Index().Adds() {
			if index.Type() == ksql.Index_Type_Primary {
				op.AddPrimary(index.Columns()[0])
				continue
			}

			op.AddIndex(index.Name(), index.Type(), index.Columns()...)
		}

		for _, index := range changes.Index().Deletes() {
			op.DropIndex(index.Name())
		}
	}

	if changes.Column() != nil {
		for _, column := range changes.Column().Adds() {
			c := op.AddColumn(column.Name(), column.Type(), column.Length(), column.Scale()).Comment(column.Comment())
			if column.AutoIncrement() {
				c.AutoIncrement()
			}
			if column.HasDefault() {
				c.Default(column.Default(), ksql.IsDefaultKeyword(column.Default()))
			}
			if column.Nullable() {
				c.Nullable()
			}
		}
		for _, change := range changes.Column().Changes() {
			column := change.New()
			c := op.ChangeColumn(change.Old().Name(), column.Name(), column.Type(), column.Length(), column.Scale()).Comment(column.Comment())
			if column.AutoIncrement() {
				c.AutoIncrement()
			}
			if column.HasDefault() {
				c.Default(column.Default(), ksql.IsDefaultKeyword(column.Default()))
			}
			if column.Nullable() {
				c.Nullable()
			}
		}
		for _, del := range changes.Column().Deletes() {
			op.DropColumn(del.Name())
		}
	}

	return op
}

func Tables(ctx context.Context, conn ksql.ConnectionInterface, dbname string) ([]schema.TableInfoInterface, error) {
	var fromTable = NewTable(conn)
	var tables []*Table
	if err := db.Build(fromTable).Table("information_schema.TABLES").Columns(fromTable.Columns()...).Where("TABLE_SCHEMA", "=", dbname).All(ctx, &tables); err != nil {
		return nil, err
	}

	tmp := make([]schema.TableInfoInterface, len(tables))
	for i, t := range tables {
		fmt.Print(".")
		if err := _tableInfo(ctx, conn, dbname, t.Name(), t); err != nil {
			return nil, err
		}

		tmp[i] = t
	}
	fmt.Println("")
	return tmp, nil
}

func _tableInfo(ctx context.Context, conn ksql.ConnectionInterface, dbname, table string, tableModel *Table) error {
	var columns []*Column
	column := &Column{base: &base{conn: conn, empty: true}}
	if err := db.Build(column).Table("information_schema.COLUMNS").Columns(column.Columns()...).Where("TABLE_SCHEMA", "=", dbname).Where("TABLE_NAME", "=", table).All(ctx, &columns); err != nil {
		return err
	}

	tmp := make([]schema.ColumnInfoInterface, len(columns))
	for i, column := range columns {
		tmp[i] = column
	}

	tableModel.SetColumns(tmp)

	raw := sql.Raw("SHOW INDEX FROM `" + table + "`")
	var indexs []*Index
	if err := db.QueryRawBy(ctx, conn, raw, &indexs); err != nil {
		return err
	}

	tmpIndexes := make([]schema.IndexMetaInterface, len(indexs))
	for i, index := range indexs {
		tmpIndexes[i] = index
	}

	tableModel.SetIndexes(tmpIndexes)
	return nil
}
