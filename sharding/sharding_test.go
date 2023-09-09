package sharding

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/itf"
	"github.com/kovey/db-go/v2/sql"
)

var (
	mysql *Mysql[*Product]
)

type Product struct {
	Id      int    `db:"id"`
	Name    string `db:"name"`
	Date    string `db:"date"`
	Time    string `db:"time"`
	Sex     int    `db:"sex"`
	Content string `db:"content"`
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

func (p *Product) Clone() itf.RowInterface {
	return &Product{}
}

func (p *Product) SetEmpty() {
}

func setup() {
	mas := make([]config.Mysql, 2)

	mas[0] = config.Mysql{
		Host: "127.0.0.1", Port: 3306, Username: "root", Password: "root", Dbname: "test", Charset: "utf8mb4", ActiveMax: 10, ConnectionMax: 10,
	}
	mas[1] = config.Mysql{
		Host: "127.0.0.1", Port: 3306, Username: "root", Password: "root", Dbname: "test", Charset: "utf8mb4", ActiveMax: 10, ConnectionMax: 10,
	}

	Init(mas, mas)

	mysql = NewMysql[*Product](true)
	sql := []string{"CREATE TABLE `{table}` (",
		"`id` INT NOT NULL AUTO_INCREMENT,",
		"`name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '名称',",
		"`date` DATE NOT NULL DEFAULT '1970-01-01' COMMENT '日期',",
		"`time` TIMESTAMP(6) NOT NULL COMMENT '时间',",
		"`sex` INT NOT NULL DEFAULT 0 COMMENT '性别',",
		"`content` JSON NOT NULL COMMENT '内容',",
		"PRIMARY KEY (`id`))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci",
	}

	mysql.AddSharding(0).AddSharding(1)

	e := mysql.Exec(0, strings.Replace(strings.Join(sql, ""), "{table}", "product_0", 1))
	if e != nil {
		fmt.Printf("init table error: %s", e)
	}

	e = mysql.Exec(1, strings.Replace(strings.Join(sql, ""), "{table}", "product_1", 1))
	if e != nil {
		fmt.Printf("init table error: %s", e)
	}
}

func teardown() {
	if err := mysql.Exec(0, "drop table product_0"); err != nil {
		fmt.Printf("drop table failure, error: %s", err)
	}
	if err := mysql.Exec(1, "drop table product_1"); err != nil {
		fmt.Printf("drop table failure, error: %s", err)
	}
}

func TestInsert(t *testing.T) {
	err := mysql.Transaction(func(tx *Tx) error {
		in := sql.NewInsert("product_0")
		in.Set("name", "golang").Set("date", "2021-01-01").Set("time", "2021-01-01 11:11:11").Set("sex", 1).Set("content", "{\"name\":\"kovey\"}")

		id, err := mysql.Insert(0, in)
		if err != nil {
			return err
		}

		t.Logf("insert id[%d]", id)

		in1 := sql.NewInsert("product_1")
		in1.Set("name", "php").Set("date", "1995-01-01").Set("time", "1995-01-01 11:11:11").Set("sex", 1).Set("content", "{\"name\":\"rust\"}")

		id, err = mysql.Insert(1, in1)
		if err != nil {
			return err
		}

		t.Logf("insert id[%d]", id)
		return nil
	})

	if err != nil {
		t.Fatalf("error: %s", err)
	}
}

func TestBatchInsert(t *testing.T) {
	batch := sql.NewBatch("product_0")
	in := sql.NewInsert("product_0")
	in.Set("name", "rust").Set("date", "2021-02-01").Set("time", "2021-01-01 11:11:11").Set("sex", 1).Set("content", "{\"name\":\"kovey2\"}")
	batch.Add(in)
	in1 := sql.NewInsert("product_0")
	in1.Set("name", "java").Set("date", "2021-03-01").Set("time", "2021-01-01 11:11:11").Set("sex", 1).Set("content", "{\"name\":\"kovey1\"}")
	batch.Add(in1)

	a, err := mysql.BatchInsert(0, batch)
	if err != nil {
		t.Errorf("batch insert fail, err: %s", err)
	}

	t.Logf("affected: %d", a)

	rows, e := mysql.FetchAll(0, "product_0", make(map[string]any), &Product{})
	if e != nil {
		t.Errorf("err: %s", err)
		return
	}

	for _, row := range rows {
		t.Logf("pro: %v", row)
	}
}

func TestQuery(t *testing.T) {

	sql := "select * from product_1"
	rows, err := mysql.Query(1, sql, &Product{})
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
	up := sql.NewUpdate("product_0")
	up.Set("name", "java").Set("time", "2021-06-18 13:21:12").Where(where)
	a, err := mysql.Update(0, up)
	if err != nil {
		t.Errorf("update fail, error: %s", err)
	}

	t.Logf("affected: %d", a)

	sql := "select * from product_0"
	rows, err := mysql.Query(0, sql, &Product{})
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
	del := sql.NewDelete("product_0")
	del.Where(where)
	a, err := mysql.Delete(0, del)
	if err != nil {
		t.Errorf("delete fail, error: %s", err)
	}

	t.Logf("affected: %d", a)

	sql := "select * from product_0"
	rows, err := mysql.Query(0, sql, &Product{})
	if err != nil {
		t.Errorf("query[%s] fail, err: %s", sql, err)
	}

	for _, row := range rows {
		t.Logf("product: %v", row)
	}
}

func TestFatchAll(t *testing.T) {
	rows, err := mysql.FetchAll(1, "product_1", make(map[string]any), &Product{})
	if err != nil {
		t.Errorf("fetch all error: %s", err)
	}

	for _, row := range rows {
		t.Logf("product: %v", row)
	}
}

func TestFatchRow(t *testing.T) {
	row := Product{}
	err := mysql.FetchRow(1, "product_1", make(map[string]any), &row)
	if err != nil {
		t.Errorf("fetch all error: %s", err)
	}

	t.Logf("product: %v", row)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
