package db

import ksql "github.com/kovey/db-go/v3"

type Row struct{}

func (Row) Values() []any                    { return nil }
func (Row) Clone() ksql.RowInterface         { return nil }
func (Row) SetConn(ksql.ConnectionInterface) {}
func (Row) FromFetch()                       {}
func (Row) Conn() ksql.ConnectionInterface   { return nil }
