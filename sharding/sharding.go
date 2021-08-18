package sharding

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kovey/config-go/config"
	"github.com/kovey/db-go/db"
	ds "github.com/kovey/db-go/sql"
	"github.com/kovey/logger-go/logger"
)

var (
	masters    []*sql.DB
	slaves     []*sql.DB
	mNodeCount int
	sNodeCount int
)

type Mysql struct {
	connections map[int]*db.Mysql
	isMaster    bool
}

func NewMysql(isMaster bool) *Mysql {
	return &Mysql{connections: make(map[int]*db.Mysql, 0), isMaster: isMaster}
}

func (m *Mysql) GetShardingKey(key interface{}) int {
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
	masters = make([]*sql.DB, 0, mNodeCount)
	slaves = make([]*sql.DB, 0, sNodeCount)

	for key, value := range mas {
		value.Dbname = fmt.Sprintf("%s_%d", value.Dbname, key)
		database, err := db.OpenDB(value)
		if err != nil {
			continue
		}

		masters[key] = database
	}

	for key, value := range sls {
		value.Dbname = fmt.Sprintf("%s_%d", value.Dbname, key)
		database, err := db.OpenDB(value)
		if err != nil {
			continue
		}

		slaves[key] = database
	}
}

func (m *Mysql) GetConnection(key int) *db.Mysql {
	logger.Debug("sharding key: %d", key)
	database, ok := m.connections[key]
	if ok {
		return database
	}

	if m.isMaster {
		database = db.NewSharding(masters[key])
	} else {
		database = db.NewSharding(slaves[key])
	}

	m.connections[key] = database

	return database
}

func (m *Mysql) rollBack(begins []int) {
	for _, id := range begins {
		m.connections[id].RollBack()
	}
}

func (m *Mysql) Begin() error {
	if len(m.connections) == 0 {
		return errors.New("connections is empty")
	}

	begins := make([]int, 0)
	for index, connection := range m.connections {
		err := connection.Begin()
		if err != nil {
			m.rollBack(begins)
			return err
		}

		begins = append(begins, index)
	}

	return nil
}

func (m *Mysql) retry(fails []int) {
	for _, id := range fails {
		m.connections[id].Commit()
	}
}

func (m *Mysql) Commit() error {
	if len(m.connections) == 0 {
		return errors.New("connections is empty")
	}

	i := 0
	fails := make([]int, 0)
	for index, connection := range m.connections {
		err := connection.Commit()
		if err != nil {
			if i == 0 {
				return err
			}

			fails = append(fails, index)
		}

		i++
	}

	if len(fails) > 0 {
		m.retry(fails)
	}

	return nil
}

func (m *Mysql) RollBack() error {
	if len(m.connections) == 0 {
		return errors.New("connections is empty")
	}

	for _, connection := range m.connections {
		connection.RollBack()
	}

	return nil
}

func (m *Mysql) Query(key interface{}, query string, t interface{}, args ...interface{}) ([]interface{}, error) {
	return m.GetConnection(m.GetShardingKey(key)).Query(query, t, args...)
}

func (m *Mysql) Exec(key interface{}, statement string) error {
	return m.GetConnection(m.GetShardingKey(key)).Exec(statement)
}

func (m *Mysql) Insert(key interface{}, insert *ds.Insert) (int64, error) {
	return m.GetConnection(m.GetShardingKey(key)).Insert(insert)
}

func (m *Mysql) Update(key interface{}, update *ds.Update) (int64, error) {
	return m.GetConnection(m.GetShardingKey(key)).Update(update)
}

func (m *Mysql) Delete(key interface{}, del *ds.Delete) (int64, error) {
	return m.GetConnection(m.GetShardingKey(key)).Delete(del)
}

func (m *Mysql) BatchInsert(key interface{}, batch *ds.Batch) (int64, error) {
	return m.GetConnection(m.GetShardingKey(key)).BatchInsert(batch)
}

func (m *Mysql) FetchRow(key interface{}, table string, where map[string]interface{}, t interface{}) (interface{}, error) {
	return m.GetConnection(m.GetShardingKey(key)).FetchRow(table, where, t)
}

func (m *Mysql) FetchAll(key interface{}, table string, where map[string]interface{}, t interface{}) ([]interface{}, error) {
	return m.GetConnection(m.GetShardingKey(key)).FetchAll(table, where, t)
}
