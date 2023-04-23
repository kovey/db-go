package tests

import(
	"github.com/kovey/db-go/v2/table"
	"github.com/kovey/db-go/v2/model"
	"github.com/shopspring/decimal"
)

type UserTable struct {
	*table.Table[*UserRow]
}

func NewUserTable() *UserTable {
	return &UserTable{Table: table.NewTable[*UserRow]("user")}
}

type UserRow struct {
	*model.Base[*UserRow]
	Id int64 `db:"id"`
	Nickname string `db:"nickname"`
	Platform_id string `db:"platform_id"`
	Balance decimal.Decimal `db:"balance"`
	Label string `db:"label"`
	Ip_addr string `db:"ip_addr"`
	Currency string `db:"currency"`
	Status int32 `db:"status"`
	Create_time int64 `db:"create_time"`
	Update_time int64 `db:"update_time"`
}

func NewUserRow() *UserRow {
	return &UserRow{Base: model.NewBase[*UserRow](NewUserTable(), model.NewPrimaryId("id", model.Int))}
}
	