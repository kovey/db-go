package model

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/kovey/db-go/table"
	"github.com/kovey/logger-go/logger"
)

type ModelInterface interface {
	Save(ModelInterface) error
}

type Base struct {
	table     table.TableInterface
	primaryId *PrimaryId
	isInsert  bool
	err       error
}

func NewBase(tb table.TableInterface, primaryId *PrimaryId) Base {
	return Base{table: tb, primaryId: primaryId, isInsert: true}
}

func (b Base) Save(t ModelInterface) error {
	logger.Debug("b.save.table: %v", b.table)
	logger.Debug("primaryId: %s", b.primaryId.Name)
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

		_, err := b.table.Update(data, where)
		return err
	}

	logger.Debug("insert data: %v", data)
	id, err := b.table.Insert(data)
	logger.Debug("save id: %d", id)
	if err == nil && id > 0 {
		vValue.FieldByName(name).SetInt(id)
	}

	return err
}

func (b Base) Delete(t ModelInterface) error {
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

	_, err := b.table.Delete(data)
	return err
}

func (b Base) FetchRow(where map[string]interface{}, t ModelInterface) error {
	vt := reflect.ValueOf(t)
	if vt.Kind() != reflect.Ptr {
		return errors.New("params is not ptr")
	}
	row, err := b.table.FetchRow(where, t)
	b.err = err
	vValue := vt.Elem()
	vValue.FieldByName("Base").Set(reflect.ValueOf(b))
	if err != nil {
		return err
	}

	b.isInsert = false
	vValue.Set(reflect.ValueOf(row))

	return nil
}

func (b Base) Empty() bool {
	return b.err == sql.ErrNoRows
}
