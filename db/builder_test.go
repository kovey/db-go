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

func (t *test_user) Values() []any {
	return []any{&t.Id, &t.Age, &t.Name, &t.CreateTime, &t.Balance}
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
	testDb, mock, err := sqlmock.New()
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
	mock.ExpectPrepare(builder.query.Prepare()).ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23))
	err = builder.First(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tu.Id)
	assert.Equal(t, 18, tu.Age)
	assert.Equal(t, "kovey", tu.Name)
	assert.Equal(t, now.Format(time.DateTime), tu.CreateTime.Format(time.DateTime))
	assert.Equal(t, float64(30.23), tu.Balance)
}

func TestBuilderAll(t *testing.T) {
	testDb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	columns := []string{"id", "age", "name", "create_time", "balance"}
	now, _ := time.Parse(time.DateTime, "2025-04-03 11:11:11")
	var rows []*test_user
	builder := Rows(&rows).(*Builder[*test_user])
	builder.Table("user").Columns(columns...).Where("id", ksql.Ge, 1)
	mock.ExpectPrepare(builder.query.Prepare()).ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23).AddRow(2, 15, "test", now, 34.13))
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
	testDb, mock, err := sqlmock.New()
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
	mock.ExpectPrepare("SELECT `u`.`id`,`us`.`age`,`us`.`name`,`us`.`create_time`,`us`.`balance`,`u`.`count` FROM").ExpectQuery().WithArgs(100, 200, 1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 18, "kovey", now, 30.23, 300).AddRow(2, 15, "test", now, 34.13, 100))
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
