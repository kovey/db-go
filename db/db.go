package db

import (
	ds "database/sql"
	"fmt"
	"reflect"

	"github.com/kovey/db-go/row"
	"github.com/kovey/db-go/sql"
	"github.com/kovey/debug-go/debug"
)

type DbInterface interface {
	Begin() error
	Commit() error
	RollBack() error
	InTransaction() bool
	Query(string, interface{}, ...interface{}) ([]interface{}, error)
	Exec(string) error
	Insert(*sql.Insert) (int64, error)
	Update(*sql.Update) (int64, error)
	Delete(*sql.Delete) (int64, error)
	BatchInsert(*sql.Batch) (int64, error)
	Select(*sql.Select, interface{}) ([]interface{}, error)
	FetchRow(string, map[string]interface{}, interface{}) (interface{}, error)
	FetchAll(string, map[string]interface{}, interface{}) ([]interface{}, error)
	FetchAllByWhere(string, *sql.Where, interface{}) ([]interface{}, error)
	FetchPage(string, map[string]interface{}, interface{}, int, int) ([]interface{}, error)
	FetchPageByWhere(string, *sql.Where, interface{}, int, int) ([]interface{}, error)
}

type ConnInterface interface {
	Query(string, ...interface{}) (*ds.Rows, error)
	Exec(string, ...interface{}) (ds.Result, error)
	Prepare(string) (*ds.Stmt, error)
	QueryRow(string, ...interface{}) *ds.Row
}

func Query(m ConnInterface, query string, t interface{}, args ...interface{}) ([]interface{}, error) {
	rows, err := m.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	res := make([]interface{}, 0)
	vType := reflect.TypeOf(t)
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}

	for rows.Next() {
		row := row.New(vType)
		row.Convert(rows)
		res = append(res, row.Value())
	}

	return res, nil
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

func Select(m ConnInterface, sel *sql.Select, t interface{}) ([]interface{}, error) {
	debug.Info("sql: %s", sel)
	return Query(m, sel.Prepare(), t, sel.Args()...)
}

func FetchRow(m ConnInterface, table string, where map[string]interface{}, t interface{}) (interface{}, error) {
	vType := reflect.TypeOf(t)
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}

	row := row.New(vType)

	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(row.Fields()...).Limit(1)

	debug.Info("sql: %s", sel)

	result := m.QueryRow(sel.Prepare(), sel.Args()...)

	if result.Err() != nil {
		return nil, result.Err()
	}

	if err := row.ConvertByRow(result); err != nil {
		return nil, err
	}

	return row.Value(), nil
}

func FetchAll(m ConnInterface, table string, where map[string]interface{}, t interface{}) ([]interface{}, error) {
	vType := reflect.TypeOf(t)
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}

	row := row.New(vType)
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(row.Fields()...)

	debug.Info("sql: %s", sel)

	return Query(m, sel.Prepare(), t, sel.Args()...)
}

func FetchAllByWhere(m ConnInterface, table string, where *sql.Where, t interface{}) ([]interface{}, error) {
	vType := reflect.TypeOf(t)
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}

	row := row.New(vType)
	sel := sql.NewSelect(table, "")
	sel.Where(where).Columns(row.Fields()...)

	debug.Info("sql: %s", sel)

	return Query(m, sel.Prepare(), t, sel.Args()...)
}

func FetchPage(m ConnInterface, table string, where map[string]interface{}, t interface{}, page int, pageSize int) ([]interface{}, error) {
	vType := reflect.TypeOf(t)
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}

	row := row.New(vType)
	sel := sql.NewSelect(table, "")
	sel.WhereByMap(where).Columns(row.Fields()...).Limit(pageSize).Offset((page - 1) * pageSize)

	debug.Info("sql: %s", sel)

	return Query(m, sel.Prepare(), t, sel.Args()...)
}

func FetchPageByWhere(m ConnInterface, table string, where *sql.Where, t interface{}, page int, pageSize int) ([]interface{}, error) {
	vType := reflect.TypeOf(t)
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}

	row := row.New(vType)
	sel := sql.NewSelect(table, "")
	sel.Where(where).Columns(row.Fields()...).Limit(pageSize).Offset((page - 1) * pageSize)

	debug.Info("sql: %s", sel)

	return Query(m, sel.Prepare(), t, sel.Args()...)
}
