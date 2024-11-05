package migrations

import (
	"context"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

type UserInfo struct {
}

func (self *UserInfo) Up(ctx context.Context) error {
	// TODO Code
	return nil
}

func (self *UserInfo) Down(ctx context.Context) error {
	// TODO Code
	if ok, err := db.HasColumn(ctx, "user", "foo"); err == nil && ok {
		return db.Table(ctx, "user", func(table ksql.TableInterface) {
			table.DropColumn("foo")
		})
	}
	return nil
}

func (self *UserInfo) Id() uint64 {
	return 1730812888306397000
}

func (self *UserInfo) Name() string {
	return "UserInfo"
}

func (self *UserInfo) Version() string {
	return "v1"
}
