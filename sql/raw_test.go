package sql

import (
	"fmt"
	"strings"
	"testing"
	"time"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

type rawTestString struct {
	name string
}

func (r *rawTestString) String() string {
	return r.name
}

func TestRaw(t *testing.T) {
	raw := Raw("SELECT * FROM `user` WHERE `id` > ?", 1)
	assert.Equal(t, "SELECT * FROM `user` WHERE `id` > ?", raw.Statement())
	assert.Equal(t, []any{1}, raw.Binds())
	assert.False(t, raw.IsExec())
}

func TestRawValue(t *testing.T) {
	var val = ""
	val = RawValue("kovey")
	assert.Equal(t, "'kovey'", val)
	val = RawValue(100)
	assert.Equal(t, "100", val)
	val = RawValue(100.01)
	assert.Equal(t, "100.010000", val)
	val = RawValue(&rawTestString{name: "kovey"})
	assert.Equal(t, "'kovey'", val)
	val = RawValue(true)
	assert.Equal(t, "true", val)
}

func TestRawBacktick(t *testing.T) {
	var build strings.Builder
	Backtick("kovey", &build)
	assert.Equal(t, "`kovey`", build.String())
	Backtick("(other)", &build)
	assert.Equal(t, "`kovey`(other)", build.String())
}

func TestRawColumn(t *testing.T) {
	var build strings.Builder
	Column("kovey", &build)
	assert.Equal(t, "`kovey`", build.String())
	Column("user.name", &build)
	assert.Equal(t, "`kovey``user`.`name`", build.String())
	Column("(other)", &build)
	assert.Equal(t, "`kovey``user`.`name`(other)", build.String())
}

func TestRawFormatSharding(t *testing.T) {
	assert.Equal(t, "(kovey)", _formatSharding("(kovey)", ksql.Sharding_Day))
	assert.Equal(t, fmt.Sprintf("kovey_%s", time.Now().Format(ksql.Day_Format)), _formatSharding("kovey", ksql.Sharding_Day))
	assert.Equal(t, fmt.Sprintf("kovey_%s", time.Now().Format(ksql.Month_Format)), _formatSharding("kovey", ksql.Sharding_Month))
	assert.Equal(t, "kovey", _formatSharding("kovey", ksql.Sharding_None))
}
