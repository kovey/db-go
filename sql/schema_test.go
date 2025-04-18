package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	sc := NewSchema().Create().Alter().Schema("user").Charset("utf8mb4").Collate("utf8mb4_general_ci")
	assert.Equal(t, "CREATE SCHEMA `user` DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci", sc.Prepare())
	assert.Nil(t, sc.Binds())
	assert.Equal(t, "CREATE SCHEMA `user` DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci", sc.Prepare())
}

func TestSchemaEncryption(t *testing.T) {
	sc := NewSchema().Create().Alter().Schema("user").Encryption("Y").ReadOnly("1").IfNotExists()
	assert.Equal(t, "CREATE SCHEMA IF NOT EXISTS `user` DEFAULT ENCRYPTION = 'Y'", sc.Prepare())
	assert.Nil(t, sc.Binds())
	assert.Equal(t, "CREATE SCHEMA IF NOT EXISTS `user` DEFAULT ENCRYPTION = 'Y'", sc.Prepare())
}

func TestSchemaAlter(t *testing.T) {
	sc := NewSchema().Alter().Create().Schema("user").Charset("utf8mb4").Collate("utf8mb4_general_ci").IfNotExists()
	assert.Equal(t, "ALTER SCHEMA `user` DEFAULT CHARACTER SET = utf8mb4 DEFAULT COLLATE = utf8mb4_general_ci", sc.Prepare())
	assert.Nil(t, sc.Binds())
	assert.Equal(t, "ALTER SCHEMA `user` DEFAULT CHARACTER SET = utf8mb4 DEFAULT COLLATE = utf8mb4_general_ci", sc.Prepare())
}

func TestSchemaAlterEncryption(t *testing.T) {
	sc := NewSchema().Alter().Create().Schema("user").Encryption("Y").ReadOnly("1")
	assert.Equal(t, "ALTER SCHEMA `user` DEFAULT ENCRYPTION = 'Y' READ ONLY = 1", sc.Prepare())
	assert.Nil(t, sc.Binds())
	assert.Equal(t, "ALTER SCHEMA `user` DEFAULT ENCRYPTION = 'Y' READ ONLY = 1", sc.Prepare())
}
