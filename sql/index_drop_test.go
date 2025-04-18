package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestIndexDrop(t *testing.T) {
	s := NewIndexDrop().Index("user").Algorithm(ksql.Index_Alg_Option_Copy).Lock(ksql.Index_Lock_Option_Default).Table("users")
	assert.Equal(t, "DROP INDEX `user` ON `users` ALGORITHM = COPY LOCK = DEFAULT", s.Prepare())
	assert.Nil(t, s.Binds())
}
