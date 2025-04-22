package sharding

import (
	"fmt"

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
