package model

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sharding"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/db-go/v2/table"
)

var (
	shardDb *sharding.Mysql[*ProductSharding]
)

type ProTableSharding struct {
	*table.TableSharding[*ProductSharding]
}

type ProductSharding struct {
	*BaseSharding[*ProductSharding]
	Id      int
	Name    string
	Date    string
	Time    string
	Sex     int
	Content string
}

func (p *ProductSharding) Columns() []*meta.Column {
	return []*meta.Column{
		meta.NewColumn("id"), meta.NewColumn("name"), meta.NewColumn("date"), meta.NewColumn("time"), meta.NewColumn("sex"), meta.NewColumn("content"),
	}
}

func (p *ProductSharding) Fields() []any {
	return []any{
		&p.Id, &p.Name, &p.Date, &p.Time, &p.Sex, &p.Content,
	}
}

func (p *ProductSharding) Values() []any {
	return []any{
		p.Id, p.Name, p.Date, p.Time, p.Sex, p.Content,
	}
}

func (p *ProductSharding) Clone() itf.RowInterface {
	return &Product{}
}

func NewProTableSharding() *ProTableSharding {
	return &ProTableSharding{table.NewTableSharding[*ProductSharding]("product", true)}
}

func NewProductSharding() *ProductSharding {
	pro := &ProductSharding{NewBaseSharding[*ProductSharding](NewProTableSharding(), NewPrimaryId("id", Int)), 0, "", "", "", 0, "{}"}

	return pro
}

func ssetup() {
	mas := make([]config.Mysql, 2)

	mas[0] = config.Mysql{
		Host: "127.0.0.1", Port: 3306, Username: "root", Password: "root", Dbname: "test", Charset: "utf8mb4", ActiveMax: 10, ConnectionMax: 10,
	}
	mas[1] = config.Mysql{
		Host: "127.0.0.1", Port: 3306, Username: "root", Password: "root", Dbname: "test", Charset: "utf8mb4", ActiveMax: 10, ConnectionMax: 10,
	}

	sharding.Init(mas, mas)

	shardDb = sharding.NewMysql[*ProductSharding](true)
	sql := []string{"CREATE TABLE `{table}` (",
		"`id` INT NOT NULL AUTO_INCREMENT,",
		"`name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '名称',",
		"`date` DATE NOT NULL DEFAULT '1970-01-01' COMMENT '日期',",
		"`time` TIMESTAMP(6) NOT NULL COMMENT '时间',",
		"`sex` INT NOT NULL DEFAULT 0 COMMENT '性别',",
		"`content` JSON NOT NULL COMMENT '内容',",
		"PRIMARY KEY (`id`))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci",
	}

	shardDb.AddSharding(0).AddSharding(1)

	e := shardDb.Exec(0, strings.Replace(strings.Join(sql, ""), "{table}", "product_0", 1))
	if e != nil {
		fmt.Printf("init table error: %s", e)
	}

	e = shardDb.Exec(1, strings.Replace(strings.Join(sql, ""), "{table}", "product_1", 1))
	if e != nil {
		fmt.Printf("init table error: %s", e)
	}
}

func steardown() {
	if err := shardDb.Exec(0, "drop table product_0"); err != nil {
		fmt.Printf("drop table failure, error: %s", err)
	}
	if err := shardDb.Exec(1, "drop table product_1"); err != nil {
		fmt.Printf("drop table failure, error: %s", err)
	}
}

func TestModelShardingSave(t *testing.T) {
	pro := NewProductSharding()
	pro.Name = "kovey"
	pro.Date = "2021-08-12"
	pro.Time = "2021-08-12 13:12:12"
	pro.Sex = 1
	pro.Content = "{\"where\":123}"

	err := pro.Save(0, pro)
	if err != nil {
		t.Errorf("product save fail, error: %s", err)
	}

	t.Logf("id: %d", pro.Id)

	pro1 := NewProductSharding()
	where := make(map[string]any)
	where["id"] = pro.Id

	if err := pro1.FetchRow(0, where, pro1); err != nil {
		t.Fatalf("FetchRow failure, error: %s", err)
	}
	pro1.Name = "chelsea"
	if err := pro1.Save(0, pro1); err != nil {
		t.Fatalf("save failure, error: %s", err)
	}
}

func TestModelShardingFetchRow(t *testing.T) {
	where := make(map[string]any)
	where["id"] = 1
	pr1 := NewProductSharding()
	err := pr1.FetchRow(0, where, pr1)
	if err != nil {
		t.Errorf("fetch row err: %s", err)
	}

	t.Logf("pr1: %v", pr1)
}

func TestModelShardingDelete(t *testing.T) {
	where := make(map[string]any)
	where["id"] = 1
	pr1 := NewProductSharding()
	err := pr1.FetchRow(0, where, pr1)
	if err != nil {
		t.Errorf("fetch row err: %s", err)
	}

	err = pr1.Delete(0, pr1)
	if err != nil {
		t.Errorf("delete row err: %s", err)
	}

	pr2 := NewProductSharding()
	if err := pr2.FetchRow(0, where, pr2); err != nil {
		t.Fatalf("FetchRow failure, error: %s", err)
	}
	t.Logf("pr2: %t", pr2.Empty())
}
