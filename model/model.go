package model

import (
	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/db-go/v2/table"
)

type Base[T itf.ModelInterface] struct {
	Table     table.TableInterface[T]
	primaryId *PrimaryId
	isInsert  bool
	isEmpty   bool
}

func NewBase[T itf.ModelInterface](tb table.TableInterface[T], primaryId *PrimaryId) *Base[T] {
	return &Base[T]{Table: tb, primaryId: primaryId, isInsert: true, isEmpty: false}
}

func (b *Base[T]) Save(model T) error {
	columns := model.Columns()
	fields := model.Fields()
	values := model.Values()
	data := meta.NewData()
	var primary any
	for index, column := range columns {
		if column.Name.Name == b.primaryId.Name {
			b.primaryId.Parse(values[index])
			if b.primaryId.Null() {
				continue
			}

			primary = fields[index]
			continue
		}

		data.Add(column.Name.Name, values[index])
	}

	if !b.isInsert {
		where := meta.NewWhere()
		where.Add(b.primaryId.Name, b.primaryId.Value())
		_, err := b.Table.Update(data, where)
		return err
	}

	id, err := b.Table.Insert(data)
	if err != nil {
		return err
	}

	if id <= 0 {
		return nil
	}

	switch tmp := primary.(type) {
	case *int:
		*tmp = int(id)
	case *int8:
		*tmp = int8(id)
	case *int16:
		*tmp = int16(id)
	case *int32:
		*tmp = int32(id)
	case *int64:
		*tmp = int64(id)
	case *uint8:
		*tmp = uint8(id)
	case *uint16:
		*tmp = uint16(id)
	case *uint32:
		*tmp = uint32(id)
	case *uint64:
		*tmp = uint64(id)
	}

	return nil
}

func (b *Base[T]) Delete(model T) error {
	columns := model.Columns()
	values := model.Values()
	where := meta.NewWhere()

	for index, column := range columns {
		if column.Name.Name == b.primaryId.Name {
			where.Add(b.primaryId.Name, values[index])
			break
		}
	}

	_, err := b.Table.Delete(where)
	return err
}

func (b *Base[T]) FetchRow(where meta.Where, model T) error {
	b.isEmpty = false
	err := b.Table.FetchRow(where, model)
	if err != nil {
		return err
	}

	b.isInsert = false
	return nil
}

func (b *Base[T]) LockRow(where meta.Where, model T) error {
	b.isEmpty = false
	err := b.Table.LockRow(where, model)
	if err != nil {
		return err
	}

	b.isInsert = false
	return nil
}

func (b *Base[T]) Empty() bool {
	return b.isEmpty
}

func (b *Base[T]) SetEmpty() {
	b.isEmpty = true
}
