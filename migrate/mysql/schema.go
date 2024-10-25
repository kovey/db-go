package mysql

import (
	"github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/migrate/schema"
)

type Schema struct {
	*base
	SCHEMA_NAME                string
	DEFAULT_CHARACTER_SET_NAME string
	DEFAULT_COLLATION_NAME     string
	tables                     []schema.TableInfoInterface
}

func (s *Schema) SetTables(tables []schema.TableInfoInterface) {
	s.tables = tables
}

func (s *Schema) Name() string {
	return s.SCHEMA_NAME
}

func (s *Schema) Charset() string {
	return s.DEFAULT_CHARACTER_SET_NAME
}

func (s *Schema) Collation() string {
	return s.DEFAULT_COLLATION_NAME
}

func (s *Schema) HasChanged(other schema.SchemaInfoInterface) bool {
	return s.Charset() != other.Charset() || s.Collation() != other.Collation()
}

func (s *Schema) HasTable(table schema.TableInfoInterface) bool {
	for _, ta := range s.tables {
		if ta.Name() == table.Name() {
			return true
		}
	}

	return false
}

func (s *Schema) Tables() []schema.TableInfoInterface {
	return s.tables
}

func (s *Schema) Clone() ksql.RowInterface {
	return &Schema{base: &base{empty: true}}
}

func (s *Schema) Values() []any {
	return []any{&s.SCHEMA_NAME, &s.DEFAULT_CHARACTER_SET_NAME, &s.DEFAULT_COLLATION_NAME}
}

func (s *Schema) Columns() []string {
	return []string{"SCHEMA_NAME", "DEFAULT_CHARACTER_SET_NAME", "DEFAULT_COLLATION_NAME"}
}
