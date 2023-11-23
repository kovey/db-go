package table

import (
	"context"
	"fmt"

	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sharding"
	"github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/pool/object"
)

type TableShardingInterface[T itf.ModelInterface] interface {
	InTransaction(tx *sharding.Tx)
	Database() *sharding.Mysql[T]
	Insert(any, meta.Data) (int64, error)
	Update(any, meta.Data, meta.Where) (int64, error)
	UpdateWhere(any, meta.Data, sql.WhereInterface) (int64, error)
	Delete(any, meta.Where) (int64, error)
	DeleteWhere(any, sql.WhereInterface) (int64, error)
	BatchInsert(any, []meta.Data) (int64, error)
	FetchRow(any, meta.Where, T) error
	LockRow(any, meta.Where, T) error
	InsertCtx(context.Context, any, meta.Data) (int64, error)
	UpdateCtx(context.Context, any, meta.Data, meta.Where) (int64, error)
	UpdateWhereCtx(context.Context, any, meta.Data, sql.WhereInterface) (int64, error)
	DeleteCtx(context.Context, any, meta.Where) (int64, error)
	DeleteWhereCtx(context.Context, any, sql.WhereInterface) (int64, error)
	BatchInsertCtx(context.Context, any, []meta.Data) (int64, error)
	FetchRowCtx(context.Context, any, meta.Where, T) error
	LockRowCtx(context.Context, any, meta.Where, T) error
}

type TableSharding[T itf.ModelInterface] struct {
	table string
	db    *sharding.Mysql[T]
}

func NewTableSharding[T itf.ModelInterface](table string) *TableSharding[T] {
	return &TableSharding[T]{db: sharding.NewMysql[T](), table: table}
}

func NewTableShardingBy[T itf.ModelInterface](table string, isMaster bool) *TableSharding[T] {
	return &TableSharding[T]{db: sharding.NewMysqlBy[T](isMaster), table: table}
}

func (t *TableSharding[T]) Set(isMaster bool, tx *sharding.Tx) {
	t.db.Set(isMaster, tx)
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
	return t.InsertCtx(context.Background(), key, data)
}

func (t *TableSharding[T]) Update(key any, data meta.Data, where meta.Where) (int64, error) {
	return t.UpdateCtx(context.Background(), key, data, where)
}

func (t *TableSharding[T]) UpdateWhere(key any, data meta.Data, where sql.WhereInterface) (int64, error) {
	return t.UpdateWhereCtx(context.Background(), key, data, where)
}

func (t *TableSharding[T]) Delete(key any, where meta.Where) (int64, error) {
	return t.DeleteCtx(context.Background(), key, where)
}

func (t *TableSharding[T]) DeleteWhere(key any, where sql.WhereInterface) (int64, error) {
	return t.DeleteWhereCtx(context.Background(), key, where)
}

func (t *TableSharding[T]) BatchInsert(key any, data []meta.Data) (int64, error) {
	return t.BatchInsertCtx(context.Background(), key, data)
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

func (t *TableSharding[T]) FetchPage(key any, where meta.Where, model T, page, pageSize int, orders ...string) (*meta.Page[T], error) {
	return t.db.FetchPage(key, t.GetTableName(key), where, model, page, pageSize, orders...)
}

func (t *TableSharding[T]) FetchPageByWhere(key any, where sql.WhereInterface, model T, page, pageSize int, orders ...string) (*meta.Page[T], error) {
	return t.db.FetchPageByWhere(key, t.GetTableName(key), where, model, page, pageSize, orders...)
}

func (t *TableSharding[T]) InsertCtx(ctx context.Context, key any, data meta.Data) (int64, error) {
	var in *sql.Insert
	if cc, ok := ctx.(object.CtxInterface); ok {
		in = sql.NewInsertBy(cc, t.GetTableName(key))
	} else {
		in = sql.NewInsert(t.GetTableName(key))
	}

	for field, value := range data {
		in.Set(field, value)
	}

	return t.db.InsertCtx(ctx, key, in)
}

func (t *TableSharding[T]) UpdateCtx(ctx context.Context, key any, data meta.Data, where meta.Where) (int64, error) {
	var up *sql.Update
	if cc, ok := ctx.(object.CtxInterface); ok {
		up = sql.NewUpdateBy(cc, t.GetTableName(key))
	} else {
		up = sql.NewUpdate(t.GetTableName(key))
	}

	for field, value := range data {
		up.Set(field, value)
	}

	up.WhereByMap(where)

	return t.db.UpdateCtx(ctx, key, up)
}

func (t *TableSharding[T]) UpdateWhereCtx(ctx context.Context, key any, data meta.Data, where sql.WhereInterface) (int64, error) {
	var up *sql.Update
	if cc, ok := ctx.(object.CtxInterface); ok {
		up = sql.NewUpdateBy(cc, t.GetTableName(key))
	} else {
		up = sql.NewUpdate(t.GetTableName(key))
	}

	for field, value := range data {
		up.Set(field, value)
	}

	up.Where(where)
	return t.db.UpdateCtx(ctx, key, up)
}

func (t *TableSharding[T]) DeleteCtx(ctx context.Context, key any, where meta.Where) (int64, error) {
	var del *sql.Delete
	if cc, ok := ctx.(object.CtxInterface); ok {
		del = sql.NewDeleteBy(cc, t.GetTableName(key))
	} else {
		del = sql.NewDelete(t.GetTableName(key))
	}
	del.WhereByMap(where)

	return t.db.DeleteCtx(ctx, key, del)
}

func (t *TableSharding[T]) DeleteWhereCtx(ctx context.Context, key any, where sql.WhereInterface) (int64, error) {
	var del *sql.Delete
	if cc, ok := ctx.(object.CtxInterface); ok {
		del = sql.NewDeleteBy(cc, t.GetTableName(key))
	} else {
		del = sql.NewDelete(t.GetTableName(key))
	}
	del.Where(where)

	return t.db.DeleteCtx(ctx, key, del)
}

func (t *TableSharding[T]) BatchInsertCtx(ctx context.Context, key any, data []meta.Data) (int64, error) {
	var batch *sql.Batch
	if cc, ok := ctx.(object.CtxInterface); ok {
		batch = sql.NewBatchBy(cc, t.GetTableName(key))
	} else {
		batch = sql.NewBatch(t.GetTableName(key))
	}
	for _, val := range data {
		in := sql.NewInsert(t.GetTableName(key))
		for field, value := range val {
			in.Set(field, value)
		}
		batch.Add(in)
	}

	return t.db.BatchInsertCtx(ctx, key, batch)
}

func (t *TableSharding[T]) FetchRowCtx(ctx context.Context, key any, where meta.Where, model T) error {
	return t.db.FetchRowCtx(ctx, key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) LockRowCtx(ctx context.Context, key any, where meta.Where, model T) error {
	return t.db.LockRowCtx(ctx, key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) FetchAllCtx(ctx context.Context, key any, where meta.Where, model T) ([]T, error) {
	return t.db.FetchAllCtx(ctx, key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) FetchAllByWhereCtx(ctx context.Context, key any, where sql.WhereInterface, model T) ([]T, error) {
	return t.db.FetchAllByWhereCtx(ctx, key, t.GetTableName(key), where, model)
}

func (t *TableSharding[T]) FetchPageCtx(ctx context.Context, key any, where meta.Where, model T, page, pageSize int, orders ...string) (*meta.Page[T], error) {
	return t.db.FetchPageCtx(ctx, key, t.GetTableName(key), where, model, page, pageSize, orders...)
}

func (t *TableSharding[T]) FetchPageByWhereCtx(ctx context.Context, key any, where sql.WhereInterface, model T, page, pageSize int, orders ...string) (*meta.Page[T], error) {
	return t.db.FetchPageByWhereCtx(ctx, key, t.GetTableName(key), where, model, page, pageSize, orders...)
}
