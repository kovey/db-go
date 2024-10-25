package sql

import "strings"

type base struct {
	binds       []any
	builder     strings.Builder
	hasPrepared bool
}

func (b *base) keyword(keyword string) {
	b.builder.WriteString(keyword)
}

func (b *base) Prepare() string {
	return b.builder.String()
}

func (b *base) Binds() []any {
	return b.binds
}
