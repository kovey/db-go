package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kovey/config-go/config"
	ds "github.com/kovey/db-go/sql"
	"github.com/kovey/logger-go/logger"
)

var (
	database *sql.DB
	dev      string
)

type Mysql struct {
	database        *sql.DB
	tx              *sql.Tx
	isInTransaction bool
}

func NewMysql() *Mysql {
	return &Mysql{database: database, tx: nil, isInTransaction: false}
}

func NewSharding(database *sql.DB) *Mysql {
	return &Mysql{database: database, tx: nil, isInTransaction: false}
}

func Init(conf config.Mysql) error {
	db, err := OpenDB(conf)
	if err != nil {
		return err
	}

	dev = conf.Dev
	database = db
	return nil
}

func OpenDB(conf config.Mysql) (*sql.DB, error) {
	logger.Debug("connection to %s:%d, dbname: %s", conf.Host, conf.Port, conf.Dbname)
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Dbname, conf.Charset,
	))
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(conf.ActiveMax)
	db.SetMaxOpenConns(conf.ConnectionMax)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m *Mysql) Begin() error {
	tx, err := m.database.Begin()
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
	err := m.tx.Commit()
	m.tx = nil
	return err
}

func (m *Mysql) RollBack() error {
	if m.tx == nil {
		return fmt.Errorf("transaction is not open or close")
	}

	m.isInTransaction = false
	err := m.tx.Rollback()
	m.tx = nil
	return err
}

func (m *Mysql) InTransaction() bool {
	return m.isInTransaction
}

func (m *Mysql) Query(query string, t interface{}, args ...interface{}) ([]interface{}, error) {
	return Query(m.getDb(), query, t, args...)
}

func (m *Mysql) Exec(statement string) error {
	return Exec(m.getDb(), statement)
}

func (m *Mysql) Insert(insert *ds.Insert) (int64, error) {
	return Insert(m.getDb(), insert)
}

func (m *Mysql) getDb() ConnInterface {
	if m.isInTransaction {
		return m.tx
	}

	return m.database
}

func (m *Mysql) Update(update *ds.Update) (int64, error) {
	return Update(m.getDb(), update)
}

func (m *Mysql) Delete(del *ds.Delete) (int64, error) {
	return Delete(m.getDb(), del)
}

func (m *Mysql) BatchInsert(batch *ds.Batch) (int64, error) {
	return BatchInsert(m.getDb(), batch)
}

func (m *Mysql) Select(sel *ds.Select, t interface{}) ([]interface{}, error) {
	return Select(m.getDb(), sel, t)
}

func (m *Mysql) FetchRow(table string, where map[string]interface{}, t interface{}) (interface{}, error) {
	return FetchRow(m.getDb(), table, where, t)
}

func (m *Mysql) FetchAll(table string, where map[string]interface{}, t interface{}) ([]interface{}, error) {
	return FetchAll(m.getDb(), table, where, t)
}

func (m *Mysql) FetchByPage(table string, where map[string]interface{}, t interface{}, page int, pageSize int) ([]interface{}, error) {
	return FetchByPage(m.getDb(), table, where, t, page, pageSize)
}
