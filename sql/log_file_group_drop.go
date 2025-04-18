package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type LogFileGroupDrop struct {
	*drop
	engine string
}

func NewLogFileGroupDrop() *LogFileGroupDrop {
	i := &LogFileGroupDrop{drop: newDrop("LOGFILE GROUP")}
	i.opChain.Append(i._otherBuild)
	return i
}

func (i *LogFileGroupDrop) _otherBuild(builder *strings.Builder) {
	if i.engine != "" {
		builder.WriteString(" ENGINE =")
		operator.BuildPureString(i.engine, builder)
	}
}

func (s *LogFileGroupDrop) LogFileGroup(event string) ksql.DropLogFileGroupInterface {
	s.name = event
	return s
}

func (s *LogFileGroupDrop) Engine(engine string) ksql.DropLogFileGroupInterface {
	s.engine = engine
	return s
}
