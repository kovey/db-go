package table

import (
	"fmt"

	"github.com/kovey/db-go/sharding"
	"github.com/kovey/db-go/sql"
)

type TableShardingInterface[T any] interface {
	Database() *sharding.Mysql[T]
	Insert(any, map[string]any) (int64, error)
	Update(any, map[string]any, map[string]any) (int64, error)
	Delete(any, map[string]any) (int64, error)
	DeleteWhere(any, sql.WhereInterface) (int64, error)
	BatchInsert(any, []map[string]any) (int64, error)
	FetchRow(any, map[string]any, T) (T, error)
	FetchAll(any, map[string]any, T) ([]T, error)
	FetchAllByWhere(any, sql.WhereInterface, T) ([]T, error)
	FetchPage(any, map[string]any, T, int, int) ([]T, error)
	FetchPageByWhere(any, sql.WhereInterface, T, int, int) ([]T, error)
}

type TableSharding[T any] struct {
	table string
	db    *sharding.Mysql[T]
}

func NewTableSharding[T any](table string, isMaster bool) *TableSharding[T] {
	return &TableSharding[T]{db: sharding.NewMysql[T](isMaster), table: table}
}

func (t *TableSharding[T]) Database() *sharding.Mysql[T] {
	return t.db
}

func (t *TableSharding[T]) GetTableName(key any) string {
	return fmt.Sprintf("%s_%d", t.table, t.db.GetShardingKey(key))
}

func (t *TableSharding[T]) Insert(key any, data map[string]any) (int64, error) {
	in := sql.NewInsert(t.GetTableName(key))
	for field, value := range data {
		in.Set(field, value)
	}

	return t.db.Insert(key, in)
}

func (t *TableSharding[T]) Update(key any, data map[string]any, where map[string]any) (int64, error) {
	up := sql.NewUpdate(t.GetTableName(key))
	for field, value := range data {
		up.Set(field, value)
	}

	up.WhereByMap(where)

	return t.db.Update(key, up)
}

func (t *TableSharding[T]) Delete(key any, where map[string]any) (int64, error) {
	del := sql.NewDelete(t.GetTableName(key))
	del.WhereByMap(where)

	return t.db.Delete(key, del)
}

func (t *TableSharding[T]) DeleteWhere(key any, where sql.WhereInterface) (int64, error) {
	del := sql.NewDelete(t.GetTableName(key))
	del.Where(where)

	return t.db.Delete(key, del)
}

func (t *TableSharding[T]) BatchInsert(key any, data []map[string]any) (int64, error) {
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

func (t *TableSharding[T]) FetchRow(key any, where map[string]any, model T) (T, error) {
	return t.db.FetchRow(key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) FetchAll(key any, where map[string]any, model T) ([]T, error) {
	return t.db.FetchAll(key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) FetchAllByWhere(key any, where sql.WhereInterface, model T) ([]T, error) {
	return t.db.FetchAllByWhere(key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) FetchPage(key any, where map[string]any, model T, page, pageSize int) ([]T, error) {
	return t.db.FetchPage(key, t.GetTableName(key), where, model, page, pageSize)
}

func (t *TableSharding[T]) FetchPageByWhere(key any, where sql.WhereInterface, model T, page, pageSize int) ([]T, error) {
	return t.db.FetchPageByWhere(key, t.GetTableName(key), where, model, page, pageSize)
}
