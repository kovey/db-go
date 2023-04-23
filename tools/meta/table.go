package meta

import (
	"strings"

	"github.com/kovey/db-go/v2/tools/tpl"
)

type Table struct {
	Package     string
	Name        string
	DbTableName string
	Fields      []*Field
	Primary     *Field
	HasSql      bool
	HasDecimal  bool
}

func NewTable(name string, p string) *Table {
	return &Table{DbTableName: name, Name: UcFirst(name), Fields: make([]*Field, 0), Package: p}
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

func UcFirst(name string) string {
	if name == "" {
		return name
	}

	return strings.ToUpper(name[:1]) + name[1:]
}

func (t *Table) Format() string {
	content := strings.ReplaceAll(tpl.Table, "{name}", t.Name)
	content = strings.ReplaceAll(content, "{package_name}", t.Package)
	content = strings.ReplaceAll(content, "{table_name}", t.DbTableName)
	content = strings.ReplaceAll(content, "{row_fields}", t.fields())
	content = strings.ReplaceAll(content, "{primary_id}", t.Primary.DbField)
	content = strings.ReplaceAll(content, "{imports}", t.imports())
	switch t.Primary.GolangType {
	case "string":
		content = strings.ReplaceAll(content, "{primary_id_type}", "Str")
	default:
		content = strings.ReplaceAll(content, "{primary_id_type}", "Int")
	}
	return content
}

func (t *Table) fields() string {
	res := make([]string, len(t.Fields))
	for index, f := range t.Fields {
		res[index] = f.Format()
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
