package table

import "strings"

const (
	Type_Bit                = "BIT"
	Type_TinyInt            = "TINYINT"
	Type_SmallInt           = "SMALLINT"
	Type_MediumInt          = "MEDIUMINT"
	Type_Int                = "INT"
	Type_BigInt             = "BIGINT"
	Type_DateTime           = "DATETIME"
	Type_Timestamp          = "TIMESTAMP"
	Type_Time               = "TIME"
	Type_Year               = "YEAR"
	Type_Char               = "CHAR"
	Type_VarChar            = "VARCHAR"
	Type_Binary             = "BINARY"
	Type_VarBinary          = "VARBINARY"
	Type_Decimal            = "DECIMAL"
	Type_Float              = "FLOAT"
	Type_Double             = "DOUBLE"
	Type_Date               = "DATE"
	Type_TinyBlob           = "TINYBLOB"
	Type_TinyText           = "TINYTEXT"
	Type_Blob               = "BLOB"
	Type_Text               = "TEXT"
	Type_MediumBlob         = "MEDIUMBLOB"
	Type_MediumText         = "MEDIUMTEXT"
	Type_LongBlob           = "LONGBLOB"
	Type_LongText           = "LONGTEXT"
	Type_GeoMetry           = "GEOMETRY"
	Type_Point              = "POINT"
	Type_LineString         = "LINESTRING"
	Type_Polygon            = "POLYGON"
	Type_MultiPoint         = "MULTIPOINT"
	Type_MultiLineString    = "MULTILINESTRING"
	Type_MultiPolygon       = "MULTIPOLYGON"
	Type_GeoMetryCollection = "GEOMETRYCOLLECTION"
	Type_Json               = "JSON"
	Type_Enum               = "ENUM"
	Type_Set                = "SET"
)

func ParseType(t string, length, scale int, sets ...string) *ColumnType {
	t = strings.ToUpper(t)
	var c *ColumnType
	switch t {
	case Type_Bit, Type_TinyInt, Type_SmallInt, Type_MediumInt, Type_Int, Type_BigInt, Type_DateTime, Type_Timestamp, Type_Time,
		Type_Year, Type_Char, Type_VarChar, Type_Binary, Type_VarBinary:
		c = &ColumnType{Name: t, Type: Scale_Type_One, Length: length, Scale: scale}
	case Type_Decimal, Type_Float, Type_Double:
		c = &ColumnType{Name: t, Type: Scale_Type_Two, Length: length, Scale: scale}
	case Type_Date, Type_TinyBlob, Type_TinyText, Type_Blob, Type_Text, Type_MediumBlob, Type_MediumText, Type_LongBlob,
		Type_LongText, Type_GeoMetry, Type_Point, Type_LineString, Type_Polygon, Type_MultiPoint,
		Type_MultiLineString, Type_MultiPolygon, Type_GeoMetryCollection, Type_Json:
		c = &ColumnType{Name: t, Type: Scale_Type_Zero, Length: length, Scale: scale}
	case Type_Enum, Type_Set:
		c = &ColumnType{Name: t, Type: Scale_Type_More, Length: length, Scale: scale}
	}

	if c != nil {
		c.Set(sets...)
	}

	return c
}

func isInteger(t string) bool {
	switch strings.ToUpper(t) {
	case Type_TinyInt, Type_SmallInt, Type_MediumInt, Type_Int, Type_BigInt:
		return true
	default:
		return false
	}
}

func isNumeric(t string) bool {
	if isInteger(t) {
		return true
	}

	switch strings.ToUpper(t) {
	case Type_Decimal, Type_Float, Type_Double:
		return true
	default:
		return false
	}
}
