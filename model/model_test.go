package model

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/db"
	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/db-go/v2/table"
)

var (
	mysql *db.Mysql[*Product]
)

type ProTable struct {
	table.Table[*Product]
}

type Product struct {
	*Base[*Product]
	Id      int
	Name    string
	Date    string
	Time    string
	Sex     int
	Content string
}

func (p *Product) Columns() []*meta.Column {
	return []*meta.Column{
		meta.NewColumn("id"), meta.NewColumn("name"), meta.NewColumn("date"), meta.NewColumn("time"), meta.NewColumn("sex"), meta.NewColumn("content"),
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

func (p *Product) Clone() itf.RowInterface {
	return &Product{}
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
	if err := mysql.Exec("drop table product"); err != nil {
		fmt.Printf("drop table failure, error: %s", err)
	}

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

	if err := pro1.FetchRow(where, pro1); err != nil {
		t.Errorf("FetchRow failure, error: %s", err)
	}
	pro1.Name = "chelsea"
	if err := pro1.Save(pro1); err != nil {
		t.Fatalf("save failure, error: %s", err)
	}
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
	if err := pr2.FetchRow(where, pr2); err != nil {
		t.Fatalf("FetchRow failure, error: %s", err)
	}
	t.Logf("pr2: %t", pr2.Empty())
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
