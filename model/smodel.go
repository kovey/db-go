package model

import (
	"errors"
	"reflect"

	"github.com/kovey/db-go/table"
	"github.com/kovey/logger-go/logger"
)

type ModelShardingInterface interface {
	Save(interface{}, ModelShardingInterface) error
}

type BaseSharding struct {
	table     table.TableShardingInterface
	primaryId string
	isInsert  bool
}

func NewBaseSharding(tb table.TableShardingInterface, primaryId string) BaseSharding {
	return BaseSharding{table: tb, primaryId: primaryId, isInsert: true}
}

func (b BaseSharding) Save(key interface{}, t ModelShardingInterface) error {
	logger.Debug("b.save.table: %v", b.table)
	logger.Debug("primaryId: %s", b.primaryId)
	vt := reflect.ValueOf(t)

	if vt.Kind() != reflect.Ptr {
		return errors.New("params is not ptr")
	}

	vValue := vt.Elem()
	vType := vValue.Type()

	var id int64 = 0
	var name string
	data := make(map[string]interface{})
	for i := 0; i < vValue.NumField(); i++ {
		tField := vType.Field(i)
		tag := tField.Tag.Get("db")
		if len(tag) == 0 {
			continue
		}

		vField := vValue.Field(i)
		if tag == b.primaryId {
			id = vField.Int()
			name = tField.Name
			if id == 0 {
				continue
			}
		}

		data[tag] = vField.Interface()
	}

	if !b.isInsert {
		where := make(map[string]interface{})
		where[b.primaryId] = id

		_, err := b.table.Update(key, data, where)
		return err
	}

	var err error
	id, err = b.table.Insert(key, data)
	logger.Debug("save id: %d", id)
	if err == nil {
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
		if field.Tag.Get("db") == b.primaryId {
			name = field.Name
			break
		}
	}

	data[b.primaryId] = vValue.FieldByName(name).Interface()

	_, err := b.table.Delete(key, data)
	return err
}

func (b BaseSharding) FetchRow(key interface{}, where map[string]interface{}, t ModelShardingInterface) error {
	vt := reflect.ValueOf(t)
	if vt.Kind() != reflect.Ptr {
		return errors.New("params is not ptr")
	}
	row, err := b.table.FetchRow(key, where, t)
	if err != nil {
		return err
	}

	b.isInsert = false
	vValue := vt.Elem()
	vValue.Set(reflect.ValueOf(row))
	vValue.FieldByName("BaseSharding").Set(reflect.ValueOf(b))

	return nil
}
