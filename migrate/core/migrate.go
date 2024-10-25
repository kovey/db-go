package core

import (
	"context"
	"sort"
	"time"

	"github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/db-go/v3/migplug"
	"github.com/kovey/db-go/v3/model"
	"github.com/kovey/debug-go/debug"
)

type MigrateType byte

const (
	Type_Up   MigrateType = 1
	Type_Down MigrateType = 2
)

func Migrate(ctx context.Context, t MigrateType) {
	if err := check(ctx); err != nil {
		debug.Erro("create migrate table error: %s", err)
		return
	}

	sort.Sort(mig)
	debug.Info("migrate begin...")
	mig.Range(func(key uint64, mi migplug.MigrateInterface) {
		m := newMigrateTable()
		err := model.Query(m).Where("migrate_id", "=", mi.Id()).First(ctx, m)
		if err != nil {
			debug.Erro("fecth migrate error: %s", err)
			return
		}

		switch t {
		case Type_Up:
			if !m.Empty() {
				return
			}

			debug.Info("migrate upgrade[%s] begin...", mi.Name())
			defer debug.Info("migrate upgrade[%s] end.", mi.Name())
			if err := mi.Up(ctx); err != nil {
				debug.Erro("migrate upgrade error: %s", err)
				return
			}
			Create(ctx, mi)
			debug.Info("migrate upgrade[%s] success.", mi.Name())
		case Type_Down:
			if m.Empty() {
				return
			}

			debug.Info("migrate downgrade[%s] begin...", mi.Name())
			defer debug.Info("migrate downgrade[%s] end.", mi.Name())
			if err := mi.Down(ctx); err != nil {
				debug.Erro("migrate downgrade error: %s", err)
				return
			}

			if err := m.Delete(ctx); err != nil {
				debug.Erro("migrate downgrade error: %s", err)
			}
			debug.Info("migrate downgrade[%s] success.", mi.Name())
		}
	})
	debug.Info("migrate end.")
}

func Has(ctx context.Context, id uint64) bool {
	check(ctx)
	ok, _ := model.Query(newMigrateTable()).Where("migrate_id", "=", id).Exist(ctx)
	return ok
}

func Create(ctx context.Context, mi migplug.MigrateInterface) {
	m := newMigrateTable()
	err := model.Query(m).Where("migrate_id", "=", mi.Id()).First(ctx, m)
	if err != nil {
		debug.Erro("create migrator error: %s", err)
		return
	}

	if !m.Empty() {
		debug.Erro("migrator[%s] is created", mi.Name())
		return
	}

	m.MigrateId = mi.Id()
	m.Name = mi.Name()
	m.Status = 0
	m.Version = mi.Version()
	m.CreateTime = time.Now().Format(time.DateTime)
	m.UpdateTime = m.CreateTime
	if err := m.Save(ctx); err != nil {
		debug.Erro("create migrator[%s]  error: %s", mi.Name(), err)
	}
}

func check(ctx context.Context) error {
	ok, err := db.HasTable(ctx, migrate_table_name)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	return db.Table(ctx, migrate_table_name, func(table ksql.TableInterface) {
		table.Create().AddColumn("id", "int", 11, 0).AutoIncrement().Comment("主键").Unsigned()
		table.AddColumn("migrate_id", "bigint", 20, 0).Comment("迁移ID").Unsigned().Default("0", false)
		table.AddColumn("name", "varchar", 255, 0).Comment("迁移文件").Default("", false)
		table.AddColumn("status", "tinyint", 3, 0).Comment("状态 0 - 已迁移 1 - 未迁移").Default("1", false)
		table.AddColumn("version", "varchar", 15, 0).Comment("版本").Default("v1", false)
		table.AddColumn("create_time", "timestamp", 0, 0).Default(ksql.CURRENT_TIMESTAMP, true).Comment("创建时间")
		table.AddColumn("update_time", "timestamp", 0, 0).Default(ksql.CURRENT_TIMESTAMP_ON_UPDATE_CURRENT_TIMESTAMP, true).Comment("创建时间")
		table.AddPrimary("id").Engine("InnoDB").Charset("utf8mb4").Collate("utf8mb4_0900_ai_ci")
		table.AddIndex("idx_migrate_id", ksql.Index_Type_Normal, "migrate_id")
	})
}
