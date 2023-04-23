package meta

import (
	"fmt"
	"strings"

	"github.com/kovey/db-go/v2/tools/tpl"
)

type Field struct {
	Name    string
	Type    string
	DbField string
}

func (f *Field) Format() string {
	return fmt.Sprintf(tpl.Field, f.Name, f.GoType(), f.DbField)
}

func (f *Field) GoType() string {
	if strings.Contains(f.Type, Mysql_Binary) || strings.Contains(f.Type, Mysql_Blob) {
		return "[]byte"
	}

	if strings.Contains(f.Type, Mysql_BigInt) {
		return "int64"
	}

	if strings.Contains(f.Type, Mysql_Int) {
		return "int32"
	}

	if strings.Contains(f.Type, Mysql_Decimal) {
		return "decimal.Decimal"
	}

	if strings.Contains(f.Type, Mysql_Double) {
		return "float64"
	}

	if strings.Contains(f.Type, Mysql_Float) {
		return "float32"
	}

	if strings.Contains(f.Type, Mysql_Bool) {
		return "bool"
	}

	return "string"
}
