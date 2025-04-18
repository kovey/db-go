package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTablespace(t *testing.T) {
	ts := NewTablespace().Tablespace("table_spaces")
	ts.OptionOnlyKey("WAIT").OptionStrWith("ADD", "DATAFILE", "data.dt")
	ts.Option("AUTOEXTEND_SIZE", "1M").OptionWith("USE", "LOGFILE GROUP", "logfile_group").OptionStr("COMMENT", "aaaa")
	assert.Equal(t, "CREATE TABLESPACE WAIT ADD DATAFILE 'data.dt' AUTOEXTEND_SIZE 1M USE LOGFILE GROUP logfile_group COMMENT 'aaaa'", ts.Prepare())
	assert.Nil(t, ts.Binds())
}

func TestTablespaceAlter(t *testing.T) {
	ts := NewTablespace().Tablespace("table_spaces").Alter()
	ts.OptionOnlyKey("WAIT").OptionStrWith("ADD", "DATAFILE", "data.dt").OptionStrWith("DROP", "DATAFILE", "old_data.dt")
	ts.Option("AUTOEXTEND_SIZE", "1M").OptionWith("USE", "LOGFILE GROUP", "logfile_group").OptionStr("COMMENT", "aaaa")
	assert.Equal(t, "ALTER TABLESPACE WAIT ADD DATAFILE 'data.dt' DROP DATAFILE 'old_data.dt' AUTOEXTEND_SIZE 1M USE LOGFILE GROUP logfile_group COMMENT 'aaaa'", ts.Prepare())
	assert.Nil(t, ts.Binds())
}
