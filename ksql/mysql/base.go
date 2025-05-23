package mysql

import (
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

type base struct {
	*db.Row
	conn          ksql.ConnectionInterface
	empty         bool
	isInitialized bool
}

func (b *base) WithConn(conn ksql.ConnectionInterface) {
	b.conn = conn
}

func (b *base) Conn() ksql.ConnectionInterface {
	return b.conn
}

func (b *base) Scan(s ksql.ScanInterface, r ksql.RowInterface) error {
	if err := s.Scan(r.Values()...); err != nil {
		return err
	}

	b.isInitialized = true
	b.empty = true
	return nil
}
