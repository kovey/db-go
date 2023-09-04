package sharding

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/db"
	"github.com/kovey/db-go/v2/itf"
	ds "github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/debug-go/debug"
)

var (
	masters    []*sql.DB
	slaves     []*sql.DB
	mNodeCount int
	sNodeCount int
)

type Mysql[T itf.ModelInterface] struct {
	connections map[int]*db.Mysql[T]
	isMaster    bool
	tx          *Tx
}

func NewMysql[T itf.ModelInterface](isMaster bool) *Mysql[T] {
	return &Mysql[T]{connections: make(map[int]*db.Mysql[T], 0), isMaster: isMaster, tx: NewTx()}
}

func (m *Mysql[T]) SetTx(tx *Tx) {
	for key, tx := range tx.txs {
		m.tx.Add(key, tx)
	}
}

func (m *Mysql[T]) Transaction(f func(tx *Tx) error) error {
	if err := m.begin(); err != nil {
		return err
	}

	if err := f(m.tx); err != nil {
		if err := m.rollBack(); err != nil {
			debug.Erro("rollBack failure, error: %s", err)
		}
		return err
	}

	return m.commit()
}

func (m *Mysql[T]) AddSharding(key any) *Mysql[T] {
	k := m.GetShardingKey(key)
	if _, ok := m.connections[k]; ok {
		return m
	}

	var database *db.Mysql[T]
	if m.isMaster {
		database = db.NewSharding[T](masters[k])
	} else {
		database = db.NewSharding[T](slaves[k])
	}

	m.connections[k] = database

	return m
}

func (m *Mysql[T]) GetShardingKey(key any) int {
	k, ok := key.(string)
	if ok {
		if m.isMaster {
			return int(getHashKey(k) % uint32(mNodeCount))
		}

		return int(getHashKey(k) % uint32(sNodeCount))
	}

	v, okk := key.(int)
	if !okk {
		return 0
	}

	if m.isMaster {
		return v % mNodeCount
	}

	return v % sNodeCount
}

func Init(mas []config.Mysql, sls []config.Mysql) {
	mNodeCount = len(mas)
	sNodeCount = len(sls)
	masters = make([]*sql.DB, mNodeCount)
	slaves = make([]*sql.DB, sNodeCount)

	for key, value := range mas {
		value.Dbname = fmt.Sprintf("%s_%d", value.Dbname, key)
		database, err := db.OpenDB(value)
		if err != nil {
			debug.Erro("open master database failure, error: %s", err)
			continue
		}

		masters[key] = database
	}

	for key, value := range sls {
		value.Dbname = fmt.Sprintf("%s_%d", value.Dbname, key)
		database, err := db.OpenDB(value)
		if err != nil {
			debug.Erro("open slave database failure, error: %s", err)
			continue
		}

		slaves[key] = database
	}
}

func (m *Mysql[T]) GetConnection(key int) *db.Mysql[T] {
	database, ok := m.connections[key]
	if ok {
		return database
	}

	if m.isMaster {
		database = db.NewSharding[T](masters[key])
	} else {
		database = db.NewSharding[T](slaves[key])
	}

	m.connections[key] = database

	return database
}

func (m *Mysql[T]) begin() error {
	if len(m.connections) == 0 {
		return errors.New("connections is empty")
	}

	for index, connection := range m.connections {
		err := connection.Begin()
		if err != nil {
			m.tx.Rollback()
			return err
		}

		m.tx.Add(index, connection.Tx())
	}

	return nil
}

func (m *Mysql[T]) commit() error {
	if len(m.connections) == 0 {
		return errors.New("connections is empty")
	}

	m.tx.Commit()
	return nil
}

func (m *Mysql[T]) rollBack() error {
	if len(m.connections) == 0 {
		return errors.New("connections is empty")
	}

	m.tx.Rollback()
	return nil
}

func (m *Mysql[T]) Query(key any, query string, model T, args ...any) ([]T, error) {
	return m.GetConnection(m.GetShardingKey(key)).Query(query, model, args...)
}

func (m *Mysql[T]) Exec(key any, statement string) error {
	return m.GetConnection(m.GetShardingKey(key)).Exec(statement)
}

func (m *Mysql[T]) Insert(key any, insert *ds.Insert) (int64, error) {
	return m.GetConnection(m.GetShardingKey(key)).Insert(insert)
}

func (m *Mysql[T]) Update(key any, update *ds.Update) (int64, error) {
	return m.GetConnection(m.GetShardingKey(key)).Update(update)
}

func (m *Mysql[T]) Delete(key any, del *ds.Delete) (int64, error) {
	return m.GetConnection(m.GetShardingKey(key)).Delete(del)
}

func (m *Mysql[T]) BatchInsert(key any, batch *ds.Batch) (int64, error) {
	return m.GetConnection(m.GetShardingKey(key)).BatchInsert(batch)
}

func (m *Mysql[T]) FetchRow(key any, table string, where meta.Where, model T) error {
	return m.GetConnection(m.GetShardingKey(key)).FetchRow(table, where, model)
}

func (m *Mysql[T]) LockRow(key any, table string, where meta.Where, model T) error {
	return m.GetConnection(m.GetShardingKey(key)).LockRow(table, where, model)
}

func (m *Mysql[T]) FetchAll(key any, table string, where meta.Where, model T) ([]T, error) {
	return m.GetConnection(m.GetShardingKey(key)).FetchAll(table, where, model)
}

func (m *Mysql[T]) FetchAllByWhere(key any, table string, where ds.WhereInterface, model T) ([]T, error) {
	return m.GetConnection(m.GetShardingKey(key)).FetchAllByWhere(table, where, model)
}

func (m *Mysql[T]) FetchPage(key any, table string, where meta.Where, model T, page, pageSize int) ([]T, error) {
	return m.GetConnection(m.GetShardingKey(key)).FetchPage(table, where, model, page, pageSize)
}

func (m *Mysql[T]) FetchPageByWhere(key any, table string, where ds.WhereInterface, model T, page, pageSize int) ([]T, error) {
	return m.GetConnection(m.GetShardingKey(key)).FetchPageByWhere(table, where, model, page, pageSize)
}
