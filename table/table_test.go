package table

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kovey/config-go/config"
	"github.com/kovey/db-go/db"
)

var (
	table *Table
	mysql *db.Mysql
)

type Product struct {
	Id      int    `db:"id"`
	Name    string `db:"name"`
	Date    string `db:"date"`
	Time    string `db:"time"`
	Sex     int    `db:"sex"`
	Content string `db:"content"`
}

func setup() {
	conf := config.Mysql{
		Host: "127.0.0.1", Port: 3306, Username: "root", Password: "123456", Dbname: "test", Charset: "utf8mb4", ActiveMax: 10, ConnectionMax: 10,
	}
	err := db.Init(conf)
	if err != nil {
		fmt.Printf("init mysql error: %s", err)
	}

	mysql = db.NewMysql()
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
	mysql.Exec("drop table product")
	steardown()
}

func TestTableInsert(t *testing.T) {
	table = NewTable("product")
	data := make(map[string]interface{}, 5)
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

	where := make(map[string]interface{})
	where["id"] = 1
	row, e := table.FetchRow(where, Product{})
	if e != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("product: %v", row.(Product))
}

func TestTableUpdate(t *testing.T) {
	table = NewTable("product")
	data := make(map[string]interface{})
	data["name"] = "test"
	where := make(map[string]interface{})
	where["id"] = 1

	a, err := table.Update(data, where)
	if err != nil {
		t.Errorf("update error: %s", err)
	}

	t.Logf("affected: %d", a)

	row, e := table.FetchRow(where, Product{})
	if e != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("product: %v", row.(Product))
}

func TestTableDelete(t *testing.T) {
	table = NewTable("product")
	where := make(map[string]interface{})
	where["id"] = 1

	a, err := table.Delete(where)
	if err != nil {
		t.Errorf("delete error: %s", err)
	}

	t.Logf("affected: %d", a)

	row, e := table.FetchRow(where, Product{})
	if e != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("product: %v", row.(Product))
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
