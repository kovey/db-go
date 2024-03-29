package ck

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sql"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/pool/object"
)

var (
	ckDb *ClickHouse[*Product]
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

func (p *Product) SetFetch() {
}

func setup() {
	conf := config.ClickHouse{
		Username: "default", Password: "", Dbname: "test", Debug: false, BlockSize: 1000000,
		Compress: 0, Timeout: config.Timeout{Read: 10, Exec: 60, Dial: 30},
		Cluster: config.Cluster{Open: "Off", Servers: make([]config.Addr, 0)},
		Server:  config.Addr{Host: "127.0.0.1", Port: 29001}, ActiveMax: 10, LifeTime: 100, ConnectionMax: 20,
	}
	err := Init(conf)
	if err != nil {
		fmt.Printf("init ckDb error: %s", err)
	}

	str := []string{
		"create TABLE product (",
		"id Int64 COMMENT '扩展ID',",
		"name String COMMENT 'APP ID', ",
		"date Date COMMENT '日期', ",
		"time DateTime COMMENT '时间', ",
		"sex UInt32 COMMENT '性别', ",
		"content String COMMENT '内容' ",
		") ENGINE=MergeTree ",
		"PARTITION BY (date) ",
		"ORDER BY (id) ",
		"SETTINGS index_granularity = 8192",
	}

	ckDb = NewClickHouse[*Product]()

	err = ckDb.Exec(strings.Join(str, ""))
	fmt.Printf("err: %s", err)
}

func teardown() {
	if err := ckDb.Exec("DROP TABLE product"); err != nil {
		fmt.Printf("drop table failure, error: %s", err)
	}
}

func TestInsert(t *testing.T) {
	in := sql.NewInsert("product")
	in.Set("id", 1).Set("name", "golang").Set("date", "2021-01-01").Set("time", "2021-01-01 11:11:11").Set("sex", 1).Set("content", "{\"name\":\"kovey\"}")

	_, err := ckDb.Insert(in)
	if err == nil {
		t.Fatalf("inser fail")
	}

	t.Logf("error tips: %s", err)
}

func TestBatchInsert(t *testing.T) {
	batch := sql.NewBatch("product")
	in := sql.NewInsert("product")
	in.Set("id", 1).Set("name", "rust").Set("date", "2021-02-01").Set("time", "2021-01-01 11:11:11").Set("sex", 1).Set("content", "{\"name\":\"kovey2\"}")
	batch.Add(in)
	in1 := sql.NewInsert("product")
	in1.Set("id", 2).Set("name", "php").Set("date", "2021-03-01").Set("time", "2021-01-01 11:11:11").Set("sex", 1).Set("content", "{\"name\":\"kovey1\"}")
	batch.Add(in1)

	a, err := ckDb.BatchInsert(batch)
	if err != nil {
		t.Errorf("batch insert fail, err: %s", err)
	}

	t.Logf("affected: %d", a)

	rows, e := ckDb.FetchAll("product", make(map[string]any), &Product{})
	if e != nil {
		t.Errorf("err: %s", err)
		return
	}

	for _, row := range rows {
		t.Logf("pro: %v", row)
	}
}

func TestQuery(t *testing.T) {

	sql := "select * from product"
	rows, err := ckDb.Query(sql, &Product{})
	if err != nil {
		t.Errorf("query[%s] fail, err: %s", sql, err)
	}

	for _, row := range rows {
		t.Logf("product: %v", row)
	}
}

func TestUpdate(t *testing.T) {
	where := sql.NewWhere()
	where.Eq("id", 1)
	up := sql.NewCkUpdate("product")
	up.Set("name", "java").Set("time", "2021-06-18 13:21:12").Where(where)
	a, err := ckDb.Update(up)
	if err != nil {
		t.Errorf("update fail, error: %s", err)
	}

	t.Logf("affected: %d", a)

	sql := "select * from product"
	rows, err := ckDb.Query(sql, &Product{})
	if err != nil {
		t.Errorf("query[%s] fail, err: %s", sql, err)
	}

	for _, row := range rows {
		t.Logf("product: %v", row)
	}
}

func TestDelete(t *testing.T) {
	where := sql.NewWhere()
	where.Eq("id", 1)
	del := sql.NewCkDelete("product")
	del.Where(where)
	a, err := ckDb.Delete(del)
	if err != nil {
		t.Errorf("delete fail, error: %s", err)
	}

	t.Logf("affected: %d", a)

	sql := "select * from product"
	rows, err := ckDb.Query(sql, &Product{})
	if err != nil {
		t.Errorf("query[%s] fail, err: %s", sql, err)
	}

	for _, row := range rows {
		t.Logf("product: %v", row)
	}
}

func TestFatchAll(t *testing.T) {
	rows, err := ckDb.FetchAll("product", make(map[string]any), &Product{})
	if err != nil {
		t.Errorf("fetch all error: %s", err)
	}

	for _, row := range rows {
		t.Logf("product: %v", row)
	}
}

func TestFatchRow(t *testing.T) {
	row := Product{}
	err := ckDb.FetchRow("product", make(map[string]any), &row)
	if err != nil {
		t.Errorf("fetch all error: %s", err)
	}

	t.Logf("product: %v", row)
}

func TestMain(m *testing.M) {
	debug.SetLevel(debug.Debug_Info)
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
