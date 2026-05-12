// Package core
package core

import (
	"context"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/model"
)

const (
	migrateTableName = "ksql_migrate_info"
)

type migrateTable struct {
	*model.Model
	ID         int
	MigrateID  uint64
	Name       string
	Status     int
	Version    string
	CreateTime string
	UpdateTime string
}

func newMigrateTable() *migrateTable {
	return &migrateTable{Model: model.NewModel("ksql_migrate_info", "id", model.Type_Int)}
}

func (m *migrateTable) Clone() ksql.RowInterface {
	return newMigrateTable()
}

func (m *migrateTable) Values() []any {
	return []any{&m.ID, &m.MigrateID, &m.Name, &m.Status, &m.Version, &m.CreateTime, &m.UpdateTime}
}

func (m *migrateTable) Columns() []string {
	return []string{"id", "migrate_id", "name", "status", "version", "create_time", "update_time"}
}

func (m *migrateTable) Save(ctx context.Context) error {
	return m.SaveBy(ctx, m)
}

func (m *migrateTable) Delete(ctx context.Context) error {
	return m.DeleteBy(ctx, m)
}
