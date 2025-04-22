package sharding

import (
	"fmt"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/db-go/v3/model"
)

type Model struct {
	*model.Model
	key any
}

func NewModel(table, primaryId string, t model.PrimaryType) *Model {
	return &Model{Model: model.NewModel(table, primaryId, t)}
}

func (m *Model) Key() any {
	return m.key
}

func (m *Model) WithKey(key any) {
	m.key = key
	m.SetTable(fmt.Sprintf("%s_%d", m.Model.Table(), database.node(m.key)))
}

func Rows[T ModelInterface](key any, models *[]T) ksql.BuilderInterface[T] {
	var m T
	tmp := m.Clone().(T)
	tmp.WithKey(key)
	return db.ShardingModels(tmp.Table(), models).WithConn(database.conn(key))
}

func Row[T ModelInterface](key any, model T) ksql.BuilderInterface[T] {
	return RowBy(key, model, database)
}

func RowBy[T ModelInterface](key any, model T, conn ConnectionInterface) ksql.BuilderInterface[T] {
	model.WithKey(key)
	return db.Model(model).WithConn(conn.Get(key))
}
