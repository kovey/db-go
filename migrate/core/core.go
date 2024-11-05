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
	"github.com/kovey/debug-go/color"
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

func (m *migrates) Keys() []uint64 {
	return m.keys
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

	ids := db.ToList(mig.Keys())
	table := gui.NewTable()
	table.Add(0, "Migrator")
	table.Add(0, "Version")
	table.Add(0, "Migrated")
	table.Add(0, "Migrate Time")
	if len(ids) == 0 {
		table.Add(1, "No migrations")
		table.Add(1, "")
		table.Add(1, "")
		table.Add(1, "")
		table.Show()
		return nil
	}

	var migrations []*migrateTable
	if err := model.Query(newMigrateTable()).WhereIn("migrate_id", ids).All(ctx, &migrations); err != nil {
		return err
	}

	i := 0
	mig.Range(func(key uint64, mi migplug.MigrateInterface) {
		i++
		table.Add(i, mi.Name())
		table.Add(i, mi.Version())
		yes := false
		var upTime = ""
		for _, migration := range migrations {
			if migration.MigrateId == key {
				yes = true
				upTime = migration.CreateTime
				break
			}
		}

		if yes {
			table.AddColor(i, "Yes", color.Color_Green)
		} else {
			table.AddColor(i, "No", color.Color_Red)
		}
		table.Add(i, upTime)
	})

	table.Show()
	return nil
}
