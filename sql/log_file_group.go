package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type LogFileGroup struct {
	*base
	name           string
	undoFile       string
	initialSize    string
	undoBufferSize string
	redoBuffSize   string
	nodeGroupId    string
	wait           string
	comment        string
	engine         string
	isCreate       bool
}

func NewLogFileGroup() *LogFileGroup {
	l := &LogFileGroup{base: newBase(), isCreate: true}
	l.opChain.Append(l._keyword, l._name, l._size, l._other)
	return l
}

func (l *LogFileGroup) _keyword(builder *strings.Builder) {
	if l.isCreate {
		keywordCreate(builder)
		return
	}

	keywordAlter(builder)
}

func (l *LogFileGroup) _name(builder *strings.Builder) {
	builder.WriteString(" LOGFILE GROUP ")
	Backtick(l.name, builder)
	builder.WriteString(" ADD UNDOFILE ")
	Quote(l.undoFile, builder)
}

func (l *LogFileGroup) _size(builder *strings.Builder) {
	if l.initialSize != "" {
		builder.WriteString(" INITIAL_SIZE = ")
		builder.WriteString(l.initialSize)
	}

	if !l.isCreate {
		return
	}

	if l.undoBufferSize != "" {
		builder.WriteString(" UNDO_BUFFER_SIZE = ")
		builder.WriteString(l.undoBufferSize)
	}

	if l.redoBuffSize != "" {
		builder.WriteString(" REDO_BUFFER_SIZE = ")
		builder.WriteString(l.redoBuffSize)
	}
}

func (l *LogFileGroup) _other(builder *strings.Builder) {
	if l.isCreate && l.nodeGroupId != "" {
		builder.WriteString(" NODEGROUP = ")
		builder.WriteString(l.nodeGroupId)
	}

	if l.wait != "" {
		builder.WriteString(" ")
		builder.WriteString(l.wait)
	}

	if l.isCreate && l.comment != "" {
		builder.WriteString(" COMMENT ")
		Quote(l.comment, builder)
	}

	if l.engine != "" {
		builder.WriteString(" ENGINE = ")
		builder.WriteString(l.engine)
	}
}

func (l *LogFileGroup) LogFileGroup(name string) ksql.LogFileGroupInterface {
	l.name = name
	return l
}

func (l *LogFileGroup) UndoFile(file string) ksql.LogFileGroupInterface {
	l.undoFile = file
	return l
}

func (l *LogFileGroup) InitialSize(size string) ksql.LogFileGroupInterface {
	l.initialSize = size
	return l
}

func (l *LogFileGroup) UndoBufferSize(size string) ksql.LogFileGroupInterface {
	l.undoBufferSize = size
	return l
}

func (l *LogFileGroup) RedoBufferSize(size string) ksql.LogFileGroupInterface {
	if !l.isCreate {
		return l
	}

	l.redoBuffSize = size
	return l
}

func (l *LogFileGroup) NodeGroupId(nodeGroupId string) ksql.LogFileGroupInterface {
	if !l.isCreate {
		return l
	}

	l.nodeGroupId = nodeGroupId
	return l
}

func (l *LogFileGroup) Wait() ksql.LogFileGroupInterface {
	l.wait = "WAIT"
	return l
}

func (l *LogFileGroup) Comment(comment string) ksql.LogFileGroupInterface {
	if !l.isCreate {
		return l
	}

	l.comment = comment
	return l
}

func (l *LogFileGroup) Engine(engine string) ksql.LogFileGroupInterface {
	l.engine = engine
	return l
}

func (l *LogFileGroup) Alter() ksql.LogFileGroupInterface {
	l.isCreate = false
	return l
}
