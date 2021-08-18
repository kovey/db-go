package table

import (
	"github.com/kovey/db-go/db"
	"github.com/kovey/db-go/sql"
)

type TableInterface interface {
	Insert(map[string]interface{}) (int64, error)
	Update(map[string]interface{}, map[string]interface{}) (int64, error)
	Delete(map[string]interface{}) (int64, error)
	BatchInsert([]map[string]interface{}) (int64, error)
	FetchRow(map[string]interface{}, interface{}) (interface{}, error)
	FetchAll(map[string]interface{}, interface{}) ([]interface{}, error)
}

type Table struct {
	table string
	db    *db.Mysql
}

func NewTable(table string) *Table {
	return &Table{db: db.NewMysql(), table: table}
}

func (t *Table) Insert(data map[string]interface{}) (int64, error) {
	in := sql.NewInsert(t.table)
	for field, value := range data {
		in.Set(field, value)
	}

	return t.db.Insert(in)
}

func (t *Table) Update(data map[string]interface{}, where map[string]interface{}) (int64, error) {
	up := sql.NewUpdate(t.table)
	for field, value := range data {
		up.Set(field, value)
	}

	up.WhereByMap(where)

	return t.db.Update(up)
}

func (t *Table) Delete(where map[string]interface{}) (int64, error) {
	del := sql.NewDelete(t.table)
	del.WhereByMap(where)

	return t.db.Delete(del)
}

func (t *Table) BatchInsert(data []map[string]interface{}) (int64, error) {
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

func (t *Table) FetchRow(where map[string]interface{}, mt interface{}) (interface{}, error) {
	return t.db.FetchRow(t.table, where, mt)
}

func (t *Table) FetchAll(where map[string]interface{}, mt interface{}) ([]interface{}, error) {
	return t.db.FetchAll(t.table, where, mt)
}
