package core

import (
	"context"

	"github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/model"
)

const (
	migrate_table_name = "ksql_migrate_info"
)

type migrateTable struct {
	*model.Model
	Id         int
	MigrateId  uint64
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
	return []any{&m.Id, &m.MigrateId, &m.Name, &m.Status, &m.Version, &m.CreateTime, &m.UpdateTime}
}

func (m *migrateTable) Columns() []string {
	return []string{"id", "migrate_id", "name", "status", "version", "create_time", "update_time"}
}

func (m *migrateTable) Save(ctx context.Context) error {
	return m.Model.Save(ctx, m)
}

func (m *migrateTable) Delete(ctx context.Context) error {
	return m.Model.Delete(ctx, m)
}
