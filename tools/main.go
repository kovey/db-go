package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/db"
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
	dt := flag.String("t", "n", "database type, n is normal, s is sharding")
	flag.Parse()

	if *dt != "s" && *dt != "n" {
		panic("database type is error")
	}

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

	ttb := desc.NewTableTable()
	tables, err := ttb.FetchAll(ms.Where{"TABLE_SCHEMA": *dbName}, &desc.Table{})
	if err != nil {
		panic(err)
	}

	showTablesBegin(tables, *dt)
	for _, tb := range tables {
		tbName := convertName(tb.Name, *dt)
		debug.Info("process table[%s] orm begin...", tbName)
		t := meta.NewTable(tbName, tb.GetComment(), *pc, *dbName, *dt)
		dTb := desc.NewDescTable()
		fields, err := dTb.FetchAll(ms.Where{"TABLE_SCHEMA": *dbName, "TABLE_NAME": tb.Name}, desc.NewDesc())
		if err != nil {
			panic(err)
		}

		for _, field := range fields {
			if field.Key.String == "PRI" {
				t.SetPrimary(meta.NewField(field.Field, field.Type, field.GetComment(), false))
				if field.Extra.String != "auto_increment" {
					t.Primary.IsAutoInc = false
				}
			}
			t.Add(meta.NewField(field.Field, field.Type, field.GetComment(), field.Null != "NO"))
		}

		path := *dist + "/" + tbName + ".go"
		if err := os.WriteFile(path, []byte(t.Format()), 0644); err != nil {
			debug.Erro("write table[%s] content to file[%s] failure, error: %s", tbName, path, err)
		}

		debug.Info("process table[%s] orm end.", tbName)
	}

	showTablesEnd(tables, *dt)
}

func showTablesBegin(tables []*desc.Table, t string) {
	names := make([]string, len(tables))
	for index, tb := range tables {
		names[index] = convertName(tb.Name, t)
	}

	debug.Info("orm prepare tables[%s]", strings.Join(names, ", "))
}

func showTablesEnd(tables []*desc.Table, t string) {
	names := make([]string, len(tables))
	for index, tb := range tables {
		names[index] = convertName(tb.Name, t)
	}

	debug.Info("orm tables[%s] end.", strings.Join(names, ", "))
}

func convertName(name, t string) string {
	if t == "n" {
		return name
	}

	info := strings.Split(name, "_")
	if len(info) == 1 {
		return name
	}

	_, err := strconv.Atoi(info[len(info)-1])
	if err != nil {
		return name
	}

	return strings.Join(info[:len(info)-1], "_")
}
