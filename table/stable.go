package table

import (
	"fmt"

	"github.com/kovey/db-go/sharding"
	"github.com/kovey/db-go/sql"
)

type TableShardingInterface interface {
	Database() *sharding.Mysql
	Insert(interface{}, map[string]interface{}) (int64, error)
	Update(interface{}, map[string]interface{}, map[string]interface{}) (int64, error)
	Delete(interface{}, map[string]interface{}) (int64, error)
	BatchInsert(interface{}, []map[string]interface{}) (int64, error)
	FetchRow(interface{}, map[string]interface{}, interface{}) (interface{}, error)
	FetchAll(interface{}, map[string]interface{}, interface{}) ([]interface{}, error)
	FetchAllByWhere(interface{}, *sql.Where, interface{}) ([]interface{}, error)
	FetchPage(interface{}, map[string]interface{}, interface{}, int, int) ([]interface{}, error)
	FetchPageByWhere(interface{}, *sql.Where, interface{}, int, int) ([]interface{}, error)
}

type TableSharding struct {
	table string
	db    *sharding.Mysql
}

func NewTableSharding(table string, isMaster bool) *TableSharding {
	return &TableSharding{db: sharding.NewMysql(isMaster), table: table}
}

func (t *TableSharding) Database() *sharding.Mysql {
	return t.db
}

func (t *TableSharding) GetTableName(key interface{}) string {
	return fmt.Sprintf("%s_%d", t.table, t.db.GetShardingKey(key))
}

func (t *TableSharding) Insert(key interface{}, data map[string]interface{}) (int64, error) {
	in := sql.NewInsert(t.GetTableName(key))
	for field, value := range data {
		in.Set(field, value)
	}

	return t.db.Insert(key, in)
}

func (t *TableSharding) Update(key interface{}, data map[string]interface{}, where map[string]interface{}) (int64, error) {
	up := sql.NewUpdate(t.GetTableName(key))
	for field, value := range data {
		up.Set(field, value)
	}

	up.WhereByMap(where)

	return t.db.Update(key, up)
}

func (t *TableSharding) Delete(key interface{}, where map[string]interface{}) (int64, error) {
	del := sql.NewDelete(t.GetTableName(key))
	del.WhereByMap(where)

	return t.db.Delete(key, del)
}

func (t *TableSharding) BatchInsert(key interface{}, data []map[string]interface{}) (int64, error) {
	batch := sql.NewBatch(t.GetTableName(key))
	for _, val := range data {
		in := sql.NewInsert(t.GetTableName(key))
		for field, value := range val {
			in.Set(field, value)
		}
		batch.Add(in)
	}

	return t.db.BatchInsert(key, batch)
}

func (t *TableSharding) FetchRow(key interface{}, where map[string]interface{}, mt interface{}) (interface{}, error) {
	return t.db.FetchRow(key, t.GetTableName(key), where, mt)
}

func (t *TableSharding) FetchAll(key interface{}, where map[string]interface{}, mt interface{}) ([]interface{}, error) {
	return t.db.FetchAll(key, t.GetTableName(key), where, mt)
}

func (t *TableSharding) FetchAllByWhere(key interface{}, where *sql.Where, mt interface{}) ([]interface{}, error) {
	return t.db.FetchAllByWhere(key, t.GetTableName(key), where, mt)
}

func (t *TableSharding) FetchPage(key interface{}, where map[string]interface{}, mt interface{}, page, pageSize int) ([]interface{}, error) {
	return t.db.FetchPage(key, t.GetTableName(key), where, mt, page, pageSize)
}

func (t *TableSharding) FetchPageByWhere(key interface{}, where *sql.Where, mt interface{}, page, pageSize int) ([]interface{}, error) {
	return t.db.FetchPageByWhere(key, t.GetTableName(key), where, mt, page, pageSize)
}
