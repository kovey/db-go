package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type IndexDrop struct {
	*drop
	table string
	alg   ksql.IndexAlgOption
	lock  ksql.IndexLockOption
}

func NewIndexDrop() *IndexDrop {
	i := &IndexDrop{drop: newDrop("INDEX")}
	i.opChain.Append(i._otherBuild)
	return i
}

func (i *IndexDrop) _otherBuild(builder *strings.Builder) {
	builder.WriteString(" ON")
	operator.BuildColumnString(i.table, builder)
	if i.alg != "" {
		builder.WriteString(" ALGORITHM =")
		operator.BuildPureString(string(i.alg), builder)
	}

	if i.lock != "" {
		builder.WriteString(" LOCK =")
		operator.BuildPureString(string(i.lock), builder)
	}
}

func (s *IndexDrop) Index(event string) ksql.DropIndexInterface {
	s.name = event
	return s
}

func (s *IndexDrop) Table(table string) ksql.DropIndexInterface {
	s.table = table
	return s
}

func (s *IndexDrop) Algorithm(alg ksql.IndexAlgOption) ksql.DropIndexInterface {
	s.alg = alg
	return s
}

func (s *IndexDrop) Lock(lock ksql.IndexLockOption) ksql.DropIndexInterface {
	s.lock = lock
	return s
}
