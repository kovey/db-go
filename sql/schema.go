package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

// CREATE SCHEMA `test_db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci ;
// ALTER SCHEMA `skillw`  DEFAULT CHARACTER SET utf8  DEFAULT COLLATE utf8_czech_ci ;

type Schema struct {
	*base
	isCreate    bool
	isAlter     bool
	schema      string
	charset     string
	collate     string
	encryption  string
	readOnly    string
	ifNotExists bool
}

func NewSchema() *Schema {
	return &Schema{base: newBase()}
}

func keywordAlter(builder *strings.Builder) {
	builder.WriteString("ALTER")
}

func keywordCreate(builder *strings.Builder) {
	builder.WriteString("CREATE")
}

func keywordSchema(builder *strings.Builder) {
	builder.WriteString(" SCHEMA")
}

func (s *Schema) IfNotExists() ksql.SchemaInterface {
	if s.isAlter {
		return s
	}

	s.ifNotExists = true
	return s
}

func (s *Schema) _schema(builder *strings.Builder) {
	if s.ifNotExists {
		builder.WriteString(" IF NOT EXISTS")
	}
	builder.WriteString(" ")
	Backtick(s.schema, builder)
	if s.isCreate {
		builder.WriteString(" DEFAULT")
	}
}

func (s *Schema) _charset(buidler *strings.Builder) {
	if s.charset == "" {
		return
	}

	if s.isAlter {
		buidler.WriteString(" DEFAULT CHARACTER SET = ")
	} else {
		buidler.WriteString(" CHARACTER SET = ")
	}

	buidler.WriteString(s.charset)
}

func (s *Schema) _collate(buidler *strings.Builder) {
	if s.collate == "" {
		return
	}

	if s.isAlter {
		buidler.WriteString(" DEFAULT COLLATE = ")
	} else {
		buidler.WriteString(" COLLATE = ")
	}
	buidler.WriteString(s.collate)
}

func (s *Schema) _encryption(buidler *strings.Builder) {
	if s.encryption == "" {
		return
	}

	if s.isAlter {
		buidler.WriteString(" DEFAULT ENCRYPTION = ")
	} else {
		buidler.WriteString(" ENCRYPTION = ")
	}
	Quote(s.encryption, buidler)
}

func (s *Schema) _readOnly(buidler *strings.Builder) {
	if s.isCreate || s.readOnly == "" {
		return
	}

	buidler.WriteString(" READ ONLY = ")
	buidler.WriteString(s.readOnly)
}

func (s *Schema) Create() ksql.SchemaInterface {
	if s.isAlter {
		return s
	}

	s.opChain.Append(keywordCreate, keywordSchema, s._schema, s._charset, s._collate, s._encryption, s._readOnly)
	s.isCreate = true
	return s
}

func (s *Schema) Alter() ksql.SchemaInterface {
	if s.isCreate {
		return s
	}
	s.opChain.Append(keywordAlter, keywordSchema, s._schema, s._charset, s._collate, s._encryption, s._readOnly)
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

func (s *Schema) Encryption(encryption string) ksql.SchemaInterface {
	s.encryption = encryption
	return s
}

func (s *Schema) ReadOnly(readOnly string) ksql.SchemaInterface {
	s.readOnly = readOnly
	return s
}
