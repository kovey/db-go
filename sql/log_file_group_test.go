package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogFileGroup(t *testing.T) {
	l := NewLogFileGroup().LogFileGroup("kovey_group").UndoFile("data.dat")
	l.InitialSize("10M").Wait().Comment("log file")
	assert.Equal(t, "CREATE LOGFILE GROUP `kovey_group` ADD UNDOFILE 'data.dat' INITIAL_SIZE = 10M WAIT COMMENT 'log file'", l.Prepare())
	assert.Nil(t, l.Binds())
}

func TestLogFileGroupOther(t *testing.T) {
	l := NewLogFileGroup().LogFileGroup("kovey_group").UndoFile("data.dat")
	l.NodeGroupId("11").RedoBufferSize("20M").UndoBufferSize("30M").Engine("engine")
	assert.Equal(t, "CREATE LOGFILE GROUP `kovey_group` ADD UNDOFILE 'data.dat' UNDO_BUFFER_SIZE = 30M REDO_BUFFER_SIZE = 20M NODEGROUP = 11 ENGINE = engine", l.Prepare())
	assert.Nil(t, l.Binds())
}

func TestLogFileGroupAlter(t *testing.T) {
	l := NewLogFileGroup().Alter().LogFileGroup("kovey_group").UndoFile("data.dat")
	l.NodeGroupId("11").RedoBufferSize("20M").UndoBufferSize("30M")
	l.InitialSize("10M").Wait().Comment("log file").Engine("engine")
	assert.Equal(t, "ALTER LOGFILE GROUP `kovey_group` ADD UNDOFILE 'data.dat' INITIAL_SIZE = 10M WAIT ENGINE = engine", l.Prepare())
	assert.Nil(t, l.Binds())
}
