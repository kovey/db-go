package model

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kovey/config-go/config"
	"github.com/kovey/db-go/db"
	"github.com/kovey/db-go/table"
)

var (
	mysql *db.Mysql[*Product]
)

type ProTable struct {
	table.Table[*Product]
}

type Product struct {
	*Base[*Product]
	Id      int    `db:"id"`
	Name    string `db:"name"`
	Date    string `db:"date"`
	Time    string `db:"time"`
	Sex     int    `db:"sex"`
	Content string `db:"content"`
}

func NewProTable() *ProTable {
	return &ProTable{*table.NewTable[*Product]("product")}
}

func NewProduct() *Product {
	return &Product{NewBase[*Product](NewProTable(), NewPrimaryId("id", Int)), 0, "", "", "", 0, "{}"}
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
	mysql.Exec("drop table product")
	steardown()
}

func TestModelSave(t *testing.T) {
	pro := NewProduct()
	pro.Name = "kovey"
	pro.Date = "2021-08-12"
	pro.Time = "2021-08-12 13:12:12"
	pro.Sex = 1
	pro.Content = "{\"where\":123}"

	err := pro.Save(pro)
	if err != nil {
		t.Errorf("product save fail, error: %s", err)
	}

	t.Logf("id: %d", pro.Id)

	pro1 := NewProduct()
	where := make(map[string]any)
	where["id"] = pro.Id

	pro1.FetchRow(where, pro1)
	pro1.Name = "chelsea"
	pro1.Save(pro1)
}

func TestModelFetchRow(t *testing.T) {
	where := make(map[string]any)
	where["id"] = 1
	pr1 := NewProduct()
	err := pr1.FetchRow(where, pr1)
	if err != nil {
		t.Errorf("fetch row err: %s", err)
	}

	t.Logf("pr1: %v", pr1)
}

func TestModelDelete(t *testing.T) {
	where := make(map[string]any)
	where["id"] = 1
	pr1 := NewProduct()
	err := pr1.FetchRow(where, pr1)
	if err != nil {
		t.Errorf("fetch row err: %s", err)
	}

	err = pr1.Delete(pr1)
	if err != nil {
		t.Errorf("delete row err: %s", err)
	}

	pr2 := NewProduct()
	pr2.FetchRow(where, pr2)
	t.Logf("pr2: %v", pr2)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
