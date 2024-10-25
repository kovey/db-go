package db

import "fmt"

type SqlErr struct {
	Sql   string
	Binds []any
	Err   error
}

func (s *SqlErr) Error() string {
	if s.Err == nil {
		return ""
	}

	return fmt.Sprintf("sql: %s, binds: %v, error: %s", s.Sql, s.Binds, s.Err)
}
