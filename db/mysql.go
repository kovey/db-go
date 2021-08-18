package db

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kovey/db-go/row"
	ds "github.com/kovey/db-go/sql"
	"github.com/kovey/logger-go/logger"
)

var (
	database *sql.DB
)

type Mysql struct {
	database        *sql.DB
	tx              *sql.Tx
	isInTransaction bool
}

func NewMysql() *Mysql {
	return &Mysql{database: database, tx: nil, isInTransaction: false}
}

func Init(host string, port int, username string, password string, dbname string, charset string, maxActive int, maxConn int) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", username, password, host, port, dbname, charset))
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(maxActive)
	db.SetMaxOpenConns(maxConn)

	err = db.Ping()
	if err != nil {
		return err
	}

	database = db
	return nil
}

func (m *Mysql) Begin() error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}

	m.tx = tx
	m.isInTransaction = true
	return nil
}

func (m *Mysql) Commit() error {
	if m.tx == nil {
		return fmt.Errorf("transaction is not open or close")
	}

	m.isInTransaction = false
	return m.tx.Commit()
}

func (m *Mysql) RollBack() error {
	if m.tx == nil {
		return fmt.Errorf("transaction is not open or close")
	}

	m.isInTransaction = false
	return m.tx.Rollback()
}

func (m *Mysql) InTransaction() bool {
	return m.isInTransaction
}

func (m *Mysql) Query(query string, t interface{}, args ...interface{}) ([]interface{}, error) {
	var rows *sql.Rows
	var err error

	if m.isInTransaction {
		rows, err = m.tx.Query(query, args...)
	} else {
		rows, err = m.database.Query(query, args...)
	}

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

func (m *Mysql) Exec(stament string) error {
	var result sql.Result
	var err error

	if m.isInTransaction {
		result, err = m.tx.Exec(stament)
	} else {
		result, err = m.database.Exec(stament)
	}

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

func (m *Mysql) prepare(pre ds.SqlInterface) (sql.Result, error) {
	var smt *sql.Stmt
	var err error
	if m.isInTransaction {
		smt, err = m.tx.Prepare(pre.Prepare())
	} else {
		smt, err = m.database.Prepare(pre.Prepare())
	}

	if err != nil {
		return nil, err
	}

	defer smt.Close()

	return smt.Exec(pre.Args()...)
}

func (m *Mysql) Insert(insert *ds.Insert) (int64, error) {
	result, err := m.prepare(insert)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (m *Mysql) Update(update *ds.Update) (int64, error) {
	result, err := m.prepare(update)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (m *Mysql) Delete(del *ds.Delete) (int64, error) {
	result, err := m.prepare(del)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (m *Mysql) BatchInsert(batch *ds.Batch) (int64, error) {
	result, err := m.prepare(batch)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (m *Mysql) Select(sel *ds.Select, t interface{}) ([]interface{}, error) {
	return m.Query(sel.Prepare(), t, sel.Args()...)
}

func (m *Mysql) FetchRow(table string, where map[string]interface{}, t interface{}) (interface{}, error) {
	vType := reflect.TypeOf(t)
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}

	row := row.New(vType)

	sel := ds.NewSelect(table, "")
	sel.WhereByMap(where).Columns(row.Fields()...).Limit(1)

	var result *sql.Row

	if m.isInTransaction {
		result = m.tx.QueryRow(sel.Prepare(), sel.Args()...)
	} else {
		result = m.database.QueryRow(sel.Prepare(), sel.Args()...)
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	row.ConvertByRow(result)

	return row.Value(), nil
}

func (m *Mysql) FetchAll(table string, where map[string]interface{}, t interface{}) ([]interface{}, error) {
	vType := reflect.TypeOf(t)
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}

	row := row.New(vType)
	sel := ds.NewSelect(table, "")
	sel.WhereByMap(where).Columns(row.Fields()...)

	logger.Debug("fetch all select: %s", sel)

	return m.Query(sel.Prepare(), t, sel.Args()...)
}

func init() {
	logger.SetLevel(logger.LOGGER_INFO)
}
