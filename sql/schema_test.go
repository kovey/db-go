package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	sc := NewSchema().Create().Schema("user").Charset("utf8mb4").Collate("utf8mb4_general_ci")
	assert.Equal(t, "CREATE SCHEMA `utf8mb4` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", sc.Prepare())
	assert.Nil(t, sc.Binds())
	assert.Equal(t, "CREATE SCHEMA `utf8mb4` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", sc.Prepare())
}
