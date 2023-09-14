package meta

import (
	"fmt"
	"strings"

	"github.com/kovey/db-go/v2/tools/tpl"
)

type Field struct {
	Name       string
	Type       string
	DbField    string
	IsNull     bool
	HasSql     bool
	HasDecimal bool
	GolangType string
	Comment    string
	IsAutoInc  bool
}

func NewField(name, t, comment string, isNull bool) *Field {
	f := &Field{Name: convert(name), DbField: name, Type: t, IsNull: isNull, Comment: comment, IsAutoInc: true}
	f.GolangType = f.parse()
	return f
}

func (f *Field) Format() string {
	return fmt.Sprintf(tpl.Field, f.Name, f.GolangType, f.DbField, f.Comment)
}

func (f *Field) ConstaintName(table string) string {
	return fmt.Sprintf("Column_%s_%s", table, f.Name)
}

func (f *Field) Constaint(table string) string {
	return fmt.Sprintf(`	%s = "%s"`, f.ConstaintName(table), f.DbField)
}

func (f *Field) ToColumn(table string) string {
	return fmt.Sprintf(tpl.Meta_Column, f.ConstaintName(table))
}

func (f *Field) ToField() string {
	return fmt.Sprintf(tpl.Meta_Fields, f.Name)
}

func (f *Field) ToValue() string {
	return fmt.Sprintf(tpl.Meta_Values, f.Name)
}

func (f *Field) parse() string {
	if strings.Contains(f.Type, Mysql_Binary) || strings.Contains(f.Type, Mysql_Blob) {
		return "[]byte"
	}

	if strings.Contains(f.Type, Mysql_BigInt) {
		if f.IsNull {
			f.HasSql = true
			return "sql.NullInt64"
		}

		return "int64"
	}

	if strings.Contains(f.Type, Mysql_Int) {
		if f.IsNull {
			f.HasSql = true
			return "sql.NullInt32"
		}

		return "int32"
	}

	if strings.Contains(f.Type, Mysql_Decimal) {
		f.HasDecimal = true
		return "decimal.Decimal"
	}

	if strings.Contains(f.Type, Mysql_Double) {
		if f.IsNull {
			f.HasSql = true
			return "sql.NullFloat64"
		}

		return "float64"
	}

	if strings.Contains(f.Type, Mysql_Float) {
		if f.IsNull {
			f.HasSql = true
			return "sql.NullFloat64"
		}

		return "float32"
	}

	if strings.Contains(f.Type, Mysql_Bool) {
		if f.IsNull {
			f.HasSql = true
			return "sql.NullBool"
		}

		return "bool"
	}

	if f.IsNull {
		f.HasSql = true
		return "sql.NullString"
	}

	return "string"
}
