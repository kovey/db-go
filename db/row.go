package db

import ksql "github.com/kovey/db-go/v3"

type Row struct{}

func (Row) Values() []any                                        { return nil }
func (Row) Clone() ksql.RowInterface                             { return nil }
func (Row) WithConn(ksql.ConnectionInterface)                    {}
func (Row) Scan(s ksql.ScanInterface, r ksql.RowInterface) error { return s.Scan(r.Values()...) }
func (Row) Conn() ksql.ConnectionInterface                       { return nil }
func (Row) Sharding(ksql.Sharding)                               {}
