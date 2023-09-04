package db

import (
	ds "database/sql"
	"fmt"

	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/debug-go/debug"
)

type DbInterface[T itf.RowInterface] interface {
	SetTx(*Tx)
	Transaction(func(*Tx) error) error
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
	FetchRow(string, meta.Where, T) error
	LockRow(string, meta.Where, T) error
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

func Query[T itf.RowInterface](m ConnInterface, query string, model T, args ...any) ([]T, error) {
	data, err := m.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer data.Close()
	rows := make([]T, 0)
	for data.Next() {
		tmp := model.Clone()
		if err := data.Scan(tmp.Fields()...); err != nil {
			return nil, err
		}

		rows = append(rows, tmp.(T))
	}

	return rows, nil
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

func ShowTables[T itf.RowInterface](m ConnInterface, show *sql.ShowTables, model T) ([]T, error) {
	data, err := m.Query(show.Prepare(), show.Args()...)
	if err != nil {
		return nil, err
	}

	defer data.Close()
	rows := make([]T, 0)
	for data.Next() {
		tmp := model.Clone()
		if err := data.Scan(tmp.Fields()...); err != nil {
			return nil, err
		}

		rows = append(rows, tmp.(T))
	}

	return rows, nil
}

func Desc[T itf.RowInterface](m ConnInterface, desc *sql.Desc, model T) ([]T, error) {
	return Query(m, desc.Prepare(), model, desc.Args()...)
}

func Select[T itf.RowInterface](m ConnInterface, sel *sql.Select, model T) ([]T, error) {
	return Query(m, sel.Prepare(), model, sel.Args()...)
}

func FetchRow[T itf.ModelInterface](m ConnInterface, table string, where meta.Where, model T) error {
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(model.Columns()...).Limit(1)
	result := m.QueryRow(sel.Prepare(), sel.Args()...)
	return parseError(result.Scan(model.Fields()...), model)
}

func parseError[T itf.ModelInterface](err error, model T) error {
	if err == nil {
		return nil
	}

	if err == ds.ErrNoRows {
		model.SetEmpty()
		return nil
	}

	return err
}

func LockRow[T itf.ModelInterface](m ConnInterface, table string, where meta.Where, model T) error {
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(model.Columns()...).Limit(1).ForUpdate()
	result := m.QueryRow(sel.Prepare(), sel.Args()...)
	return parseError(result.Scan(model.Fields()...), model)
}

func FetchAll[T itf.RowInterface](m ConnInterface, table string, where meta.Where, model T) ([]T, error) {
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(model.Columns()...)

	return Select(m, sel, model)
}

func FetchAllByWhere[T itf.RowInterface](m ConnInterface, table string, where sql.WhereInterface, model T) ([]T, error) {
	sel := sql.NewSelect(table, "")
	sel.Where(where).Columns(model.Columns()...)

	return Select(m, sel, model)
}

func FetchPage[T itf.RowInterface](m ConnInterface, table string, where meta.Where, model T, page int, pageSize int) ([]T, error) {
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(model.Columns()...).Limit(pageSize).Offset((page - 1) * pageSize)

	return Select(m, sel, model)
}

func FetchPageByWhere[T itf.RowInterface](m ConnInterface, table string, where sql.WhereInterface, model T, page int, pageSize int) ([]T, error) {
	sel := sql.NewSelect(table, "")
	sel.Where(where).Columns(model.Columns()...).Limit(pageSize).Offset((page - 1) * pageSize)

	return Select(m, sel, model)
}
