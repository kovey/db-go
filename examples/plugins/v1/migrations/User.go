package migrations

import (
	"context"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

type User struct {
}

func (self *User) Up(ctx context.Context) error {
	// TODO Code
	if ok, err := db.HasColumn(ctx, "user", "foo"); err == nil && !ok {
		return db.Table(ctx, "user", func(table ksql.TableInterface) {
			table.AddColumn("foo", "varchar", 10, 0).Default("").Comment("test").Nullable()
		})
	}

	return nil
}

func (self *User) Down(ctx context.Context) error {
	// TODO Code
	if ok, err := db.HasColumn(ctx, "user", "foo"); err == nil && ok {
		return db.Table(ctx, "user", func(table ksql.TableInterface) {
			table.DropColumn("foo")
		})
	}

	return nil
}

func (self *User) Id() uint64 {
	return 1729871031576555000
}

func (self *User) Name() string {
	return "User"
}

func (self *User) Version() string {
	return "v1"
}
