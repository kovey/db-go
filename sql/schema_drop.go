package sql

import (
	ksql "github.com/kovey/db-go/v3"
)

type SchemaDrop struct {
	*drop
}

func NewSchemaDrop() *SchemaDrop {
	return &SchemaDrop{drop: newDrop("SCHEMA")}
}

func (s *SchemaDrop) Schema(schema string) ksql.DropSchemaInterface {
	s.name = schema
	return s
}

func (s *SchemaDrop) IfExists() ksql.DropSchemaInterface {
	s.ifExists = true
	return s
}
