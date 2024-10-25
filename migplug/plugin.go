package migplug

import "context"

type MigrateInterface interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
	Id() uint64
	Name() string
	Version() string
}

type CoreInterface interface {
	Add(migrate MigrateInterface)
}

type PluginInterface interface {
	Register(core CoreInterface)
}
