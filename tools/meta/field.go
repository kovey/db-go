package meta

import (
	"fmt"
	"strings"

	"github.com/kovey/db-go/v2/tools/tpl"
)

type Field struct {
	Name          string
	Type          string
	DbField       string
	IsNull        bool
	HasSql        bool
	HasDecimal    bool
	GolangType    string
	Comment       string
	IsAutoInc     bool
	GolangDefalut string
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
	return fmt.Sprintf(`	%s = "%s" %s`, f.ConstaintName(table), f.DbField, f.Comment)
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
		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Nil, f.Name)
		return "[]byte"
	}

	if strings.Contains(f.Type, Mysql_BigInt) {
		if f.IsNull {
			f.HasSql = true
			f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Sql_Null, f.Name, "Int64")
			return "ds.NullInt64"
		}

		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Num, f.Name)
		return "int64"
	}

	if strings.Contains(f.Type, Mysql_Int) {
		if f.IsNull {
			f.HasSql = true
			f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Sql_Null, f.Name, "Int32")
			return "ds.NullInt32"
		}

		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Num, f.Name)
		return "int32"
	}

	if strings.Contains(f.Type, Mysql_TinyInt) {
		if f.IsNull {
			f.HasSql = true
			f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Sql_Null, f.Name, "Int8")
			return "ds.NullInt8"
		}

		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Num, f.Name)
		return "int8"
	}

	if strings.Contains(f.Type, Mysql_SmallInt) {
		if f.IsNull {
			f.HasSql = true
			f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Sql_Null, f.Name, "Int16")
			return "ds.NullInt16"
		}

		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Num, f.Name)
		return "int16"
	}

	if strings.Contains(f.Type, Mysql_MediumInt) {
		if f.IsNull {
			f.HasSql = true
			f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Sql_Null, f.Name, "Int")
			return "ds.NullInt"
		}

		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Num, f.Name)
		return "int"
	}

	if strings.Contains(f.Type, Mysql_Decimal) {
		f.HasDecimal = true
		if f.IsNull {
			f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Decimal_Null, f.Name)
			return "decimal.NullDecimal"
		}

		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Decimal, f.Name)
		return "decimal.Decimal"
	}

	if strings.Contains(f.Type, Mysql_Double) {
		if f.IsNull {
			f.HasSql = true
			f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Sql_Null, f.Name, "Float64")
			return "ds.NullFloat64"
		}

		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Num, f.Name)
		return "float64"
	}

	if strings.Contains(f.Type, Mysql_Float) {
		if f.IsNull {
			f.HasSql = true
			f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Sql_Null, f.Name, "Float32")
			return "ds.NullFloat32"
		}

		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Num, f.Name)
		return "float32"
	}

	if strings.Contains(f.Type, Mysql_Bool) {
		if f.IsNull {
			f.HasSql = true
			f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Sql_Null, f.Name, "Bool")
			return "ds.NullBool"
		}

		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Bool, f.Name)
		return "bool"
	}

	if f.IsNull {
		f.HasSql = true
		f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Sql_Null, f.Name, "String")
		return "ds.NullString"
	}

	f.GolangDefalut = fmt.Sprintf(tpl.Field_Reset_Str, f.Name)
	return "string"
}
