package table

import (
	"github.com/kovey/db-go/v2/db"
	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/sql/meta"
)

type TableInterface[T itf.ModelInterface] interface {
	Database() db.DbInterface[T]
	InTransation(*db.Tx)
	Insert(meta.Data) (int64, error)
	Update(meta.Data, meta.Where) (int64, error)
	Delete(meta.Where) (int64, error)
	DeleteWhere(sql.WhereInterface) (int64, error)
	BatchInsert([]meta.Data) (int64, error)
	FetchRow(meta.Where, T) error
	LockRow(meta.Where, T) error
}

type Table[T itf.ModelInterface] struct {
	table string
	db    db.DbInterface[T]
}

func NewTable[T itf.ModelInterface](table string) *Table[T] {
	return NewTableByDb[T](table, db.NewMysql[T]())
}

func NewTableByDb[T itf.ModelInterface](table string, database db.DbInterface[T]) *Table[T] {
	return &Table[T]{db: database, table: table}
}

func (t *Table[T]) Database() db.DbInterface[T] {
	return t.db
}

func (t *Table[T]) InTransation(tx *db.Tx) {
	t.db.SetTx(tx)
}

func (t *Table[T]) Insert(data meta.Data) (int64, error) {
	in := sql.NewInsert(t.table)
	for field, value := range data {
		in.Set(field, value)
	}

	return t.db.Insert(in)
}

func (t *Table[T]) Update(data meta.Data, where meta.Where) (int64, error) {
	up := sql.NewUpdate(t.table)
	for field, value := range data {
		up.Set(field, value)
	}

	up.WhereByMap(where)

	return t.db.Update(up)
}

func (t *Table[T]) Delete(where meta.Where) (int64, error) {
	del := sql.NewDelete(t.table)
	del.WhereByMap(where)

	return t.db.Delete(del)
}

func (t *Table[T]) DeleteWhere(where sql.WhereInterface) (int64, error) {
	del := sql.NewDelete(t.table)
	del.Where(where)

	return t.db.Delete(del)
}

func (t *Table[T]) BatchInsert(data []meta.Data) (int64, error) {
	batch := sql.NewBatch(t.table)
	for _, val := range data {
		in := sql.NewInsert(t.table)
		for field, value := range val {
			in.Set(field, value)
		}
		batch.Add(in)
	}

	return t.db.BatchInsert(batch)
}

func (t *Table[T]) FetchRow(where meta.Where, model T) error {
	return t.db.FetchRow(t.table, where, model)
}

func (t *Table[T]) LockRow(where meta.Where, model T) error {
	return t.db.LockRow(t.table, where, model)
}

func (t *Table[T]) FetchAll(where meta.Where, model T) ([]T, error) {
	return t.db.FetchAll(t.table, where, model)
}

func (t *Table[T]) FetchAllByWhere(where sql.WhereInterface, model T) ([]T, error) {
	return t.db.FetchAllByWhere(t.table, where, model)
}

func (t *Table[T]) FetchBySelect(sel *sql.Select, model T) ([]T, error) {
	return t.db.FetchBySelect(sel, model)
}

func (t *Table[T]) FetchPage(where meta.Where, model T, page, pageSize int) (*meta.Page[T], error) {
	return t.db.FetchPage(t.table, where, model, page, pageSize)
}

func (t *Table[T]) FetchPageByWhere(where sql.WhereInterface, model T, page, pageSize int) (*meta.Page[T], error) {
	return t.db.FetchPageByWhere(t.table, where, model, page, pageSize)
}
