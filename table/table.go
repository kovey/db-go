package table

import (
	"github.com/kovey/db-go/db"
	"github.com/kovey/db-go/sql"
)

type TableInterface[T any] interface {
	Database() db.DbInterface[T]
	Insert(map[string]any) (int64, error)
	Update(map[string]any, map[string]any) (int64, error)
	Delete(map[string]any) (int64, error)
	DeleteWhere(*sql.Where) (int64, error)
	BatchInsert([]map[string]any) (int64, error)
	FetchRow(map[string]any, T) (T, error)
	FetchAll(map[string]any, T) ([]T, error)
	FetchAllByWhere(*sql.Where, T) ([]T, error)
	FetchPage(map[string]any, T, int, int) ([]T, error)
	FetchPageByWhere(*sql.Where, T, int, int) ([]T, error)
}

type Table[T any] struct {
	table string
	db    db.DbInterface[T]
}

func NewTable[T any](table string) *Table[T] {
	return NewTableByDb[T](table, db.NewMysql[T]())
}

func NewTableByDb[T any](table string, database db.DbInterface[T]) *Table[T] {
	return &Table[T]{db: database, table: table}
}

func (t *Table[T]) Database() db.DbInterface[T] {
	return t.db
}

func (t *Table[T]) Insert(data map[string]any) (int64, error) {
	in := sql.NewInsert(t.table)
	for field, value := range data {
		in.Set(field, value)
	}

	return t.db.Insert(in)
}

func (t *Table[T]) Update(data map[string]any, where map[string]any) (int64, error) {
	up := sql.NewUpdate(t.table)
	for field, value := range data {
		up.Set(field, value)
	}

	up.WhereByMap(where)

	return t.db.Update(up)
}

func (t *Table[T]) Delete(where map[string]any) (int64, error) {
	del := sql.NewDelete(t.table)
	del.WhereByMap(where)

	return t.db.Delete(del)
}

func (t *Table[T]) DeleteWhere(where *sql.Where) (int64, error) {
	del := sql.NewDelete(t.table)
	del.Where(where)

	return t.db.Delete(del)
}

func (t *Table[T]) BatchInsert(data []map[string]any) (int64, error) {
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

func (t *Table[T]) FetchRow(where map[string]any, model T) (T, error) {
	return t.db.FetchRow(t.table, where, model)
}

func (t *Table[T]) FetchAll(where map[string]any, model T) ([]T, error) {
	return t.db.FetchAll(t.table, where, model)
}

func (t *Table[T]) FetchAllByWhere(where *sql.Where, model T) ([]T, error) {
	return t.db.FetchAllByWhere(t.table, where, model)
}

func (t *Table[T]) FetchPage(where map[string]any, model T, page, pageSize int) ([]T, error) {
	return t.db.FetchPage(t.table, where, model, page, pageSize)
}

func (t *Table[T]) FetchPageByWhere(where *sql.Where, model T, page, pageSize int) ([]T, error) {
	return t.db.FetchPageByWhere(t.table, where, model, page, pageSize)
}
