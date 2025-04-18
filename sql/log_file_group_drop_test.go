package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogFileGroupDrop(t *testing.T) {
	s := NewLogFileGroupDrop().LogFileGroup("user").Engine("users")
	assert.Equal(t, "DROP LOGFILE GROUP `user` ENGINE = users", s.Prepare())
	assert.Nil(t, s.Binds())
}
