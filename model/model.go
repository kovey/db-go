package model

import (
	"database/sql"
	"reflect"

	"github.com/kovey/db-go/v2/rows"
	"github.com/kovey/db-go/v2/table"
)

type ModelInterface interface {
	Empty() bool
}

type Base[T ModelInterface] struct {
	Table     table.TableInterface[T]
	primaryId *PrimaryId
	isInsert  bool
	err       error
}

func NewBase[T ModelInterface](tb table.TableInterface[T], primaryId *PrimaryId) *Base[T] {
	return &Base[T]{Table: tb, primaryId: primaryId, isInsert: true}
}

func (b *Base[T]) Save(model T) error {
	vValue := reflect.ValueOf(model)
	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
	}

	vType := vValue.Type()
	var name string

	data := make(map[string]any)
	for i := 0; i < vValue.NumField(); i++ {
		tField := vType.Field(i)
		tag := tField.Tag.Get(rows.Tag_Db)
		if tag == "" {
			continue
		}

		vField := vValue.Field(i)
		if tag == b.primaryId.Name {
			name = tField.Name
			b.primaryId.Parse(vField)
			if b.primaryId.Null() {
				continue
			}

			continue
		}

		data[tag] = vField.Interface()
	}

	if !b.isInsert {
		where := make(map[string]any)
		where[b.primaryId.Name] = b.primaryId.Value()
		_, err := b.Table.Update(data, where)
		return err
	}

	id, err := b.Table.Insert(data)
	if err != nil {
		return err
	}

	if err == nil && id > 0 {
		vValue.FieldByName(name).SetInt(id)
	}
	return nil
}

func (b *Base[T]) Delete(model T) error {
	where := make(map[string]any)
	vValue := reflect.ValueOf(model)
	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
	}

	vType := vValue.Type()
	var val any

	for i := 0; i < vType.NumField(); i++ {
		tField := vType.Field(i)
		if tField.Tag.Get(rows.Tag_Db) == b.primaryId.Name {
			val = vValue.Field(i).Interface()
			break
		}
	}

	where[b.primaryId.Name] = val
	_, err := b.Table.Delete(where)
	return err
}

func (b *Base[T]) FetchRow(where map[string]any, model T) error {
	vValue := reflect.ValueOf(model)
	isPointer := false
	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
		isPointer = true
	}

	row, err := b.Table.FetchRow(where, model)
	b.err = err
	if err != nil {
		vValue.FieldByName("Base").Set(reflect.ValueOf(b))
		return err
	}

	b.isInsert = false
	if isPointer {
		vValue.Set(reflect.ValueOf(row).Elem())
	} else {
		vValue.Set(reflect.ValueOf(row))
	}

	vValue.FieldByName("Base").Set(reflect.ValueOf(b))
	return nil
}

func (b *Base[T]) Empty() bool {
	return b.err == sql.ErrNoRows
}
