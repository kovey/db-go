package meta

import (
	"strings"
	"time"

	"github.com/kovey/db-go/v2/tools/tpl"
	"github.com/kovey/debug-go/debug"
)

type Table struct {
	Package     string
	Name        string
	DbTableName string
	Fields      []*Field
	Primary     *Field
	HasSql      bool
	HasDecimal  bool
	Database    string
	Comment     string
	Type        string
}

func NewTable(name, comment, p, database, t string) *Table {
	return &Table{DbTableName: name, Name: convert(name), Fields: make([]*Field, 0), Package: p, Database: database, Comment: comment, Type: t}
}

func (t *Table) SetPrimary(f *Field) {
	t.Primary = f
}

func (t *Table) Add(f *Field) {
	if f.HasDecimal {
		t.HasDecimal = true
	}
	if f.HasSql {
		t.HasSql = true
	}

	t.Fields = append(t.Fields, f)
}

func convert(name string) string {
	info := strings.Split(name, "_")
	for i := 0; i < len(info); i++ {
		info[i] = UcFirst(info[i])
	}

	return strings.Join(info, "")
}

func UcFirst(name string) string {
	if name == "" {
		return name
	}

	return strings.ToUpper(name[:1]) + name[1:]
}
func now() string {
	return time.Now().Format(time.DateTime)
}

func (t *Table) Format() string {
	if t.Primary == nil {
		debug.Erro("table[%s] primary not found", t.DbTableName)
		panic(t.DbTableName)
	}

	tplStr := tpl.Table
	if t.Type == "s" {
		tplStr = tpl.Table_Sharding
	}

	content := strings.ReplaceAll(tplStr, "{name}", t.Name)
	content = strings.ReplaceAll(content, "{package_name}", t.Package)
	content = strings.ReplaceAll(content, "{orm_version}", Version)
	content = strings.ReplaceAll(content, "{database}", t.Database)
	content = strings.ReplaceAll(content, "{created_date}", now())
	content = strings.ReplaceAll(content, "{table_name}", t.DbTableName)
	content = strings.ReplaceAll(content, "{table_comment}", t.Comment)
	content = strings.ReplaceAll(content, "{column_const}", t.constaints())
	content = strings.ReplaceAll(content, "{func_columns}", t.columns())
	content = strings.ReplaceAll(content, "{func_fields}", t.tofields())
	content = strings.ReplaceAll(content, "{func_values}", t.values())
	content = strings.ReplaceAll(content, "{row_fields}", t.fields())
	content = strings.ReplaceAll(content, "{primary_id}", t.Primary.ConstaintName(t.Name))
	content = strings.ReplaceAll(content, "{imports}", t.imports())
	content = strings.ReplaceAll(content, "{row_fields_reset}", t.resets())
	switch t.Primary.GolangType {
	case "string":
		content = strings.ReplaceAll(content, "{primary_id_type}", "Str")
	default:
		content = strings.ReplaceAll(content, "{primary_id_type}", "Int")
	}

	closeAuto := ""
	if !t.Primary.IsAutoInc {
		closeAuto = tpl.Close_Auto_Inc
	}

	content = strings.ReplaceAll(content, "{close_auto_inc}", closeAuto)
	return content
}

func (t *Table) resets() string {
	res := make([]string, len(t.Fields))
	for index, f := range t.Fields {
		res[index] = f.GolangDefalut
	}

	return strings.Join(res, "\n")
}

func (t *Table) fields() string {
	res := make([]string, len(t.Fields))
	for index, f := range t.Fields {
		res[index] = f.Format()
	}

	return strings.Join(res, "\n")
}

func (t *Table) constaints() string {
	res := make([]string, len(t.Fields))
	for index, f := range t.Fields {
		res[index] = f.Constaint(t.Name)
	}

	return strings.Join(res, "\n")
}

func (t *Table) columns() string {
	res := make([]string, len(t.Fields))
	for index, f := range t.Fields {
		res[index] = f.ToColumn(t.Name)
	}

	return strings.Join(res, "\n")
}

func (t *Table) tofields() string {
	res := make([]string, len(t.Fields))
	for index, f := range t.Fields {
		res[index] = f.ToField()
	}

	return strings.Join(res, "\n")
}

func (t *Table) values() string {
	res := make([]string, len(t.Fields))
	for index, f := range t.Fields {
		res[index] = f.ToValue()
	}

	return strings.Join(res, "\n")
}

func (t *Table) imports() string {
	res := make([]string, 0)
	if t.HasSql {
		res = append(res, tpl.Sql)
	}
	if t.HasDecimal {
		res = append(res, tpl.Decimal)
	}

	return strings.Join(res, "\n")
}
