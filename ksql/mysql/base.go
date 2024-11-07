package mysql

import "github.com/kovey/db-go/v3"

type base struct {
	conn  ksql.ConnectionInterface
	empty bool
}

func (b *base) SetConn(conn ksql.ConnectionInterface) {
	b.conn = conn
}

func (b *base) FromFetch() {
	b.empty = false
}

func (b *base) Conn() ksql.ConnectionInterface {
	return b.conn
}
