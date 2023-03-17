package model

import (
	"database/sql"
	"reflect"

	"github.com/kovey/db-go/v2/rows"
	"github.com/kovey/db-go/v2/table"
)

type ModelShardingInterface interface {
	Empty() bool
}

type BaseSharding[T ModelShardingInterface] struct {
	Table     table.TableShardingInterface[T]
	primaryId *PrimaryId
	isInsert  bool
	err       error
}

func NewBaseSharding[T ModelShardingInterface](tb table.TableShardingInterface[T], primaryId *PrimaryId) *BaseSharding[T] {
	return &BaseSharding[T]{Table: tb, primaryId: primaryId, isInsert: true}
}

func (b *BaseSharding[T]) Save(key any, model T) error {
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
		_, err := b.Table.Update(key, data, where)
		return err
	}

	id, err := b.Table.Insert(key, data)
	if err != nil {
		return err
	}

	if err == nil && id > 0 {
		vValue.FieldByName(name).SetInt(id)
	}
	return nil
}

func (b *BaseSharding[T]) Delete(key any, model T) error {
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
	_, err := b.Table.Delete(key, where)
	return err
}

func (b *BaseSharding[T]) FetchRow(key any, where map[string]any, model T) error {
	vValue := reflect.ValueOf(model)
	isPointer := false
	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
		isPointer = true
	}

	row, err := b.Table.FetchRow(key, where, model)
	b.err = err
	if err != nil {
		vValue.FieldByName("BaseSharding").Set(reflect.ValueOf(b))
		return err
	}

	b.isInsert = false
	if isPointer {
		vValue.Set(reflect.ValueOf(row).Elem())
	} else {
		vValue.Set(reflect.ValueOf(row))
	}

	vValue.FieldByName("BaseSharding").Set(reflect.ValueOf(b))
	return nil
}

func (b *BaseSharding[T]) Empty() bool {
	return b.err == sql.ErrNoRows
}
