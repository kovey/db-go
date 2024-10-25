package sql

import "github.com/kovey/db-go/v3"

// CREATE SCHEMA `test_db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci ;

// ALTER SCHEMA `skillw`  DEFAULT CHARACTER SET utf8  DEFAULT COLLATE utf8_czech_ci ;

type Schema struct {
	*base
	isCreate bool
	isAlter  bool
	schema   string
	charset  string
	collate  string
}

func NewSchema() *Schema {
	return &Schema{base: &base{hasPrepared: false}}
}

func (s *Schema) Create() ksql.SchemaInterface {
	if s.isAlter {
		return s
	}
	s.isCreate = true
	return s
}

func (s *Schema) Alter() ksql.SchemaInterface {
	if s.isCreate {
		return s
	}
	s.isAlter = true
	return s
}

func (s *Schema) Schema(schema string) ksql.SchemaInterface {
	s.schema = schema
	return s
}

func (s *Schema) Charset(charset string) ksql.SchemaInterface {
	s.charset = charset
	return s
}

func (s *Schema) Collate(collate string) ksql.SchemaInterface {
	s.collate = collate
	return s
}

func (s *Schema) Prepare() string {
	if s.hasPrepared {
		return s.base.Prepare()
	}

	s.hasPrepared = true
	if s.isCreate {
		s.keyword("CREATE SCHEMA ")
	} else if s.isAlter {
		s.keyword("ALTER SCHEMA ")
	}

	Backtick(s.charset, &s.builder)
	if s.charset != "" {
		s.keyword(" DEFAULT CHARACTER SET ")
		s.builder.WriteString(s.charset)
		s.builder.WriteString(" ")
	}

	if s.collate != "" {
		s.keyword("COLLATE ")
		s.builder.WriteString(s.collate)
	}

	return s.base.Prepare()
}
