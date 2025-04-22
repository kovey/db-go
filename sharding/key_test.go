package sharding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type test_key struct {
	name string
}

func (t *test_key) String() string {
	return t.name
}

func TestKey(t *testing.T) {
	total := 10
	assert.Equal(t, 1, node(1, total))
	assert.Equal(t, 1, node(int8(1), total))
	assert.Equal(t, 1, node(int16(1), total))
	assert.Equal(t, 1, node(int32(1), total))
	assert.Equal(t, 1, node(int64(1), total))
	assert.Equal(t, 1, node(uint(1), total))
	assert.Equal(t, 1, node(uint8(1), total))
	assert.Equal(t, 1, node(uint16(1), total))
	assert.Equal(t, 1, node(uint32(1), total))
	assert.Equal(t, 1, node(uint64(1), total))
	assert.Equal(t, 1, node("ABCDEFGHIJKL!@##@$@#$@#^&*", total))
	assert.Equal(t, 8, node(&test_key{name: "adlsjdflsflsdjflsf80980380!@#@!#!"}, total))
	assert.Equal(t, 0, node(newTestModel(), total))
}
