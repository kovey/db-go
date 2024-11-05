package main

import (
	"github.com/kovey/db-go/v3/examples/plugins/v1/migrations"
	"github.com/kovey/db-go/v3/migplug"
)

type Migrator struct {
}

func (m *Migrator) Register(c migplug.CoreInterface) {
	c.Add(&migrations.User{})
	c.Add(&migrations.UserInfo{})

}

func Migrate() migplug.PluginInterface {
	return &Migrator{}
}
