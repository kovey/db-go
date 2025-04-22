package model

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/stretchr/testify/assert"
)

func TestTableManagerDay(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	err = db.InitBy(testDb, "mysql")
	assert.Nil(t, err)
	now := time.Now()
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, -3).Format(ksql.Day_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}).AddRow("user_"))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, -2).Format(ksql.Day_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, 0, -2).Format(ksql.Day_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, -1).Format(ksql.Day_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, 0, -1).Format(ksql.Day_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, 0).Format(ksql.Day_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, 0, 0).Format(ksql.Day_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, 1).Format(ksql.Day_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, 0, 1).Format(ksql.Day_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, 2).Format(ksql.Day_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, 0, 2).Format(ksql.Day_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, 3).Format(ksql.Day_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, 0, 3).Format(ksql.Day_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	tm := NewTableManager(nil)
	tm.Append(&Template{Table: "user", Keep: 86400 * 7, Type: ksql.Sharding_Day, TemplateTable: "user_tpl"})
	tm.Create(context.Background())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestTableManagerMonth(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	err = db.InitBy(testDb, "mysql")
	assert.Nil(t, err)
	now := time.Now()
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, -3, 0).Format(ksql.Month_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}).AddRow("user_"))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, -2, 0).Format(ksql.Month_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, -2, 0).Format(ksql.Month_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, -1, 0).Format(ksql.Month_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, -1, 0).Format(ksql.Month_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, 0).Format(ksql.Month_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, 0, 0).Format(ksql.Month_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 1, 0).Format(ksql.Month_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, 1, 0).Format(ksql.Month_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 2, 0).Format(ksql.Month_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, 2, 0).Format(ksql.Month_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 3, 0).Format(ksql.Month_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	mock.ExpectPrepare(fmt.Sprintf("CREATE TABLE `user_%s` (LIKE `user_tpl`)", now.AddDate(0, 3, 0).Format(ksql.Month_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	tm := NewTableManager(nil)
	tm.Append(&Template{Table: "user", Keep: 86400 * 7, Type: ksql.Sharding_Month, TemplateTable: "user_tpl"})
	tm.Create(context.Background())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestTableManagerDayDelete(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	err = db.InitBy(testDb, "mysql")
	assert.Nil(t, err)
	now := time.Now()
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, -1).Format(ksql.Day_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}).AddRow("user_"))
	mock.ExpectPrepare(fmt.Sprintf("DROP TABLE `user_%s`", now.AddDate(0, 0, -1).Format(ksql.Day_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, -2).Format(ksql.Day_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}).AddRow("user_"))
	mock.ExpectPrepare(fmt.Sprintf("DROP TABLE `user_%s`", now.AddDate(0, 0, -2).Format(ksql.Day_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, 0, -3).Format(ksql.Day_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	tm := NewTableManager(nil)
	tm.Append(&Template{Table: "user", Keep: 1, Type: ksql.Sharding_Day, TemplateTable: "user_tpl"})
	tm.Delete(context.Background())
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestTableManagerMonthDelete(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	err = db.InitBy(testDb, "mysql")
	assert.Nil(t, err)
	now := time.Now()
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, -1, 0).Format(ksql.Month_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}).AddRow("user_"))
	mock.ExpectPrepare(fmt.Sprintf("DROP TABLE `user_%s`", now.AddDate(0, -1, 0).Format(ksql.Month_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, -2, 0).Format(ksql.Month_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}).AddRow("user_"))
	mock.ExpectPrepare(fmt.Sprintf("DROP TABLE `user_%s`", now.AddDate(0, -2, 0).Format(ksql.Month_Format))).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare(fmt.Sprintf("SHOW TABLES LIKE 'user_%s'", now.AddDate(0, -3, 0).Format(ksql.Month_Format))).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"table_name"}))
	tm := NewTableManager(nil)
	tm.Append(&Template{Table: "user", Keep: 1, Type: ksql.Sharding_Month, TemplateTable: "user_tpl"})
	tm.Delete(context.Background())
	assert.Nil(t, mock.ExpectationsWereMet())
}
