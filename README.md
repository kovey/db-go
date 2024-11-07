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
		DataSourceName: "root:password@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true",
		MaxIdleTime:    time.Second * 60,
		MaxLifeTime:    time.Second * 120,
		MaxIdleConns:   10,
		MaxOpenConns:   50,
	}

	if err := db.Init(conf); err != nil {
		panic(err)
	}

    db.Table(context.Background(), "user", func (table ksql.TableInterface) {
		table.Create()
		table.AddString("account", 63).Comment("账号").Default("")
		table.AddDate("create_date").Comment("创建日期").Default("")
		table.AddBigInt("create_time").Comment("创建时间").Unsigned().Default("0")
		table.AddInt("id").AutoIncrement().Comment("主键").Unsigned()
		table.AddString("nickname", 63).Comment("昵称").Default("")
		table.AddString("password", 64).Comment("密码").Default("")
		table.AddTinyInt("status").Comment("状态 0 - 正常 1 - 删除").Default("0")
		table.AddBigInt("update_time").Comment("更新时间").Unsigned().Default("0")
		table.AddPrimary("id").Engine("InnoDB").Charset("utf8mb4").Collate("utf8mb4_0900_ai_ci")
		table.AddUnique("idx_account", "account")
    })

	ssetup()
}

func TestCreateTable(t *testing.T) {
    err := db.Table(context.Background(), "user", func (table ksql.TableInterface) {
		table.Create()
		table.AddString("account", 63).Comment("账号").Default("")
		table.AddDate("create_date").Comment("创建日期").Default("")
		table.AddBigInt("create_time").Comment("创建时间").Unsigned().Default("0")
		table.AddInt("id").AutoIncrement().Comment("主键").Unsigned()
		table.AddString("nickname", 63).Comment("昵称").Default("")
		table.AddString("password", 64).Comment("密码").Default("")
		table.AddTinyInt("status").Comment("状态 0 - 正常 1 - 删除").Default("0")
		table.AddBigInt("update_time").Comment("更新时间").Unsigned().Default("0")
		table.AddPrimary("id").Engine("InnoDB").Charset("utf8mb4").Collate("utf8mb4_0900_ai_ci")
		table.AddUnique("idx_account", "account")
    })

    if err != nil {
        t.Fatal(err)
    }
}

func TestAlterTable(t *testing.T) {
    err := db.Table(context.Background(), "user", func (table ksql.TableInterface) {
	    table.Alter()
	    table.AddString("foo", 63).Comment("foo").Default("")
            table.DropColumn("boo").DropIndex("idx_xxxx")
	    table.ChangeColumn("nickname", "nick", "varchar", 31, 0).Comment("nickname").Default("")
	    table.Engine("InnoDB").Charset("utf8mb4").Collate("utf8mb4_0900_ai_ci")
	    table.AddUnique("idx_nick", "nick")
    })

    if err != nil {
        t.Fatal(err)
    }
}

func TestInsert(t *testing.T) {
    u := NewUser()
    u.Account = "kovey"
    u.Nickname = "kovey_nickname"
    u.Password = "1232555"
    u.Status = 0
    u.CreateTime = time.Now().Unix()
    u.CreateDate = time.Now()
    u.UpdateTime = u.CreateTime

    if err := u.Save(context.Background()); err != nil {
	    t.Fatal(err)
    }
}

func TestFetchRow(t *testing.T) {
	ctx := context.Background()
	u := NewUser()
	if err := model.Query(u).Where("id", "=", 1).First(ctx, u); err != nil {
        t.Fatal(err)
	}

	t.Log("user: ", u)
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	u := NewUser()
	if err := model.Query(u).Where("id", "=", 1).First(ctx, u); err != nil {
        t.Fatal(err)
	}

	t.Log("user: ", u)
        u.Account = "kovey"
        u.Nickname = "kovey_nickname"
        u.Password = "1232555"
        u.Status = 1
        u.UpdateTime = time.Now().Unix()

        if err := u.Save(context.Background()); err != nil {
	        t.Fatal(err)
        }
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	u := NewUser()
	if err := model.Query(u).Where("id", "=", 1).First(ctx, u); err != nil {
        t.Fatal(err)
	}

	t.Log("user: ", u)
        if err := u.Delete(context.Background()); err != nil {
	        t.Fatal(err)
        }
}

func TestFetchAll(t *testing.T) {
	ctx := context.Background()
        var users []*User
	if err := model.Query(NewUser()).Where("id", "<", 100).Limit(10).All(ctx, &users); err != nil {
        t.Fatal(err)
	}

        for _, u := range users {
            t.Log(u)
        }
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
