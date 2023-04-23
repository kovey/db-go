package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/db"
	"github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/tools/desc"
	"github.com/kovey/db-go/v2/tools/meta"
)

func main() {
	host := flag.String("a", "", "database addr")
	port := flag.Int("P", 0, "database port")
	user := flag.String("u", "", "database user name")
	password := flag.String("p", "", "database password")
	dbName := flag.String("db", "", "database name")
	charset := flag.String("c", "", "charset")
	pc := flag.String("pk", "", "package name")
	dist := flag.String("d", "", "dist dir")
	flag.Parse()

	conf := config.Mysql{
		Host:          *host,
		Port:          *port,
		Username:      *user,
		Password:      *password,
		Dbname:        *dbName,
		Charset:       *charset,
		ActiveMax:     100,
		ConnectionMax: 100,
		LifeTime:      100,
	}

	if err := db.Init(conf); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(*dist, 0755); err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	if sta, err := os.Stat(*dist); err != nil {
		panic(err)
	} else {
		if !sta.IsDir() {
			panic(fmt.Errorf("dist[%s] not dir", *dist))
		}
	}

	mysql := db.NewMysql[*desc.Table]()
	tables, err := mysql.ShowTables(sql.NewShowTables(), desc.NewTable())
	if err != nil {
		panic(err)
	}

	dm := db.NewMysql[*desc.Desc]()
	for _, tb := range tables {
		t := meta.NewTable(tb.Name, *pc)
		fields, err := dm.Desc(sql.NewDesc(tb.Name), desc.NewDesc(tb.Name))
		if err != nil {
			panic(err)
		}

		for _, field := range fields {
			if field.Key.String == "PRI" {
				t.SetPrimary(meta.NewField(field.Field, field.Type, false))
			}
			t.Add(meta.NewField(field.Field, field.Type, field.Null != "NO"))
		}

		path := *dist + "/" + tb.Name + ".go"
		os.WriteFile(path, []byte(t.Format()), 0644)
	}
}
