package model

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kovey/db-go/db"
	"github.com/kovey/db-go/table"
)

var (
	mysql *db.Mysql
)

type ProTable struct {
	table.Table
}

type Product struct {
	Base
	Id      int
	Name    string
	Date    string
	Time    string
	Sex     int
	Content string
}

func NewProTable() *ProTable {
	return &ProTable{*table.NewTable("product")}
}

func NewProduct() Product {
	pro := Product{Base{table: NewProTable(), primaryId: "Id"}, 0, "", "", "", 0, "{}"}

	return pro
}

func setup() {
	err := db.Init("127.0.0.1", 3306, "root", "123456", "test", "utf8mb4", 10, 10)
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
}

func teardown() {
	mysql.Exec("drop table product")
}

func TestModelSave(t *testing.T) {
	pro := NewProduct()
	pro.Name = "kovey"
	pro.Date = "2021-08-12"
	pro.Time = "2021-08-12 13:12:12"
	pro.Sex = 1
	pro.Content = "{\"where\":123}"

	err := pro.Save(&pro)
	if err != nil {
		t.Errorf("product save fail, error: %s", err)
	}

	t.Logf("id: %d", pro.Id)

}

func TestModelFetchRow(t *testing.T) {
	where := make(map[string]interface{})
	where["id"] = 1
	pr1 := NewProduct()
	err := pr1.FetchRow(where, &pr1)
	if err != nil {
		t.Errorf("fetch row err: %s", err)
	}

	t.Logf("pr1: %v", pr1)
}

func TestModelDelete(t *testing.T) {
	where := make(map[string]interface{})
	where["id"] = 1
	pr1 := NewProduct()
	err := pr1.FetchRow(where, &pr1)
	if err != nil {
		t.Errorf("fetch row err: %s", err)
	}

	err = pr1.Delete(pr1)
	if err != nil {
		t.Errorf("delete row err: %s", err)
	}

	pr2 := NewProduct()
	pr2.FetchRow(where, &pr2)
	t.Logf("pr2: %v", pr2)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
