package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql"
	"github.com/stretchr/testify/assert"
)

func TestDbInsert(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	dt := NewData()
	dt.Set("name", "kovey")
	dt.Set("age", 18)
	n, err := Insert(context.Background(), "user", dt)
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestDbInsertErr(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	insErr := errors.New("insert error")
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewErrorResult(insErr))
	dt := NewData()
	dt.Set("name", "kovey")
	dt.Set("age", 18)
	n, err := Insert(context.Background(), "user", dt)
	assert.Equal(t, int64(0), n)
	assert.Equal(t, insErr, err.(*SqlErr).Err)
}

func TestDbUpdate(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	now := time.Now().Unix()
	mock.ExpectPrepare("UPDATE `user_ext` SET `last_time` = ? WHERE `id` = ?").ExpectExec().WithArgs(now, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	dt := NewData()
	dt.Set("last_time", now)
	where := sql.NewWhere()
	where.Where("id", ksql.Eq, 1)
	n, err := Update(context.Background(), "user_ext", dt, where)
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestDbUpdateErr(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	upErr := errors.New("update error")
	now := time.Now().Unix()
	mock.ExpectPrepare("UPDATE `user_ext` SET `last_time` = ? WHERE `id` = ?").ExpectExec().WithArgs(now, 1).WillReturnResult(sqlmock.NewErrorResult(upErr))
	dt := NewData()
	dt.Set("last_time", now)
	where := sql.NewWhere()
	where.Where("id", ksql.Eq, 1)
	n, err := Update(context.Background(), "user_ext", dt, where)
	assert.Equal(t, int64(0), n)
	assert.Equal(t, upErr, err.(*SqlErr).Err)
}

func TestDbDelete(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	mock.ExpectPrepare("DELETE FROM `user_ext` WHERE `id` = ?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	where := sql.NewWhere()
	where.Where("id", ksql.Eq, 1)
	n, err := Delete(context.Background(), "user_ext", where)
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestDbDeleteErr(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	delErr := errors.New("del error")
	mock.ExpectPrepare("DELETE FROM `user_ext` WHERE `id` = ?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(delErr))
	where := sql.NewWhere()
	where.Where("id", ksql.Eq, 1)
	n, err := Delete(context.Background(), "user_ext", where)
	assert.Equal(t, int64(0), n)
	assert.Equal(t, delErr, err.(*SqlErr).Err)
}

func TestDbInsertFrom(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) SELECT `name`, `age` FROM `users` WHERE `id` > ?").ExpectExec().WithArgs(100).WillReturnResult(sqlmock.NewResult(1, 1))
	query := NewQuery().Table("users").Columns("name", "age").Where("id", ">", 100)
	n, err := InsertFrom(context.Background(), "user", []string{"name", "age"}, query)
	assert.Equal(t, int64(1), n)
	assert.Nil(t, err)
}

func TestDbInsertFromErr(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn
	delErr := errors.New("del error")
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) SELECT `name`, `age` FROM `users` WHERE `id` > ?").ExpectExec().WithArgs(100).WillReturnResult(sqlmock.NewErrorResult(delErr))
	query := NewQuery().Table("users").Columns("name", "age").Where("id", ">", 100)
	n, err := InsertFrom(context.Background(), "user", []string{"name", "age"}, query)
	assert.Equal(t, int64(0), n)
	assert.Equal(t, delErr, err.(*SqlErr).Err)
}

func TestDbQuery(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user
	query := NewQuery()
	query.Table("user").Columns(columns...).Order("id").OrderDesc("age").Where("id", ksql.Ge, 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` >= ? ORDER BY `id` ASC, `age` DESC").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23).AddRow(2, 15, "test", now, 34.13))
	err = Query(context.Background(), query, &rows)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rows))
	tu := rows[0]
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	tu = rows[1]
	assert.Equal(t, int64(2), tu.Id)
	assert.Equal(t, 15, tu.Age)
	assert.Equal(t, "test", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(34.13), tu.Balance)
}

func TestDbQueryRow(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	query := NewQuery()
	query.Table("user").Columns(columns...).Where("id", ksql.Eq, 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = QueryRow(context.Background(), query, tu)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestDbFind(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = Find(context.Background(), tu, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestDbLock(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ? FOR UPDATE").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("UPDATE `user` SET `balance` = `balance` + ? WHERE `id` = ?").ExpectExec().WithArgs(100, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	up := NewUpdate()
	up.Table("user").IncColumn("balance", 100).Where(NewWhere().Where("id", ksql.Eq, 1))
	in := NewInsert()
	in.Table("user").Add("name", "kovey").Add("age", 18)
	err = Transaction(context.Background(), func(ctx context.Context, db ksql.ConnectionInterface) error {
		tu := newTestUser()
		if err := Lock(ctx, db, tu, 1); err != nil {
			return err
		}

		if _, err := db.Insert(ctx, in); err != nil {
			return err
		}

		_, err := db.Update(ctx, up)
		return err
	})

	if err != nil {
		t.Fatal(err.(*TxErr).callErr)
	}

	assert.Nil(t, err)
}

func TestDbTable(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("ALTER TABLE `user` DROP COLUMN `user_name`, ADD COLUMN `level` BIGINT(20) DEFAULT '0' KEY COMMENT '等级', ADD FOREIGN KEY `idx_user_id` (`user_id`), DROP INDEX `idx_level`").ExpectExec().WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	err = Table(context.Background(), "user", func(table ksql.TableInterface) {
		table.DropColumn("user_name")
		table.AddBigInt("level").Index().Default("0").Comment("等级")
		table.AddIndex("idx_user_id").Foreign().Algorithm(ksql.Index_Alg_BTree).Comment("索引").Columns("user_id")
		table.DropIndex("idx_level")
	})
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDbCreate(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("CREATE TABLE `user` (`level` BIGINT(20) DEFAULT '0' KEY COMMENT '等级', FOREIGN KEY `idx_user_id` (`user_id`))").ExpectExec().WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	err = Create(context.Background(), "user", func(table ksql.TableInterface) {
		table.DropColumn("user_name")
		table.AddBigInt("level").Index().Default("0").Comment("等级")
		table.AddIndex("idx_user_id").Foreign().Algorithm(ksql.Index_Alg_BTree).Comment("索引").Columns("user_id")
		table.DropIndex("idx_level")
	})
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDbDropTable(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("DROP TABLE `user`").ExpectExec().WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	err = DropTable(context.Background(), "user")
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDbDropTableIfExists(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("DROP TABLE IF EXISTS `user`").ExpectExec().WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	err = DropTableIfExists(context.Background(), "user")
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDbShowDLL(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("SHOW CREATE TABLE `user`").ExpectQuery().WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{"table", "ddl"}).AddRow("user", "ddl"))
	ddl, err := ShowDDL(context.Background(), "user")
	assert.Nil(t, err)
	assert.Equal(t, "ddl", ddl)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDbScan(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	query := NewQuery().Table("user").Columns("id", "name").Where("id", "=", 1)
	mock.ExpectPrepare("SELECT `id`, `name` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "kovey"))
	var id int
	var name string
	err = Scan(context.Background(), query, &id, &name)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, id)
	assert.Equal(t, "kovey", name)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDbFindBy(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = FindBy(context.Background(), tu, func(query ksql.QueryInterface) {
		query.Where("id", ksql.Eq, 1)
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestDbLockShare(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ? FOR SHARE").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("UPDATE `user` SET `balance` = `balance` + ? WHERE `id` = ?").ExpectExec().WithArgs(100, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	up := NewUpdate()
	up.Table("user").IncColumn("balance", 100).Where(NewWhere().Where("id", ksql.Eq, 1))
	in := NewInsert()
	in.Table("user").Add("name", "kovey").Add("age", 18)
	err = Transaction(context.Background(), func(ctx context.Context, db ksql.ConnectionInterface) error {
		tu := newTestUser()
		if err := LockShare(ctx, db, tu, 1); err != nil {
			return err
		}

		if _, err := db.Insert(ctx, in); err != nil {
			return err
		}

		_, err := db.Update(ctx, up)
		return err
	})

	if err != nil {
		t.Fatal(err.(*TxErr).callErr)
	}

	assert.Nil(t, err)
}

func TestDbLockShareBy(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ? FOR SHARE").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("UPDATE `user` SET `balance` = `balance` + ? WHERE `id` = ?").ExpectExec().WithArgs(100, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	up := NewUpdate()
	up.Table("user").IncColumn("balance", 100).Where(NewWhere().Where("id", ksql.Eq, 1))
	in := NewInsert()
	in.Table("user").Add("name", "kovey").Add("age", 18)
	err = Transaction(context.Background(), func(ctx context.Context, db ksql.ConnectionInterface) error {
		tu := newTestUser()
		if err := LockByShare(ctx, db, tu, func(query ksql.QueryInterface) {
			query.Where("id", ksql.Eq, 1)
		}); err != nil {
			return err
		}

		if _, err := db.Insert(ctx, in); err != nil {
			return err
		}

		_, err := db.Update(ctx, up)
		return err
	})

	if err != nil {
		t.Fatal(err.(*TxErr).callErr)
	}

	assert.Nil(t, err)
}

func TestDbLockBy(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ? FOR UPDATE").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	mock.ExpectPrepare("INSERT INTO `user` (`name`, `age`) VALUES (?, ?)").ExpectExec().WithArgs("kovey", 18).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("UPDATE `user` SET `balance` = `balance` + ? WHERE `id` = ?").ExpectExec().WithArgs(100, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	up := NewUpdate()
	up.Table("user").IncColumn("balance", 100).Where(NewWhere().Where("id", ksql.Eq, 1))
	in := NewInsert()
	in.Table("user").Add("name", "kovey").Add("age", 18)
	err = Transaction(context.Background(), func(ctx context.Context, db ksql.ConnectionInterface) error {
		tu := newTestUser()
		if err := LockBy(ctx, db, tu, func(query ksql.QueryInterface) {
			query.Where("id", ksql.Eq, 1)
		}); err != nil {
			return err
		}

		if _, err := db.Insert(ctx, in); err != nil {
			return err
		}

		_, err := db.Update(ctx, up)
		return err
	})

	if err != nil {
		t.Fatal(err.(*TxErr).callErr)
	}

	assert.Nil(t, err)
}
