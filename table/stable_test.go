package table

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/sharding"
)

var (
	shardTable *TableSharding[*Product]
	shardDb    *sharding.Mysql[*Product]
)

func ssetup() {
	mas := make([]config.Mysql, 2)

	mas[0] = config.Mysql{
		Host: "127.0.0.1", Port: 3306, Username: "root", Password: "root", Dbname: "test", Charset: "utf8mb4", ActiveMax: 10, ConnectionMax: 10,
	}
	mas[1] = config.Mysql{
		Host: "127.0.0.1", Port: 3306, Username: "root", Password: "root", Dbname: "test", Charset: "utf8mb4", ActiveMax: 10, ConnectionMax: 10,
	}

	sharding.Init(mas, mas)

	shardDb = sharding.NewMysql[*Product](true)
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
	shardDb.Exec(0, "drop table product_0")
	shardDb.Exec(1, "drop table product_1")
}

func TestTableShardingInsert(t *testing.T) {
	shardTable = NewTableSharding[*Product]("product", true)
	data := make(map[string]any, 5)
	data["name"] = "kovey"
	data["date"] = "2021-01-01"
	data["time"] = "2021-01-02 11:11:11"
	data["sex"] = 1
	data["content"] = "{\"a\":3}"

	a, err := shardTable.Insert(0, data)

	if err != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("id: %d", a)

	where := make(map[string]any)
	where["id"] = 1
	row, e := shardTable.FetchRow(0, where, &Product{})
	if e != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("product: %v", row)
}

func TestTableShardingUpdate(t *testing.T) {
	shardTable = NewTableSharding[*Product]("product", true)
	data := make(map[string]any)
	data["name"] = "test"
	where := make(map[string]any)
	where["id"] = 1

	a, err := shardTable.Update(0, data, where)
	if err != nil {
		t.Errorf("update error: %s", err)
	}

	t.Logf("affected: %d", a)

	row, e := shardTable.FetchRow(0, where, &Product{})
	if e != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("product: %v", row)
}

func TestTableShardingDelete(t *testing.T) {
	shardTable = NewTableSharding[*Product]("product", true)
	where := make(map[string]any)
	where["id"] = 1

	a, err := shardTable.Delete(0, where)
	if err != nil {
		t.Errorf("delete error: %s", err)
	}

	t.Logf("affected: %d", a)

	row, e := shardTable.FetchRow(0, where, &Product{})
	if e != nil {
		t.Errorf("err: %s", err)
	}

	t.Logf("product: %v", row)
}
