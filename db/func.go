package db

import (
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql"
)

type NewQueryFun func() ksql.QueryInterface
type RawFun func(raw string, binds ...any) ksql.ExpressInterface
type NewInsertFun func() ksql.InsertInterface
type NewUpdateFun func() ksql.UpdateInterface
type NewDeleteFun func() ksql.DeleteInterface
type NewSchemaFun func() ksql.SchemaInterface
type NewDropTableFun func() ksql.DropTableInterface
type NewCreateTableFun func() ksql.CreateTableInterface
type NewAlterTableFun func() ksql.AlterInterface
type NewTableFun func() ksql.TableInterface
type NewWhereFun func() ksql.WhereInterface

var NewWhere NewWhereFun = func() ksql.WhereInterface {
	return sql.NewWhere()
}
var NewQuery NewQueryFun = func() ksql.QueryInterface {
	return sql.NewQuery()
}
var Raw RawFun = sql.Raw
var NewInsert NewInsertFun = func() ksql.InsertInterface {
	return sql.NewInsert()
}
var NewUpdate NewUpdateFun = func() ksql.UpdateInterface {
	return sql.NewUpdate()
}
var NewDelete NewDeleteFun = func() ksql.DeleteInterface {
	return sql.NewDelete()
}
var NewSchema NewSchemaFun = func() ksql.SchemaInterface {
	return sql.NewSchema()
}
var NewDropTable NewDropTableFun = func() ksql.DropTableInterface {
	return sql.NewDropTable()
}
var NewCreateTable NewCreateTableFun = func() ksql.CreateTableInterface {
	return sql.NewTable()
}
var NewAlterTable NewAlterTableFun = func() ksql.AlterInterface {
	return sql.NewAlter()
}
var NewTable NewTableFun = func() ksql.TableInterface {
	return NewTableBuilder()
}
