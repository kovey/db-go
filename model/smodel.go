package model

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/kovey/db-go/table"
)

type ModelShardingInterface interface {
	Save(interface{}, ModelShardingInterface) error
}

type BaseSharding struct {
	table     table.TableShardingInterface
	primaryId *PrimaryId
	isInsert  bool
	err       error
}

func NewBaseSharding(tb table.TableShardingInterface, primaryId *PrimaryId) BaseSharding {
	return BaseSharding{table: tb, primaryId: primaryId, isInsert: true}
}

func (b BaseSharding) Save(key interface{}, t ModelShardingInterface) error {
	vt := reflect.ValueOf(t)

	if vt.Kind() != reflect.Ptr {
		return errors.New("params is not ptr")
	}

	vValue := vt.Elem()
	vType := vValue.Type()

	var name string
	data := make(map[string]interface{})
	for i := 0; i < vValue.NumField(); i++ {
		tField := vType.Field(i)
		tag := tField.Tag.Get("db")
		if len(tag) == 0 {
			continue
		}

		vField := vValue.Field(i)
		if tag == b.primaryId.Name {
			b.primaryId.Parse(vField)
			name = tField.Name
			if b.primaryId.Null() {
				continue
			}
		}

		data[tag] = vField.Interface()
	}

	if !b.isInsert {
		where := make(map[string]interface{})
		where[b.primaryId.Name] = b.primaryId.Value()

		_, err := b.table.Update(key, data, where)
		return err
	}

	var err error
	id, err := b.table.Insert(key, data)
	if err == nil && id > 0 {
		vValue.FieldByName(name).SetInt(id)
	}

	return err
}

func (b BaseSharding) Delete(key interface{}, t ModelShardingInterface) error {
	data := make(map[string]interface{})
	vValue := reflect.ValueOf(t)

	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
	}

	var name string
	vType := vValue.Type()
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		if field.Tag.Get("db") == b.primaryId.Name {
			name = field.Name
			break
		}
	}

	data[b.primaryId.Name] = vValue.FieldByName(name).Interface()

	_, err := b.table.Delete(key, data)
	return err
}

func (b BaseSharding) FetchRow(key interface{}, where map[string]interface{}, t ModelShardingInterface) error {
	vt := reflect.ValueOf(t)
	if vt.Kind() != reflect.Ptr {
		return errors.New("params is not ptr")
	}
	row, err := b.table.FetchRow(key, where, t)
	b.err = err
	vValue := vt.Elem()
	if err != nil {
		vValue.FieldByName("BaseSharding").Set(reflect.ValueOf(b))
		return err
	}

	b.isInsert = false
	vValue.Set(reflect.ValueOf(row))
	vValue.FieldByName("BaseSharding").Set(reflect.ValueOf(b))

	return nil
}

func (b BaseSharding) Empty() bool {
	return b.err == sql.ErrNoRows
}
