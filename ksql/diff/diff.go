package diff

import (
	"context"
	"database/sql"
	"strings"

	"github.com/kovey/db-go/ksql/mysql"
	"github.com/kovey/db-go/ksql/schema"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/debug-go/debug"
)

func Diff(ctx context.Context, driverName, fromDsn, toDsn, fromDbname, toDbname string) ([]ksql.SqlInterface, error) {
	var ops []ksql.SqlInterface
	from, err := getConn(driverName, fromDsn)
	if err != nil {
		return ops, err
	}

	to, err := getConn(driverName, toDsn)
	if err != nil {
		return ops, err
	}

	fromTables, err := getTables(ctx, from, fromDbname)
	if err != nil {
		return ops, err
	}
	toTables, err := getTables(ctx, to, toDbname)
	if err != nil {
		return ops, err
	}

	op, err := diffSchema(ctx, driverName, from, to, fromDbname, toDbname)
	if err != nil {
		return ops, err
	}

	if op != nil {
		ops = append(ops, op)
	}

	debug.Info("diff table begin...")
	defer debug.Info("diff table end.")
	for _, fromTable := range fromTables {
		if fromTable.Name() == "ksql_migrate_info" {
			continue
		}
		op, err := diffTable(ctx, driverName, fromTable, getTable(fromTable.Name(), toTables))
		if err != nil {
			return nil, err
		}

		if op != nil {
			ops = append(ops, op)
		}
	}

	for _, toTable := range toTables {
		if toTable.Name() == "ksql_migrate_info" {
			continue
		}
		fromTable := getTable(toTable.Name(), fromTables)
		if fromTable != nil {
			continue
		}

		op, err := diffTable(ctx, driverName, fromTable, toTable)
		if err != nil {
			return nil, err
		}

		if op != nil {
			ops = append(ops, op)
		}
	}

	return ops, nil
}

func getTable(table string, tables []schema.TableInfoInterface) schema.TableInfoInterface {
	for _, t := range tables {
		if table == t.Name() {
			return t
		}
	}

	return nil
}

func diffSchema(ctx context.Context, driverName string, from, to ksql.ConnectionInterface, fromDbname, toDbname string) (ksql.SqlInterface, error) {
	debug.Info("diff schema begin, [%s] -> [%s]", fromDbname, toDbname)
	defer debug.Info("diff schema end, [%s] -> [%s]", fromDbname, toDbname)
	switch strings.ToLower(driverName) {
	case "mysql":
		return mysql.DiffSchema(ctx, from, to, fromDbname, toDbname), nil

	}
	return nil, nil
}

func diffTable(ctx context.Context, driverName string, from, to schema.TableInfoInterface) (ksql.SqlInterface, error) {
	if from != nil && to != nil {
		debug.Info("diff table begin, [%s] -> [%s]", from.Name(), to.Name())
		defer debug.Info("diff schema end, [%s] -> [%s]", from.Name(), to.Name())
	} else if from != nil {
		debug.Info("diff table begin, [%s] -> [null]", from.Name())
		defer debug.Info("diff schema end, [%s] -> [null]", from.Name())
	} else {
		debug.Info("diff table begin, [null] -> [%s]", to.Name())
		defer debug.Info("diff schema end, [null] -> [%s]", to.Name())
	}

	switch strings.ToLower(driverName) {
	case "mysql":
		return mysql.DiffTable(ctx, from, to), nil
	}

	return nil, nil
}

func getConn(driverName, dsn string) (ksql.ConnectionInterface, error) {
	conn, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	return db.Open(conn, driverName)
}

func getTables(ctx context.Context, conn ksql.ConnectionInterface, dbname string) ([]schema.TableInfoInterface, error) {
	switch strings.ToLower(conn.DriverName()) {
	case "mysql":
		return mysql.Tables(ctx, conn, dbname)
	}

	return nil, nil
}
