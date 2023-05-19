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
}

func NewField(name, t, comment string, isNull bool) *Field {
	f := &Field{Name: UcFirst(name), DbField: name, Type: t, IsNull: isNull, Comment: comment}
	f.GolangType = f.parse()
	return f
}

func (f *Field) Format() string {
	return fmt.Sprintf(tpl.Field, f.Name, f.GolangType, f.DbField, f.Comment)
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
