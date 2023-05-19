package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/db"
	"github.com/kovey/db-go/v2/sql"
	ms "github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/db-go/v2/tools/desc"
	"github.com/kovey/db-go/v2/tools/meta"
	"github.com/kovey/debug-go/debug"
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

	showTablesBegin(tables)
	for _, tb := range tables {
		debug.Info("process table[%s] orm begin...", tb.Name)
		t := meta.NewTable(tb.Name, *pc, *dbName)
		dTb := desc.NewDescTable()
		fields, err := dTb.FetchAll(ms.Where{"TABLE_SCHEMA": *dbName, "TABLE_NAME": tb.Name}, desc.NewDesc())
		if err != nil {
			panic(err)
		}

		for _, field := range fields {
			if field.Key.String == "PRI" {
				t.SetPrimary(meta.NewField(field.Field, field.Type, field.Comment.String, false))
			}
			t.Add(meta.NewField(field.Field, field.Type, field.Comment.String, field.Null != "NO"))
		}

		path := *dist + "/" + tb.Name + ".go"
		if err := os.WriteFile(path, []byte(t.Format()), 0644); err != nil {
			debug.Erro("write table[%s] content to file[%s] failure, error: %s", tb.Name, path, err)
		}

		debug.Info("process table[%s] orm end.", tb.Name)
	}

	showTablesEnd(tables)
}

func showTablesBegin(tables []*desc.Table) {
	names := make([]string, len(tables))
	for index, tb := range tables {
		names[index] = tb.Name
	}

	debug.Info("orm prepare tables[%s]", strings.Join(names, ", "))
}

func showTablesEnd(tables []*desc.Table) {
	names := make([]string, len(tables))
	for index, tb := range tables {
		names[index] = tb.Name
	}

	debug.Info("orm tables[%s] end.", strings.Join(names, ", "))
}
