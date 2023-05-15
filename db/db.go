package db

import (
	ds "database/sql"
	"fmt"

	"github.com/kovey/db-go/v2/rows"
	"github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/debug-go/debug"
)

type DbInterface[T any] interface {
	Begin() error
	Commit() error
	RollBack() error
	InTransaction() bool
	Query(string, T, ...any) ([]T, error)
	Exec(string) error
	Desc(*sql.Desc, T) ([]T, error)
	ShowTables(*sql.ShowTables, T) ([]T, error)
	Insert(*sql.Insert) (int64, error)
	Update(*sql.Update) (int64, error)
	Delete(*sql.Delete) (int64, error)
	BatchInsert(*sql.Batch) (int64, error)
	Select(*sql.Select, T) ([]T, error)
	FetchRow(string, meta.Where, T) (T, error)
	FetchAll(string, meta.Where, T) ([]T, error)
	FetchAllByWhere(string, sql.WhereInterface, T) ([]T, error)
	FetchPage(string, meta.Where, T, int, int) ([]T, error)
	FetchPageByWhere(string, sql.WhereInterface, T, int, int) ([]T, error)
}

type ConnInterface interface {
	Query(string, ...any) (*ds.Rows, error)
	Exec(string, ...any) (ds.Result, error)
	Prepare(string) (*ds.Stmt, error)
	QueryRow(string, ...any) *ds.Row
}

func Query[T any](m ConnInterface, query string, model T, args ...any) ([]T, error) {
	data, err := m.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer data.Close()
	res := rows.NewRows[T]()
	if err := res.Scan(data, model); err != nil {
		return nil, err
	}

	return res.All(), nil
}

func Exec(m ConnInterface, stament string) error {
	debug.Info("sql: %s", stament)
	result, err := m.Exec(stament)

	if err != nil {
		return err
	}

	lastId, _ := result.LastInsertId()
	affected, _ := result.RowsAffected()

	if lastId < 1 && affected < 1 {
		return fmt.Errorf("lastId or affectedId is zero")
	}

	return nil
}

func prepare(m ConnInterface, pre sql.SqlInterface) (ds.Result, error) {
	debug.Info("sql: %s", pre)
	smt, err := m.Prepare(pre.Prepare())
	if err != nil {
		return nil, err
	}

	defer smt.Close()

	return smt.Exec(pre.Args()...)
}

func Insert(m ConnInterface, insert *sql.Insert) (int64, error) {
	result, err := prepare(m, insert)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func Update(m ConnInterface, update *sql.Update) (int64, error) {
	result, err := prepare(m, update)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func Delete(m ConnInterface, del *sql.Delete) (int64, error) {
	result, err := prepare(m, del)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func BatchInsert(m ConnInterface, batch *sql.Batch) (int64, error) {
	result, err := prepare(m, batch)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func ShowTables[T any](m ConnInterface, show *sql.ShowTables, model T) ([]T, error) {
	return Query(m, show.Prepare(), model, show.Args()...)
}

func Desc[T any](m ConnInterface, desc *sql.Desc, model T) ([]T, error) {
	return Query(m, desc.Prepare(), model, desc.Args()...)
}

func Select[T any](m ConnInterface, sel *sql.Select, model T) ([]T, error) {
	return Query(m, sel.Prepare(), model, sel.Args()...)
}

func FetchRow[T any](m ConnInterface, table string, where meta.Where, model T) (T, error) {
	row := rows.NewRow(model)
	row.Table = table
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(row.Fields...).Limit(1)
	result := m.QueryRow(sel.Prepare(), sel.Args()...)
	if result.Err() != nil {
		return row.Model, result.Err()
	}

	if err := row.Scan(result); err != nil {
		return row.Model, err
	}

	return row.Model, nil
}

func FetchAll[T any](m ConnInterface, table string, where meta.Where, model T) ([]T, error) {
	row := rows.NewRow(model)
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(row.Fields...)

	return Select(m, sel, model)
}

func FetchAllByWhere[T any](m ConnInterface, table string, where sql.WhereInterface, model T) ([]T, error) {
	row := rows.NewRow(model)
	sel := sql.NewSelect(table, "")
	sel.Where(where).Columns(row.Fields...)

	return Select(m, sel, model)
}

func FetchPage[T any](m ConnInterface, table string, where meta.Where, model T, page int, pageSize int) ([]T, error) {
	row := rows.NewRow(model)
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(row.Fields...).Limit(pageSize).Offset((page - 1) * pageSize)

	return Select(m, sel, model)
}

func FetchPageByWhere[T any](m ConnInterface, table string, where sql.WhereInterface, model T, page int, pageSize int) ([]T, error) {
	row := rows.NewRow(model)
	sel := sql.NewSelect(table, "")
	sel.Where(where).Columns(row.Fields...).Limit(pageSize).Offset((page - 1) * pageSize)

	return Select(m, sel, model)
}
