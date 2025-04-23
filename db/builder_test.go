package db

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

type test_user struct {
	*Row
	Id         int64
	Age        int
	Name       string
	CreateTime time.Time
	Balance    float64
}

func newTestUser() *test_user {
	return &test_user{Row: &Row{}}
}

func (t *test_user) Clone() ksql.RowInterface {
	return newTestUser()
}

func (t *test_user) Columns() []string {
	return []string{"id", "age", "name", "create_time", "balance"}
}

func (t *test_user) Values() []any {
	return []any{&t.Id, &t.Age, &t.Name, &t.CreateTime, &t.Balance}
}

func (t *test_user) Delete(ctx context.Context) error {
	return nil
}

func (t *test_user) Empty() bool {
	return t.Id == 0
}

func (t *test_user) Table() string {
	return "user"
}
func (t *test_user) PrimaryId() string {
	return "id"
}
func (t *test_user) Save(ctx context.Context) error {
	return nil
}

func (t *test_user) OnUpdateBefore(conn ksql.ConnectionInterface) error {
	return nil
}

func (t *test_user) OnUpdateAfter(conn ksql.ConnectionInterface) error {
	return nil
}

func (t *test_user) OnCreateBefore(conn ksql.ConnectionInterface) error {
	return nil
}

func (t *test_user) OnCreateAfter(conn ksql.ConnectionInterface) error {
	return nil
}

func (t *test_user) OnDeleteBefore(conn ksql.ConnectionInterface) error {
	return nil
}

func (t *test_user) OnDeleteAfter(conn ksql.ConnectionInterface) error {
	return nil
}

type test_user_count struct {
	*Row
	Id         int64
	Age        int
	Name       string
	CreateTime time.Time
	Balance    float64
	Count      int
}

func newTestUserCount() *test_user_count {
	return &test_user_count{Row: &Row{}}
}

func (t *test_user_count) Clone() ksql.RowInterface {
	return newTestUserCount()
}

func (t *test_user_count) Values() []any {
	return []any{&t.Id, &t.Age, &t.Name, &t.CreateTime, &t.Balance, &t.Count}
}

func TestBuilderFirst(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	builder := NewBuilder(tu)
	builder.Table("user").Columns(columns...).Where("id", ksql.Eq, 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = builder.First(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestBuilderAll(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user
	builder := Rows(&rows).(*Builder[*test_user])
	builder.Table("user").Columns(columns...).Order("id").OrderDesc("age").Where("id", ksql.Ge, 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` >= ? ORDER BY `id` ASC, `age` DESC").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23).AddRow(2, 15, "test", now, 34.13))
	err = builder.All(context.Background())
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

func TestBuilderWithSubSelect(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	sub := NewQuery().Table("users").Columns("id").Func("COUNT", "1", "count").Group("id").Between("create_time", 100, 200)
	columns := []string{"u.id", "us.age", "us.name", "us.create_time", "us.balance", "u.count"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user_count
	builder := Rows(&rows).(*Builder[*test_user_count])
	builder.TableBy(sub, "u").Columns(columns...).Where("us.id", ksql.Ge, 1).LeftJoin("user").As("us").On("us.id", "=", "u.id")
	mock.ExpectPrepare("SELECT `u`.`id`, `us`.`age`, `us`.`name`, `us`.`create_time`, `us`.`balance`, `u`.`count` FROM (SELECT `id`, COUNT(`1`) AS `count` FROM `users` WHERE `create_time` BETWEEN ? AND ? GROUP BY `id`) LEFT JOIN `user` AS `us` ON (`us`.`id` = `u`.`id`) WHERE `us`.`id` >= ?").ExpectQuery().WithArgs(100, 200, 1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23, 300).AddRow(2, 15, "test", now, 34.13, 100))
	err = builder.All(context.Background())
	if err != nil {
		t.Fatal(err, builder.query.Binds())
	}
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rows))
	tu := rows[0]
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	assert.Equal(t, 300, tu.Count)
	tu = rows[1]
	assert.Equal(t, int64(2), tu.Id)
	assert.Equal(t, 15, tu.Age)
	assert.Equal(t, "test", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(34.13), tu.Balance)
	assert.Equal(t, 100, tu.Count)
	tu = rows[1]
}

func TestBuilderColumnFunc(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user_count
	builder := Rows(&rows).(*Builder[*test_user_count])
	builder.Table("user").As("u").Columns(columns...).Column("create_time", "ct").Func("SUM", "balance", "balance").ColumnsExpress(Raw("COUNT(1) as count"))
	builder.Where("id", ksql.Ge, 1).WhereExpress(Raw("name LIKE ?", "%kovey%")).OrWhere(func(wi ksql.WhereInterface) {
		wi.Between("create_time", "2025-04-03 01:11:11", "2025-04-03 11:15:11")
	})
	builder.WhereIsNull("update_time").WhereIsNotNull("age").WhereIn("id", []any{1, 3, 5}).WhereNotIn("id", []any{2, 4, 6})
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time` AS `ct`, SUM(`balance`) AS `balance`, (COUNT(1) as count) FROM `user` AS `u` WHERE `id` >= ? AND name LIKE ? AND `update_time` IS NULL AND `age` IS NULL AND `id` IN (?, ?, ?) AND `id` NOT IN (?, ?, ?) OR (`create_time` BETWEEN ? AND ?)").ExpectQuery().WithArgs(1, "%kovey%", 1, 3, 5, 2, 4, 6, "2025-04-03 01:11:11", "2025-04-03 11:15:11").WillReturnRows(sqlmock.NewRows([]string{"id", "age", "name", "ct", "balance", "count"}).AddRow(1, 18, "kovey", now, 30.23, 300).AddRow(2, 15, "test", now, 34.13, 100))
	err = builder.All(context.Background())
	if err != nil {
		t.Fatal(err, builder.query.Binds())
	}
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rows))
	tu := rows[0]
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	assert.Equal(t, 300, tu.Count)
	tu = rows[1]
	assert.Equal(t, int64(2), tu.Id)
	assert.Equal(t, 15, tu.Age)
	assert.Equal(t, "test", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(34.13), tu.Balance)
	assert.Equal(t, 100, tu.Count)
	tu = rows[1]
}

func TestBuilderWhereSubSelect(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	sub := NewQuery().Table("users").Columns("id").Between("create_time", 100, 200)
	columns := []string{"u.id", "us.age", "us.name", "us.create_time", "us.balance", "u.count"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user_count
	builder := Rows(&rows).(*Builder[*test_user_count])
	builder.Table("user_ext").As("us").Columns(columns...).Where("us.id", ksql.Ge, 1).LeftJoin("user").As("us").On("us.id", "=", "u.id")
	builder.WhereInBy("us.id", sub).WhereNotInBy("u.id", sub)
	builder.AndWhere(func(w ksql.WhereInterface) {
		w.Where("us.age", ">=", 1).Where("us.age", "<=", 18)
	})
	mock.ExpectPrepare("SELECT `u`.`id`, `us`.`age`, `us`.`name`, `us`.`create_time`, `us`.`balance`, `u`.`count` FROM `user_ext` AS `us` LEFT JOIN `user` AS `us` ON (`us`.`id` = `u`.`id`) WHERE `us`.`id` >= ? AND `us`.`id` IN (SELECT `id` FROM `users` WHERE `create_time` BETWEEN ? AND ?) AND `u`.`id` NOT IN (SELECT `id` FROM `users` WHERE `create_time` BETWEEN ? AND ?) AND (`us`.`age` >= ? AND `us`.`age` <= ?)").ExpectQuery().WithArgs(1, 100, 200, 100, 200, 1, 18).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23, 300).AddRow(2, 15, "test", now, 34.13, 100))
	err = builder.All(context.Background())
	if err != nil {
		t.Fatal(err, builder.query.Binds())
	}
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rows))
	tu := rows[0]
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	assert.Equal(t, 300, tu.Count)
	tu = rows[1]
	assert.Equal(t, int64(2), tu.Id)
	assert.Equal(t, 15, tu.Age)
	assert.Equal(t, "test", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(34.13), tu.Balance)
	assert.Equal(t, 100, tu.Count)
	tu = rows[1]
}

func TestBuilderHavingSubSelect(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	sub := NewQuery().Table("users").Columns("id").Between("create_time", 100, 200)
	columns := []string{"u.id", "us.age", "us.name", "us.create_time", "us.balance", "u.count"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user_count
	builder := Rows(&rows).(*Builder[*test_user_count])
	builder.Table("user_ext").As("us").Columns(columns...).Between("u.id", 1, 1000).NotBetween("u.id", 1001, 2000).Where("us.id", ksql.Ge, 1).LeftJoin("user").As("us").On("us.id", "=", "u.id")
	builder.Having("name", "LIKE", "%kovey%")
	builder.HavingInBy("us.id", sub).HavingNotInBy("u.id", sub)
	builder.AndHaving(func(w ksql.HavingInterface) {
		w.Having("us.age", ">=", 1).Having("us.age", "<=", 18)
	})
	builder.OrHaving(func(w ksql.HavingInterface) {
		w.Having("u.count", ">=", 1).Having("u.count", "<=", 18)
	})
	mock.ExpectPrepare("SELECT `u`.`id`, `us`.`age`, `us`.`name`, `us`.`create_time`, `us`.`balance`, `u`.`count` FROM `user_ext` AS `us` LEFT JOIN `user` AS `us` ON (`us`.`id` = `u`.`id`) WHERE `u`.`id` BETWEEN ? AND ? AND `u`.`id` NOT BETWEEN ? AND ? AND `us`.`id` >= ? HAVING `name` LIKE ? AND `us`.`id` IN (SELECT `id` FROM `users` WHERE `create_time` BETWEEN ? AND ?) AND `u`.`id` NOT IN (SELECT `id` FROM `users` WHERE `create_time` BETWEEN ? AND ?) AND (`us`.`age` >= ? AND `us`.`age` <= ?) OR (`u`.`count` >= ? AND `u`.`count` <= ?)").ExpectQuery().WithArgs(1, 1000, 1001, 2000, 1, "%kovey%", 100, 200, 100, 200, 1, 18, 1, 18).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23, 300).AddRow(2, 15, "test", now, 34.13, 100))
	err = builder.All(context.Background())
	if err != nil {
		t.Fatal(err, builder.query.Binds())
	}
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rows))
	tu := rows[0]
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	assert.Equal(t, 300, tu.Count)
	tu = rows[1]
	assert.Equal(t, int64(2), tu.Id)
	assert.Equal(t, 15, tu.Age)
	assert.Equal(t, "test", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(34.13), tu.Balance)
	assert.Equal(t, 100, tu.Count)
	tu = rows[1]
}

func TestBuilderHaving(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	builder := NewBuilder(tu)
	builder.Table("user").Columns(columns...).Where("id", ksql.Eq, 1).HavingExpress(Raw("name LIKE ?", "%kovey%")).HavingIsNull("age").HavingIsNotNull("name").HavingIn("balance", []any{1, 3, 5})
	builder.HavingNotIn("count", []any{1, 3, 5}).HavingBetween("cc", 100, 200).HavingNotBetween("aa", 200, 300)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ? HAVING name LIKE ? AND `age` IS NULL AND `name` IS NULL AND `balance` IN (?, ?, ?) AND `count` NOT IN (?, ?, ?) AND `cc` BETWEEN ? AND ? AND `aa` NOT BETWEEN ? AND ?").ExpectQuery().WithArgs(1, "%kovey%", 1, 3, 5, 1, 3, 5, 100, 200, 200, 300).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = builder.First(context.Background())
	if err != nil {
		t.Fatal(err, builder.query.Binds())
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestBuilderPage(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	builder := Rows(&[]*test_user{}).(*Builder[*test_user])
	builder.Table("user").Columns(columns...).Where("id", ksql.Eq, 1).Distinct()
	mock.ExpectPrepare("SELECT DISTINCT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ? LIMIT ? OFFSET ?").ExpectQuery().WithArgs(1, 2, 0).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23).AddRow(2, 15, "test", now, 34.13))
	mock.ExpectPrepare("SELECT DISTINCT (COUNT(1) as count) FROM `user` WHERE `id` = ? LIMIT ? OFFSET ?").ExpectQuery().WithArgs(1, 1, 0).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(9))
	pageInfo, err := builder.Pagination(context.Background(), 1, 2)
	if err != nil {
		t.Fatal(err, builder.query.Binds())
	}
	tu := pageInfo.List()[0]
	assert.Equal(t, 2, len(pageInfo.List()))
	assert.Equal(t, uint64(9), pageInfo.TotalCount())
	assert.Equal(t, uint64(5), pageInfo.TotalPage())
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	tu = pageInfo.List()[1]
	assert.Equal(t, int64(2), tu.Id)
	assert.Equal(t, 15, tu.Age)
	assert.Equal(t, "test", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(34.13), tu.Balance)
}

func TestBuilderAllGroup(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user_count
	builder := Rows(&rows).(*Builder[*test_user_count])
	builder.Table("user").Columns(columns...).FuncDistinct("COUNT", "id", "count").Where("id", ksql.Ge, 1).Group("id", "age", "name", "create_time", "balance")
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance`, DISTINCT `id` AS `count` FROM `user` WHERE `id` >= ? GROUP BY `id`, `age`, `name`, `create_time`, `balance`").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "age", "name", "create_time", "balance", "count"}).AddRow(1, 18, "kovey", now, 30.23, 12).AddRow(2, 15, "test", now, 34.13, 13))
	err = builder.All(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rows))
	tu := rows[0]
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	assert.Equal(t, 12, tu.Count)
	tu = rows[1]
	assert.Equal(t, int64(2), tu.Id)
	assert.Equal(t, 15, tu.Age)
	assert.Equal(t, "test", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(34.13), tu.Balance)
	assert.Equal(t, 13, tu.Count)
}

func TestBuilderAllWithJoin(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"u.id", "u.age", "ue.name", "us.create_time", "b.balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user_count
	builder := Rows(&rows).(*Builder[*test_user_count])
	builder.Table("user").Columns(columns...).As("u")
	builder.Join("user_ext").As("ue").On("ue.id", "=", "u.id")
	builder.RightJoin("balance").As("b").On("b.id", "=", "u.id")
	builder.JoinExpress(Raw("LEFT JOIN user_sys as us ON us.id = u.id"))
	builder.Where("u.id", ">=", 1)
	mock.ExpectPrepare("SELECT `u`.`id`, `u`.`age`, `ue`.`name`, `us`.`create_time`, `b`.`balance` FROM `user` AS `u` INNER JOIN `user_ext` AS `ue` ON (`ue`.`id` = `u`.`id`) RIGHT JOIN `balance` AS `b` ON (`b`.`id` = `u`.`id`) LEFT JOIN user_sys as us ON us.id = u.id WHERE `u`.`id` >= ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "age", "name", "create_time", "balance", "count"}).AddRow(1, 18, "kovey", now, 30.23, 12).AddRow(2, 15, "test", now, 34.13, 13))
	err = builder.All(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(rows))
	tu := rows[0]
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	assert.Equal(t, 12, tu.Count)
	tu = rows[1]
	assert.Equal(t, int64(2), tu.Id)
	assert.Equal(t, 15, tu.Age)
	assert.Equal(t, "test", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(34.13), tu.Balance)
	assert.Equal(t, 13, tu.Count)
}

func TestBuilderExists(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	builder := NewBuilder(newTestUser()).WithConn(conn)
	builder.Table("user").Columns(columns...).Where("id", ksql.Eq, 1).ForUpdate()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ? LIMIT ? FOR UPDATE").ExpectQuery().WithArgs(1, 1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	ok, err := builder.Exist(context.Background())
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestBuilderFirstWithConn(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUser()
	builder := NewBuilder(tu).WithConn(conn)
	builder.Table("user").Columns(columns...).Where("id", ksql.Eq, 1).For().Share()
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance` FROM `user` WHERE `id` = ? FOR SHARE").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = builder.First(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestBuilderMax(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUserCount()
	builder := Build(tu).WithConn(conn)
	builder.Table("user").Columns(columns...).Where("id", ksql.Eq, 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance`, MAX(`count`) AS `count` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "age", "name", "create_time", "balance", "count"}).AddRow(1, 18, "kovey", now, 30.23, 12))
	err = builder.Max(context.Background(), "count")
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	assert.Equal(t, 12, tu.Count)
}

func TestBuilderMin(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	tu := newTestUserCount()
	builder := NewBuilder(tu).WithConn(conn)
	builder.Table("user").Columns(columns...).Where("id", ksql.Eq, 1)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `balance`, MIN(`count`) AS `count` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "age", "name", "create_time", "balance", "count"}).AddRow(1, 18, "kovey", now, 30.23, 12))
	err = builder.Min(context.Background(), "count")
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
	assert.Equal(t, 12, tu.Count)
}

func TestBuilderSumFloat(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)

	builder := NewBuilder(newTestUserCount()).WithConn(conn)
	builder.Table("user").Where("id", ksql.Eq, 1)
	mock.ExpectPrepare("SELECT SUM(`balance`) AS `balance` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(30.23))
	ba, err := builder.SumFloat(context.Background(), "balance")
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 30.23, ba)
}

func TestBuilderSum(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)

	builder := NewBuilder(newTestUserCount()).WithConn(conn)
	builder.Table("user").Where("id", ksql.Eq, 1)
	mock.ExpectPrepare("SELECT SUM(`balance`) AS `balance` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(30))
	ba, err := builder.SumInt(context.Background(), "balance")
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, uint64(30), ba)
}
