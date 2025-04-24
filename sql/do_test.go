package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	d := NewDo().Do(Raw("SLEEP ?", 5))
	d.Do(Raw("CONTRACT(1, 2)"))

	assert.Equal(t, "DO SLEEP ?, CONTRACT(1, 2)", d.Prepare())
	assert.Equal(t, []any{5}, d.Binds())
}
