package model

import (
	"context"
	"errors"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

var Err_Affect_No_Rows = errors.New("affect no rows")

type PrimaryType byte

const (
	Type_Int PrimaryType = 1
	Type_Str PrimaryType = 2
)

type Model struct {
	table         string
	conn          ksql.ConnectionInterface
	primaryId     string
	primaryType   PrimaryType
	isAutoInc     bool
	fromFecth     bool
	isInitialized bool
	data          *db.Data
	shardingType  ksql.Sharding
}

func NewModel(table, primaryId string, t PrimaryType) *Model {
	return &Model{table: table, primaryId: primaryId, primaryType: t, isAutoInc: true, isInitialized: false, data: db.NewData(), shardingType: ksql.Sharding_None}
}

func (m *Model) OnUpdateBefore(conn ksql.ConnectionInterface) error { return nil }
func (m *Model) OnUpdateAfter(conn ksql.ConnectionInterface) error  { return nil }
func (m *Model) OnCreateBefore(conn ksql.ConnectionInterface) error { return nil }
func (m *Model) OnCreateAfter(conn ksql.ConnectionInterface) error  { return nil }
func (m *Model) OnDeleteBefore(conn ksql.ConnectionInterface) error { return nil }
func (m *Model) OnDeleteAfter(conn ksql.ConnectionInterface) error  { return nil }

func (m *Model) Empty() bool {
	return !m.isInitialized
}

func (m *Model) Scan(s ksql.ScanInterface, r ksql.RowInterface) error {
	if err := s.Scan(r.Values()...); err != nil {
		return err
	}

	m.isInitialized = true
	m.fromFecth = true
	values := r.Values()
	tmp, ok := r.(ksql.ModelInterface)
	if !ok {
		return nil
	}

	for i, column := range tmp.Columns() {
		m.data.Set(column, values[i])
	}

	return nil
}

func (m *Model) Sharding(sharding ksql.Sharding) {
	m.shardingType = sharding
}

func (m *Model) Table() string {
	return ksql.FormatSharding(m.table, m.shardingType)
}

func (m *Model) WithConn(conn ksql.ConnectionInterface) {
	m.conn = conn
}

func (m *Model) Conn() ksql.ConnectionInterface {
	return m.conn
}

func (m *Model) NoAutoInc() {
	m.isAutoInc = false
}

func (m *Model) PrimaryId() string {
	return m.primaryId
}

func (m *Model) setPrimary(model ksql.ModelInterface, id int64) {
	if !m.isAutoInc || m.primaryType != Type_Int {
		return
	}

	for i, column := range model.Columns() {
		if column == m.primaryId {
			val := model.Values()[i]
			switch tmp := val.(type) {
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
			case *uint:
				*tmp = uint(id)
			case *uint8:
				*tmp = uint8(id)
			case *uint16:
				*tmp = uint16(id)
			case *uint32:
				*tmp = uint32(id)
			case *uint64:
				*tmp = uint64(id)
			}
		}
	}
}

func (m *Model) hasChanged(model ksql.ModelInterface) bool {
	columns := model.Columns()
	values := model.Values()
	for i, column := range columns {
		if m.data.Changed(column, values[i]) {
			return true
		}
	}

	return false
}

func (m *Model) toData(model ksql.ModelInterface) *db.Data {
	data := db.NewData()
	columns := model.Columns()
	values := model.Values()
	for i, column := range columns {
		if !m.data.Changed(column, values[i]) {
			continue
		}

		data.Set(column, values[i])
	}

	return data
}

func (m *Model) insert(ctx context.Context, data *db.Data) (int64, error) {
	if m.conn == nil {
		return db.Insert(ctx, m.table, data)
	}

	op := db.NewInsert()
	op.Table(m.table)
	data.Range(func(key string, val any) {
		op.Add(key, val)
	})

	return m.conn.Insert(ctx, op)
}

func (m *Model) update(ctx context.Context, data *db.Data) (int64, error) {
	w := db.NewWhere()
	w.Where(m.primaryId, "=", data.Get(m.primaryId))
	if m.conn == nil {
		return db.Update(ctx, m.table, data, w)
	}

	u := db.NewUpdate()
	u.Table(m.Table())
	data.Range(func(key string, val any) {
		if key == m.primaryId {
			return
		}

		u.Set(key, val)
	})
	u.Where(w)
	return m.conn.Update(ctx, u)
}

func (m *Model) _conn() ksql.ConnectionInterface {
	if m.conn == nil {
		database, _ := db.Get()
		return database
	}

	return m.conn
}

func (m *Model) Save(ctx context.Context, model ksql.ModelInterface) error {
	if !m.hasChanged(model) {
		return nil
	}

	if !m.fromFecth {
		if err := m.OnCreateBefore(m._conn()); err != nil {
			return err
		}
	} else {
		if err := m.OnUpdateBefore(m._conn()); err != nil {
			return err
		}
	}

	data := m.toData(model)
	if !m.fromFecth {
		id, err := m.insert(ctx, data)
		if err != nil {
			return err
		}

		m.setPrimary(model, id)
		return m.OnCreateAfter(m._conn())
	}

	id, err := m.update(ctx, data)
	if err != nil {
		return err
	}

	if id == 0 {
		return Err_Affect_No_Rows
	}

	return m.OnUpdateAfter(m._conn())
}

func (m *Model) primaryValue(model ksql.ModelInterface) any {
	columns := model.Columns()
	for i, val := range model.Values() {
		if columns[i] == m.primaryId {
			return val
		}
	}

	return nil
}

func (m *Model) Delete(ctx context.Context, model ksql.ModelInterface) error {
	if err := m.OnDeleteBefore(m._conn()); err != nil {
		return err
	}

	w := db.NewWhere()
	w.Where(m.primaryId, "=", m.primaryValue(model))
	if m.conn == nil {
		id, err := db.Delete(ctx, m.table, w)
		if err != nil {
			return err
		}

		if id == 0 {
			return Err_Affect_No_Rows
		}

		return nil
	}

	op := db.NewDelete()
	op.Table(m.table).Where(w)
	id, err := m.conn.Delete(ctx, op)
	if err != nil {
		return err
	}

	if id == 0 {
		return Err_Affect_No_Rows
	}

	return m.OnDeleteAfter(m._conn())
}

func Rows[T ksql.ModelInterface](models *[]T) ksql.BuilderInterface[T] {
	return db.Models(models)
}

func Row[T ksql.ModelInterface](model T) ksql.BuilderInterface[T] {
	return db.Model(model)
}
