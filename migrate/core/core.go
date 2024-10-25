package core

import (
	"context"
	"fmt"
	"plugin"
	"time"

	"github.com/kovey/cli-go/gui"
	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/db-go/v3/migplug"
	"github.com/kovey/db-go/v3/model"
)

type migrates struct {
	data map[uint64]migplug.MigrateInterface
	keys []uint64
}

func newMigrates() *migrates {
	return &migrates{data: make(map[uint64]migplug.MigrateInterface)}
}

func (m *migrates) Len() int {
	return len(m.keys)
}

func (m *migrates) Swap(i, j int) {
	m.keys[i], m.keys[j] = m.keys[j], m.keys[i]
}

func (m *migrates) Less(i, j int) bool {
	return m.keys[i] < m.keys[j]
}

func (m *migrates) Add(mi migplug.MigrateInterface) {
	if _, ok := m.data[mi.Id()]; ok {
		return
	}

	m.data[mi.Id()] = mi
	m.keys = append(m.keys, mi.Id())
}

func (m *migrates) Range(call func(key uint64, mi migplug.MigrateInterface)) {
	for _, key := range m.keys {
		call(key, m.data[key])
	}
}

func (m *migrates) Get(key uint64) migplug.MigrateInterface {
	return m.data[key]
}

var mig = newMigrates()

func AddMigrate(mi migplug.MigrateInterface) {
	mig.Add(mi)
}

func _load(driverName, dsn, path string) error {
	plug, err := plugin.Open(path)
	if err != nil {
		return err
	}

	fun, err := plug.Lookup("Migrate")
	if err != nil {
		return err
	}

	f, ok := fun.(func() migplug.PluginInterface)
	if !ok {
		return fmt.Errorf("plugin[%s] not core.PluginInterface", path)
	}

	migrate := f()
	if err := db.Init(db.Config{DriverName: driverName, DataSourceName: dsn, MaxIdleTime: 120 * time.Second, MaxLifeTime: 120 * time.Second, MaxIdleConns: 10, MaxOpenConns: 10}); err != nil {
		return err
	}

	migrate.Register(mig)
	return nil
}

func LoadPlugin(driverName, dsn, path string, t MigrateType) error {
	if err := _load(driverName, dsn, path); err != nil {
		return err
	}

	Migrate(context.Background(), t)
	return nil
}

func Show(driverName, dsn, path string) error {
	if err := _load(driverName, dsn, path); err != nil {
		return err
	}

	ctx := context.Background()
	if err := check(ctx); err != nil {
		return err
	}

	var ids []any
	mig.Range(func(key uint64, mi migplug.MigrateInterface) {
		ids = append(ids, mi.Id())
	})

	table := gui.NewTable()
	if len(ids) == 0 {
		table.Add("No migrations")
		table.Show()
		return nil
	}

	var migrations []*migrateTable
	if err := model.Query(newMigrateTable()).WhereIn("migrate_id", ids).All(ctx, &migrations); err != nil {
		return err
	}

	mig.Range(func(key uint64, mi migplug.MigrateInterface) {
		for _, migration := range migrations {
			if migration.MigrateId == key {
				table.Add(fmt.Sprintf("%s: Yes", mi.Name()))
				return
			}
		}

		table.Add(fmt.Sprintf("%s: No", mi.Name()))
	})

	table.Show()
	return nil
}
