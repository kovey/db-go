package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	s := NewServer().Server("serv_kovey")
	s.WrapperName("mysql").Option(ksql.Serv_Opt_Key_Host, "127.0.0.1").Option(ksql.Serv_Opt_Key_Database, "kovey")
	assert.Equal(t, "CREATE SERVER `serv_kovey` FOREIGN DATA WRAPPER mysql OPTIONS (HOST '127.0.0.1', DATABASE 'kovey')", s.Prepare())
	assert.Nil(t, s.Binds())
}

func TestServerAlter(t *testing.T) {
	s := NewServer().Server("serv_kovey").Alter()
	s.WrapperName("mysql").Option(ksql.Serv_Opt_Key_Host, "127.0.0.1").Option(ksql.Serv_Opt_Key_Database, "kovey")
	assert.Equal(t, "ALTER SERVER `serv_kovey` OPTIONS (HOST '127.0.0.1', DATABASE 'kovey')", s.Prepare())
	assert.Nil(t, s.Binds())
}
