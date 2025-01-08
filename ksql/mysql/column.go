package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kovey/db-go/ksql/schema"
	ksql "github.com/kovey/db-go/v3"
)

type Column struct {
	*base
	COLUMN_NAME              string
	COLUMN_DEFAULT           sql.NullString
	IS_NULLABLE              string
	DATA_TYPE                string
	CHARACTER_MAXIMUM_LENGTH sql.NullInt64
	NUMERIC_PRECISION        sql.NullInt32
	NUMERIC_SCALE            sql.NullInt32
	DATETIME_PRECISION       sql.NullInt32
	COLUMN_KEY               string
	EXTRA                    string
	COLUMN_COMMENT           sql.NullString
}

func (c *Column) Key() string {
	return c.COLUMN_KEY
}

func (c *Column) Name() string {
	return c.COLUMN_NAME
}

func (c *Column) HasDefault() bool {
	return c.COLUMN_DEFAULT.Valid
}

func (c *Column) AutoIncrement() bool {
	return strings.ToLower(c.EXTRA) == "auto_increment"
}

func (c *Column) Default() string {
	if strings.ToUpper(c.EXTRA) == "DEFAULT_GENERATED ON UPDATE CURRENT_TIMESTAMP" {
		return fmt.Sprintf("%s ON UPDATE %s", c.COLUMN_DEFAULT.String, c.COLUMN_DEFAULT.String)
	}

	return c.COLUMN_DEFAULT.String
}

func (c *Column) Nullable() bool {
	return c.IS_NULLABLE == "YES" || c.IS_NULLABLE == "yes"
}

func (c *Column) Type() string {
	return c.DATA_TYPE
}

func (c *Column) Length() int {
	if c.NUMERIC_PRECISION.Valid {
		return int(c.NUMERIC_PRECISION.Int32)
	}

	return int(c.CHARACTER_MAXIMUM_LENGTH.Int64)
}

func (c *Column) NumLen() int {
	return int(c.NUMERIC_PRECISION.Int32)
}

func (c *Column) Scale() int {
	return int(c.NUMERIC_SCALE.Int32)
}

func (c *Column) DateTimeLen() int {
	return int(c.DATETIME_PRECISION.Int32)
}

func (c *Column) Comment() string {
	return c.COLUMN_COMMENT.String
}

func (c *Column) Extra() string {
	return c.EXTRA
}

func (c *Column) HasChanged(other schema.ColumnInfoInterface) bool {
	return c.Name() != other.Name() || c.Default() != other.Default() || c.Nullable() != other.Nullable() || c.Type() != other.Type() || c.Length() != other.Length() || c.NumLen() != other.NumLen() || c.Scale() != other.Scale() || c.DateTimeLen() != other.DateTimeLen() || c.Comment() != other.Comment() || c.Extra() != other.Extra()
}

func (c *Column) Values() []any {
	return []any{&c.COLUMN_NAME, &c.COLUMN_DEFAULT, &c.IS_NULLABLE, &c.DATA_TYPE, &c.CHARACTER_MAXIMUM_LENGTH, &c.NUMERIC_PRECISION, &c.NUMERIC_SCALE, &c.DATETIME_PRECISION, &c.COLUMN_KEY, &c.EXTRA, &c.COLUMN_COMMENT}
}

func (c *Column) Clone() ksql.RowInterface {
	return &Column{base: &base{isInitialized: true}}
}

func (c *Column) Columns() []string {
	return []string{"COLUMN_NAME", "COLUMN_DEFAULT", "IS_NULLABLE", "DATA_TYPE", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE", "DATETIME_PRECISION", "COLUMN_KEY", "EXTRA", "COLUMN_COMMENT"}
}
