package express

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type Statement struct {
	raw   string
	binds []any
	typ   ksql.SqlType
}

func NewStatement(raw string, binds []any) *Statement {
	s := &Statement{raw: raw, binds: binds, typ: ksql.Sql_Type_Query}
	s.parse()
	return s
}

func (s *Statement) parse() {
	count := len(s.raw)
	if count < 4 {
		return
	}
	first := ksql.SqlType(strings.ToUpper(s.raw[:4]))
	if first == ksql.Sql_Type_Drop {
		s.typ = ksql.Sql_Type_Drop
		return
	}

	if count < 5 {
		return
	}
	first = ksql.SqlType(strings.ToUpper(s.raw[:5]))
	switch first {
	case ksql.Sql_Type_Alter, ksql.Sql_Type_Query:
		s.typ = first
		return
	}

	if count < 6 {
		return
	}

	first = ksql.SqlType(strings.ToUpper(s.raw[:6]))
	switch first {
	case ksql.Sql_Type_Insert, ksql.Sql_Type_Update, ksql.Sql_Type_Delete, ksql.Sql_Type_Create:
		s.typ = first
	}

	if count < 7 {
		return
	}
	first = ksql.SqlType(strings.ToUpper(s.raw[:7]))
	switch first {
	case ksql.Sql_Type_Release:
		s.typ = first
	}

	if count < 8 {
		return
	}
	first = ksql.SqlType(strings.ToUpper(s.raw[:8]))
	switch first {
	case ksql.Sql_Type_Rollback:
		s.typ = first
	}

	if count < 9 {
		return
	}

	first = ksql.SqlType(strings.ToUpper(s.raw[:9]))
	switch first {
	case ksql.Sql_Type_Save_Point:
		s.typ = first
	}
}

func (s *Statement) Statement() string {
	return s.raw
}

func (s *Statement) Binds() []any {
	return s.binds
}

func (s *Statement) IsExec() bool {
	return s.typ != ksql.Sql_Type_Query
}

func (s *Statement) Type() ksql.SqlType {
	return s.typ
}
