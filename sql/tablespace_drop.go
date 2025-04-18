package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type TablespaceDrop struct {
	*drop
	engine string
}

func NewTablespaceDrop() *TablespaceDrop {
	i := &TablespaceDrop{drop: newDrop("TABLESPACE")}
	i.opChain.Append(i._otherBuild)
	return i
}

func (i *TablespaceDrop) _otherBuild(builder *strings.Builder) {
	if i.engine != "" {
		builder.WriteString(" ENGINE =")
		operator.BuildPureString(i.engine, builder)
	}
}

func (s *TablespaceDrop) Tablespace(event string) ksql.DropTablespaceInterface {
	s.name = event
	return s
}

func (s *TablespaceDrop) Engine(engine string) ksql.DropTablespaceInterface {
	s.engine = engine
	return s
}

func (s *TablespaceDrop) Undo() ksql.DropTablespaceInterface {
	s.keyword = "UNDO TABLESPACE"
	return s
}
