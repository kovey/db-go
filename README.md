## kovey mysql database by golang
#### Description
###### This is a mysql database library
###### Usage
    go get -u github.com/kovey/db-go
### Examples
```golang
    package main
    import (
        "fmt"
        "os"
        "strings"
        "testing"

        "github.com/kovey/config-go/config"
        "github.com/kovey/db-go/db"
        "github.com/kovey/db-go/table"
        "github.com/kovey/db-go/model"
    )

    type ProTable struct {
        table.Table
    }

    type Product struct {
        Base    model.Base
        Id      int    `db:"id"`
        Name    string `db:"name"`
        Date    string `db:"date"`
        Time    string `db:"time"`
        Sex     int    `db:"sex"`
        Content string `db:"content"`
    }

    func NewProTable() *ProTable {
        return &ProTable{*table.NewTable("product")}
    }

    func NewProduct() Product {
        pro := Product{model.NewBase(NewProTable(), model.NewPrimaryId("id", model.Int)), 0, "", "", "", 0, "{}"}

        return pro
    }

    func TestModelDelete(t *testing.T) {
    }

    func TestMain(m *testing.M) {
        setup()
        code := m.Run()
        teardown()
        os.Exit(code)
    }


    func main() {
        conf := config.Mysql{
            Host: "127.0.0.1", Port: 3306, Username: "root", Password: "root", Dbname: "test", Charset: "utf8mb4", ActiveMax: 10, ConnectionMax: 10,
        }
        err := db.Init(conf)
        if err != nil {
            fmt.Printf("init mysql error: %s", err)
        }

        mysql := db.NewMysql()
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

        // save
        pro := NewProduct()
        pro.Name = "kovey"
        pro.Date = "2021-08-12"
        pro.Time = "2021-08-12 13:12:12"
        pro.Sex = 1
        pro.Content = "{\"where\":123}"

        if err := pro.Save(&pro); err != nil {
            fmt.Printf("product save fail, error: %s\n", err)
        }

        fmt.Printf("id: %d\n", pro.Id)

        pro1 := NewProduct()
        where := make(map[string]interface{})
        where["id"] = pro.Id

        // update
        pro1.FetchRow(where, &pro1)
        pro1.Name = "chelsea"
        pro1.Save(&pro1)

        // select
        where = make(map[string]interface{})
        where["id"] = 1
        pr1 = NewProduct()
        if err := pr1.FetchRow(where, &pr1); err != nil {
            fmt.Printf("fetch row err: %s\n", err)
        }

        fmt.Printf("pr1: %v\n", pr1)
        
        // delete
        where = make(map[string]interface{})
        where["id"] = 1
        pr1 = NewProduct()
        if err := pr1.FetchRow(where, &pr1); err != nil {
            t.Errorf("fetch row err: %s", err)
        }

        if err = pr1.Delete(pr1); err != nil {
            fmt.Printf("delete row err: %s\n", err)
        }

        pr2 := NewProduct()
        pr2.FetchRow(where, &pr2)
        fmt.Printf("pr2: %v\n", pr2)
    }
```
