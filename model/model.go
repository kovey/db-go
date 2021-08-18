package model

import (
	"errors"
	"reflect"
	"strings"

	"github.com/kovey/db-go/table"
	"github.com/kovey/logger-go/logger"
)

type ModelInterface interface {
	Save(ModelInterface) error
}

type Base struct {
	table     table.TableInterface
	primaryId string
}

func NewBase(tb table.TableInterface, primaryId string) Base {
	return Base{table: tb, primaryId: primaryId}
}

func (b Base) Save(t ModelInterface) error {
	logger.Debug("b.save.table: %v", b.table)
	logger.Debug("primaryId: %s", b.primaryId)
	vt := reflect.ValueOf(t)

	if vt.Kind() != reflect.Ptr {
		return errors.New("params is not ptr")
	}

	vValue := vt.Elem()
	vType := vValue.Type()

	var id int64 = 0
	data := make(map[string]interface{})
	for i := 0; i < vValue.NumField(); i++ {
		tField := vType.Field(i)
		if tField.Name == "Base" {
			continue
		}

		vField := vValue.Field(i)
		data[strings.ToLower(tField.Name)] = vField.Interface()
		if tField.Name == b.primaryId {
			id = vField.Int()
		}
	}

	if id > 0 {
		where := make(map[string]interface{})
		where[strings.ToLower(b.primaryId)] = id

		_, err := b.table.Update(data, where)
		return err
	}

	var err error
	id, err = b.table.Insert(data)
	logger.Debug("save id: %d", id)
	if err == nil {
		vValue.FieldByName(b.primaryId).SetInt(id)
	}

	return err
}

func (b Base) Delete(t ModelInterface) error {
	data := make(map[string]interface{})
	vValue := reflect.ValueOf(t)

	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
	}

	data[strings.ToLower(b.primaryId)] = vValue.FieldByName(b.primaryId).Interface()

	_, err := b.table.Delete(data)
	return err
}

func (b Base) FetchRow(where map[string]interface{}, t ModelInterface) error {
	vt := reflect.ValueOf(t)
	if vt.Kind() != reflect.Ptr {
		return errors.New("params is not ptr")
	}
	row, err := b.table.FetchRow(where, t)
	if err != nil {
		return err
	}

	vValue := vt.Elem()
	tmp := reflect.New(vValue.Type()).Elem()
	tmp.Set(vValue)

	vValue.Set(reflect.ValueOf(row))
	vValue.FieldByName("Base").Set(tmp.FieldByName("Base"))

	return nil
}
