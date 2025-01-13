package orm

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kovey/db-go/ksql/dir"
	"github.com/kovey/db-go/ksql/mysql"
	"github.com/kovey/db-go/ksql/schema"
	"github.com/kovey/db-go/ksql/version"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/debug-go/debug"
)

func Orm(driverName, dsn, directory, dbname string) error {
	if err := os.MkdirAll(directory, 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	if err := db.Init(db.Config{DriverName: driverName, DataSourceName: dsn, MaxIdleTime: 120 * time.Second, MaxLifeTime: 120 * time.Second, MaxIdleConns: 10, MaxOpenConns: 10}); err != nil {
		return err
	}

	info := strings.Split(directory, dir.Sep())
	packageName := info[len(info)-1]
	conn, _ := db.Get()
	tables, err := getTables(context.Background(), conn, dbname)
	if err != nil {
		return err
	}

	debug.Info("create orm from db[%s] begin...", dbname)
	defer debug.Info("create orm from db[%s] end.", dbname)
	ctx := context.Background()
	for _, table := range tables {
		if table.Name() == "ksql_migrate_info" {
			continue
		}

		debug.Info("create table[%s] begin...", table.Name())
		defer debug.Info("create table[%s] end.", table.Name())
		ddlSql, err := db.ShowDDLBy(ctx, conn, table.Name())
		tpl := &modelTpl{Package: packageName, Table: table.Name(), Comment: table.Comment(), CreateTime: time.Now().Format(time.DateTime), Version: version.Version(), Sql: ddlSql}
		tpl.ModelTag = tag("-")
		name := formatName(table.Name())
		tpl.Name = name
		tpl.DbName = dbname
		var columns []string
		var values []string
		for _, column := range table.Fields() {
			f := field{Name: formatName(column.Name()), Comment: column.Comment(), Type: getType(column.Type()), Tag: tag(column.Name())}
			f.CanNull = column.Nullable()
			tpl.Fields = append(tpl.Fields, f)
			tpl.Consts = append(tpl.Consts, constInfo{Table: tpl.Name, Name: column.Name(), Column: f.Name, Comment: f.Comment})
			if f.Type == "time.Time" {
				tpl.Imports = append(tpl.Imports, "time")
			}
			if column.Key() == "PRI" {
				tpl.PrimaryId = column.Name()
				if f.Type == "string" {
					tpl.PrimaryType = "Type_Str"
				} else {
					tpl.PrimaryType = "Type_Int"
				}
			}

			columns = append(columns, fmt.Sprintf("Table_%s_%s", tpl.Name, f.Name))
			values = append(values, "&self."+formatName(column.Name()))
		}

		tpl.Columns = strings.Join(columns, ",")
		tpl.Values = strings.Join(values, ",")
		res, err := tpl.Parse()
		if err != nil {
			return err
		}

		if err := os.WriteFile(directory+dir.Sep()+table.Name()+".go", res, 0644); err != nil {
			return err
		}

		debug.Info("create table[%s] success.", table.Name())
	}

	return nil
}

func tag(name string) string {
	var tag strings.Builder
	tag.WriteString("`db:")
	tag.WriteByte('"')
	tag.WriteString(name)
	tag.WriteByte('"')
	tag.WriteString(" json:")
	tag.WriteByte('"')
	tag.WriteString(name)
	tag.WriteByte('"')
	tag.WriteString("`")
	return tag.String()
}

func getType(mysqlType string) string {
	switch strings.ToUpper(mysqlType) {
	case "BIT":
		return "byte"
	case "TINYINT":
		return "int8"
	case "SMALLINT":
		return "int16"
	case "MEDIUMINT", "INT":
		return "int"
	case "BIGINT":
		return "int64"
	case "DATETIME":
		return "time.Time"
	case "TIMESTAMP":
		return "time.Time"
	case "TIME":
		return "time.Time"
	case "YEAR":
		return "int16"
	case "CHAR":
		return "byte"
	case "VARCHAR":
		return "string"
	case "BINARY":
		return "[]byte"
	case "VARBINARY":
		return "[]byte"
	case "DECIMAL":
		return "float64"
	case "FLOAT":
		return "float32"
	case "DOUBLE":
		return "float64"
	case "DATE":
		return "time.Time"
	case "TINYBLOB", "TINYTEXT", "BLOB", "TEXT", "MEDIUMBLOB", "MEDIUMTEXT", "LONGBLOB", "LONGBTEXT", "GEOMETRY", "POINT", "LINESTRING", "POLYGON", "MULTIPOINT",
		"MULTILINESTRING", "MULTIPOLYGON", "GEOMETRYCOLLECTION", "JSON":
		return "string"
	case "ENUM":
		return "int16"
	case "SET":
		return "string"
	}

	return "string"
}

func getTables(ctx context.Context, conn ksql.ConnectionInterface, dbname string) ([]schema.TableInfoInterface, error) {
	switch strings.ToLower(conn.DriverName()) {
	case "mysql":
		return mysql.Tables(ctx, conn, dbname)
	}

	return nil, nil
}

func formatName(name string) string {
	info := strings.Split(name, "_")
	for i := 0; i < len(info); i++ {
		info[i] = FirstUpper(info[i])
	}

	return strings.Join(info, "")
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
