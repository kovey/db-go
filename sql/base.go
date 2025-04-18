package sql

import (
	"strings"

	"github.com/kovey/db-go/v3/sql/operator"
)

type base struct {
	binds       []any
	builder     strings.Builder
	hasPrepared bool
	opChain     *operator.Chain
}

func newBase() *base {
	return &base{opChain: operator.NewChain()}
}

func (b *base) keyword(keyword string) {
	b.builder.WriteString(keyword)
}

func (b *base) Prepare() string {
	if b.opChain != nil {
		b.opChain.Call(&b.builder)
	}

	return b.builder.String()
}

func (b *base) Binds() []any {
	return b.binds
}
