package table

import (
	"fmt"

	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sharding"
	"github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/sql/meta"
)

type TableShardingInterface[T itf.ModelInterface] interface {
	InTransaction(tx *sharding.Tx)
	Database() *sharding.Mysql[T]
	Insert(any, meta.Data) (int64, error)
	Update(any, meta.Data, meta.Where) (int64, error)
	Delete(any, meta.Where) (int64, error)
	DeleteWhere(any, sql.WhereInterface) (int64, error)
	BatchInsert(any, []meta.Data) (int64, error)
	FetchRow(any, meta.Where, T) error
	LockRow(any, meta.Where, T) error
	FetchAll(any, meta.Where, T) ([]T, error)
	FetchAllByWhere(any, sql.WhereInterface, T) ([]T, error)
	FetchPage(any, meta.Where, T, int, int) ([]T, error)
	FetchPageByWhere(any, sql.WhereInterface, T, int, int) ([]T, error)
}

type TableSharding[T itf.ModelInterface] struct {
	table string
	db    *sharding.Mysql[T]
}

func NewTableSharding[T itf.ModelInterface](table string, isMaster bool) *TableSharding[T] {
	return &TableSharding[T]{db: sharding.NewMysql[T](isMaster), table: table}
}

func (t *TableSharding[T]) InTransaction(tx *sharding.Tx) {
	t.db.SetTx(tx)
}

func (t *TableSharding[T]) Database() *sharding.Mysql[T] {
	return t.db
}

func (t *TableSharding[T]) GetTableName(key any) string {
	return fmt.Sprintf("%s_%d", t.table, t.db.GetShardingKey(key))
}

func (t *TableSharding[T]) Insert(key any, data meta.Data) (int64, error) {
	in := sql.NewInsert(t.GetTableName(key))
	for field, value := range data {
		in.Set(field, value)
	}

	return t.db.Insert(key, in)
}

func (t *TableSharding[T]) Update(key any, data meta.Data, where meta.Where) (int64, error) {
	up := sql.NewUpdate(t.GetTableName(key))
	for field, value := range data {
		up.Set(field, value)
	}

	up.WhereByMap(where)

	return t.db.Update(key, up)
}

func (t *TableSharding[T]) Delete(key any, where meta.Where) (int64, error) {
	del := sql.NewDelete(t.GetTableName(key))
	del.WhereByMap(where)

	return t.db.Delete(key, del)
}

func (t *TableSharding[T]) DeleteWhere(key any, where sql.WhereInterface) (int64, error) {
	del := sql.NewDelete(t.GetTableName(key))
	del.Where(where)

	return t.db.Delete(key, del)
}

func (t *TableSharding[T]) BatchInsert(key any, data []meta.Data) (int64, error) {
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

func (t *TableSharding[T]) FetchRow(key any, where meta.Where, model T) error {
	return t.db.FetchRow(key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) LockRow(key any, where meta.Where, model T) error {
	return t.db.LockRow(key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) FetchAll(key any, where meta.Where, model T) ([]T, error) {
	return t.db.FetchAll(key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) FetchAllByWhere(key any, where sql.WhereInterface, model T) ([]T, error) {
	return t.db.FetchAllByWhere(key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) FetchPage(key any, where meta.Where, model T, page, pageSize int) ([]T, error) {
	return t.db.FetchPage(key, t.GetTableName(key), where, model, page, pageSize)
}

func (t *TableSharding[T]) FetchPageByWhere(key any, where sql.WhereInterface, model T, page, pageSize int) ([]T, error) {
	return t.db.FetchPageByWhere(key, t.GetTableName(key), where, model, page, pageSize)
}
