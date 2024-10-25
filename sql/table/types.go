package table

import "strings"

func ParseType(t string, length, scale int, sets ...string) *ColumnType {
	t = strings.ToUpper(t)
	var c *ColumnType
	switch t {
	case "BIT", "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "BIGINT", "DATETIME", "TIMESTAMP", "TIME", "YEAR", "CHAR", "VARCHAR", "BINARY", "VARBINARY":
		c = &ColumnType{Name: t, Type: Scale_Type_One, Length: length, Scale: scale}
	case "DECIMAL", "FLOAT", "DOUBLE":
		c = &ColumnType{Name: t, Type: Scale_Type_Two, Length: length, Scale: scale}
	case "DATE", "TINYBLOB", "TINYTEXT", "BLOB", "TEXT", "MEDIUMBLOB", "MEDIUMTEXT", "LONGBLOB", "LONGBTEXT", "GEOMETRY", "POINT", "LINESTRING", "POLYGON", "MULTIPOINT",
		"MULTILINESTRING", "MULTIPOLYGON", "GEOMETRYCOLLECTION", "JSON":
		c = &ColumnType{Name: t, Type: Scale_Type_Zero, Length: length, Scale: scale}
	case "ENUM", "SET":
		c = &ColumnType{Name: t, Type: Scale_Type_More, Length: length, Scale: scale}
	}

	if c != nil {
		c.Set(sets...)
	}

	return c
}
