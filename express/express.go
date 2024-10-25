package express

type Statement struct {
	raw   string
	binds []any
}

func NewStatement(raw string, binds []any) *Statement {
	return &Statement{raw: raw, binds: binds}
}

func (s *Statement) Statement() string {
	return s.raw
}

func (s *Statement) Binds() []any {
	return s.binds
}
