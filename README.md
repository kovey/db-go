## kovey sql by golang
#### Description
###### This is a database library, no reflect
###### Usage
    go get -u github.com/kovey/db-go/v3
### Examples
```golang
import (
	"context"

	"github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/model"
	"github.com/kovey/db-go/v3/db"
	"time"
)

type User struct {
	*model.Model  `db:"-" json:"-"` // model
	Account       string            `db:"account" json:"account"`                 // 账号
	CreateDate    time.Time         `db:"create_date" json:"create_date"`         // 创建日期
	CreateTime    int64             `db:"create_time" json:"create_time"`         // 创建时间
	Id            int64             `db:"id" json:"id"`                           // 主键
	Nickname      string            `db:"nickname" json:"nickname"`               // 昵称
	Password      string            `db:"password" json:"password"`               // 密码
	Status        int               `db:"status" json:"status"`                   // 状态 0 - 正常 1 - 封禁
	UpdateTime    int64             `db:"update_time" json:"update_time"`         // 更新时间
}

func NewUser() *User {
	return &User{Model: model.NewModel("user", "id", model.Type_Int)}
}

func (self *User) Save(ctx context.Context) error {
	return self.Model.Save(ctx, self)
}

func (self *User) Clone() ksql.RowInterface {
	return NewUser()
}

func (self *User) Values() []any {
	return []any{&self.Account, &self.CreateDate, &self.CreateTime, &self.Email, &self.Id, &self.Nickname, &self.Password, &self.Status, &self.UpdateTime}
}

func (self *User) Columns() []string {
	return []string{"account", "create_date", "create_time", "id", "nickname", "password", "status", "update_time"}
}

func (self *User) Delete(ctx context.Context) error {
	return self.Model.Delete(ctx, self)
}

func setup() {
	conf := db.Config{
		DriverName:     "mysql",
		DataSourceName: "root:some34QA@123@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true",
		MaxIdleTime:    time.Second * 60,
		MaxLifeTime:    time.Second * 120,
		MaxIdleConns:   10,
		MaxOpenConns:   50,
	}

	if err := db.Init(conf); err != nil {
		panic(err)
	}

    // TODO
    db.Table(context.Background(), "user", func (table ksql.TableInterface) {
    })
	ssetup()
}

func teardown() {
	if err := db.DropTable(context.Background(), "user"); err != nil {
		fmt.Printf("drop table failure, error: %s", err)
	}

	steardown()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
```
