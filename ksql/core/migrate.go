package core

import (
	"context"
	"sort"
	"time"

	ksql "github.com/kovey/db-go/v3"
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
	beginTime := time.Now().UnixMilli()
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
	debug.Info("migrate end. time used: %.3fms", float64(time.Now().UnixMilli()-beginTime)*0.001)
}

func Has(ctx context.Context, id uint64) (bool, error) {
	if err := check(ctx); err != nil {
		return false, err
	}

	return model.Query(newMigrateTable()).Where("migrate_id", "=", id).Exist(ctx)
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
		table.Create()
		table.AddInt("id").AutoIncrement().Comment("主键").Unsigned()
		table.AddBigInt("migrate_id").Comment("迁移ID").Unsigned().Default("0")
		table.AddString("name", 255).Comment("迁移文件").Default("")
		table.AddTinyInt("status").Comment("状态 0 - 已迁移 1 - 未迁移").Default("1")
		table.AddString("version", 15).Comment("版本").Default("v1")
		table.AddTimestamp("create_time").Default(ksql.CURRENT_TIMESTAMP).Comment("创建时间")
		table.AddTimestamp("update_time").Default(ksql.CURRENT_TIMESTAMP_ON_UPDATE_CURRENT_TIMESTAMP).Comment("创建时间")
		table.AddPrimary("id").Engine("InnoDB").Charset("utf8mb4").Collate("utf8mb4_0900_ai_ci")
		table.AddUnique("idx_migrate_id", "migrate_id")
	})
}
