package driver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportSavePoint(t *testing.T) {
	assert.True(t, SupportSavePoint("mysql"))
	assert.True(t, SupportSavePoint("postgresql"))
	assert.True(t, SupportSavePoint("oracle"))
	assert.True(t, SupportSavePoint("sqlserver"))
	assert.True(t, SupportSavePoint("db2"))
	assert.True(t, SupportSavePoint("sqlite"))
	assert.True(t, SupportSavePoint("firebird"))
	assert.True(t, SupportSavePoint("h2"))
	assert.False(t, SupportSavePoint("tidb"))
}
