package table

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/db"
	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/pool/object"
)

var (
	table *Table[*Product]
	mysql *db.Mysql[*Product]
)

type Product struct {
	Id      int
	Name    string
	Date    string
	Time    string
	Sex     int
	Content string
}

func (p *Product) Columns() []string {
	return []string{
		"id", "name", "date", "time", "sex", "content",
	}
}

func (p *Product) Fields() []any {
	return []any{
		&p.Id, &p.Name, &p.Date, &p.Time, &p.Sex, &p.Content,
	}
}

func (p *Product) Values() []any {
	return []any{
		p.Id, p.Name, p.Date, p.Time, p.Sex, p.Content,
	}
}

func (p *Product) Clone(object.CtxInterface) itf.RowInterface {
	return &Product{}
}

func (p *Product) SetEmpty() {
}

func setup() {
	conf := config.Mysql{
		Host: "127.0.0.1", Port: 3306, Username: "root", Password: "root", Dbname: "test", Charset: "utf8mb4", ActiveMax: 10, ConnectionMax: 10,
	}
	err := db.Init(conf)
	if err != nil {
		fmt.Printf("init mysql error: %s", err)
	}

	mysql = db.NewMysql[*Product]()
	sql := []string{"CREATE TABLE `test`.`product` (",
		"`id` INT NOT NULL AUTO_INCREMENT,",
		"`name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '名称',",
		"`date` DATE NOT NULL DEFAULT '1970-01-01' COMMENT '日期',",
		"`time` TIMESTAMP(6) NOT NULL COMMENT '时间',",
		"`sex` INT NOT NULL DEFAULT 0 COMMENT '性别',",
		"`content` JSON NOT NULL COMMENT '内容',",
		"PRIMARY KEY (`id`))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci",
	}

	e := mysql.Exec(strings.Join(sql, ""))
	if e != nil {
		fmt.Printf("init table error: %s", e)
	}

	ssetup()
}

func teardown() {
	if err := mysql.Exec("drop table product"); err != nil {
		fmt.Printf("drop table faiure, error: %s", err)
	}
	steardown()
}

func TestTableInsert(t *testing.T) {
	table = NewTable[*Product]("product")
	data := make(map[string]any, 5)
	data["name"] = "kovey"
	data["date"] = "2021-01-01"
	data["time"] = "2021-01-02 11:11:11"
	data["sex"] = 1
	data["content"] = "{\"a\":3}"

	a, err := table.Insert(data)

	if err != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("id: %d", a)

	where := make(map[string]any)
	where["id"] = 1
	row := &Product{}
	e := table.FetchRow(where, row)
	if e != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("product: %v", row)
}

func TestTableUpdate(t *testing.T) {
	table = NewTable[*Product]("product")
	data := make(map[string]any)
	data["name"] = "test"
	where := make(map[string]any)
	where["id"] = 1

	a, err := table.Update(data, where)
	if err != nil {
		t.Errorf("update error: %s", err)
	}

	t.Logf("affected: %d", a)

	row := Product{}
	e := table.FetchRow(where, &row)
	if e != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("product: %v", row)
}

func TestTableDelete(t *testing.T) {
	table = NewTable[*Product]("product")
	where := make(map[string]any)
	where["id"] = 1

	a, err := table.Delete(where)
	if err != nil {
		t.Errorf("delete error: %s", err)
	}

	t.Logf("affected: %d", a)

	row := Product{}
	e := table.FetchRow(where, &row)
	if e != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("product: %v", row)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
