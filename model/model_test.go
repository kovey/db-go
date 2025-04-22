package model

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/stretchr/testify/assert"
)

type test_model struct {
	*Model
	Id         int
	Age        int
	Name       string
	CreateTime string
	Sex        *int
}

func newTestmModel() *test_model {
	return &test_model{Model: NewModel("user", "id", Type_Int)}
}

func (t *test_model) Clone() ksql.RowInterface {
	return &test_model{Model: NewModel("user", "id", Type_Int)}
}

func (t *test_model) Columns() []string {
	return []string{"id", "age", "name", "create_time", "sex"}
}

func (t *test_model) Values() []any {
	return []any{&t.Id, &t.Age, &t.Name, &t.CreateTime, &t.Sex}
}

func (t *test_model) Save(ctx context.Context) error {
	return t.Model.Save(ctx, t)
}

func (t *test_model) Delete(ctx context.Context) error {
	return t.Model.Delete(ctx, t)
}

func TestModel(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := db.Open(testDb, "mysql")
	assert.Nil(t, err)
	m := &test_model{Model: NewModel("user", "id", Type_Int)}
	m.WithConn(conn)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `sex` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(m.Columns()).AddRow(1, 18, "kovey", "2025-04-03 11:11:11", 1))
	err = db.Model(m).Where("id", ksql.Eq, 1).First(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 18, m.Age)
	assert.Equal(t, "kovey", m.Name)
	assert.Equal(t, "2025-04-03 11:11:11", m.CreateTime)
	assert.Equal(t, 1, *m.Sex)
}

func TestModelSave(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := db.Open(testDb, "mysql")
	assert.Nil(t, err)
	m := &test_model{Model: NewModel("user", "id", Type_Int)}
	m.WithConn(conn)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `sex` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(m.Columns()).AddRow(1, 18, "kovey", "2025-04-03 11:11:11", 1))
	mock.ExpectPrepare("UPDATE `user` SET `age` = ?, `name` = ? WHERE `id` = ?").ExpectExec().WithArgs(19, "kovey save", 1).WillReturnResult(sqlmock.NewResult(0, 1))
	err = db.Model(m).Where("id", ksql.Eq, 1).First(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 18, m.Age)
	assert.Equal(t, "kovey", m.Name)
	assert.Equal(t, "2025-04-03 11:11:11", m.CreateTime)
	assert.Equal(t, 1, *m.Sex)

	m.Age = 19
	m.Name = "kovey save"
	err = m.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 19, m.Age)
	assert.Equal(t, "kovey save", m.Name)
	assert.Equal(t, "2025-04-03 11:11:11", m.CreateTime)
	assert.Equal(t, 1, *m.Sex)
	err = m.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestModelSaveFull(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := db.Open(testDb, "mysql")
	assert.Nil(t, err)
	m := &test_model{Model: NewModel("user", "id", Type_Int)}
	m.WithConn(conn)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `sex` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(m.Columns()).AddRow(1, 18, "kovey", "2025-04-03 11:11:11", 1))
	mock.ExpectPrepare("UPDATE `user` SET `age` = ?, `name` = ?, `create_time` = ?, `sex` = ? WHERE `id` = ?").ExpectExec().WithArgs(19, "kovey save", "2025-04-03 11:11:12", 2, 1).WillReturnResult(sqlmock.NewResult(0, 1))
	err = db.Model(m).Where("id", ksql.Eq, 1).First(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 18, m.Age)
	assert.Equal(t, "kovey", m.Name)
	assert.Equal(t, "2025-04-03 11:11:11", m.CreateTime)
	assert.Equal(t, 1, *m.Sex)

	m.Age = 19
	m.Name = "kovey save"
	m.CreateTime = "2025-04-03 11:11:12"
	sex := 2
	m.Sex = &sex
	err = m.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 19, m.Age)
	assert.Equal(t, "kovey save", m.Name)
	assert.Equal(t, "2025-04-03 11:11:12", m.CreateTime)
	assert.Equal(t, 2, *m.Sex)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestModelDelete(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := db.Open(testDb, "mysql")
	assert.Nil(t, err)
	m := &test_model{Model: NewModel("user", "id", Type_Int)}
	m.WithConn(conn)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `sex` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(m.Columns()).AddRow(1, 18, "kovey", "2025-04-03 11:11:11", 1))
	mock.ExpectPrepare("DELETE FROM `user` WHERE `id` = ?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	err = db.Model(m).Where("id", ksql.Eq, 1).First(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 18, m.Age)
	assert.Equal(t, "kovey", m.Name)
	assert.Equal(t, "2025-04-03 11:11:11", m.CreateTime)
	assert.Equal(t, 1, *m.Sex)

	err = m.Delete(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestModelRow(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := db.Open(testDb, "mysql")
	assert.Nil(t, err)
	m := &test_model{Model: NewModel("user", "id", Type_Int)}
	m.WithConn(conn)
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `sex` FROM `user` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(m.Columns()).AddRow(1, 18, "kovey", "2025-04-03 11:11:11", 1))
	mock.ExpectPrepare("UPDATE `user` SET `age` = ?, `name` = ? WHERE `id` = ?").ExpectExec().WithArgs(19, "kovey save", 1).WillReturnResult(sqlmock.NewResult(0, 1))
	err = Row(m).Where("id", ksql.Eq, 1).First(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 18, m.Age)
	assert.Equal(t, "kovey", m.Name)
	assert.Equal(t, "2025-04-03 11:11:11", m.CreateTime)
	assert.Equal(t, 1, *m.Sex)
	m.Age = 19
	m.Name = "kovey save"
	err = m.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 19, m.Age)
	assert.Equal(t, "kovey save", m.Name)
	assert.Equal(t, "2025-04-03 11:11:11", m.CreateTime)
	assert.Equal(t, 1, *m.Sex)
	err = m.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestModelRows(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := db.Open(testDb, "mysql")
	assert.Nil(t, err)
	var rows []*test_model
	mock.ExpectPrepare("SELECT `id`, `age`, `name`, `create_time`, `sex` FROM `user` WHERE `id` >= ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(newTestmModel().Columns()).AddRow(1, 18, "kovey", "2025-04-03 11:11:11", 1).AddRow(2, 19, "kovey22", "2025-04-03 12:11:11", 0))
	mock.ExpectPrepare("DELETE FROM `user` WHERE `id` = ?").ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectPrepare("DELETE FROM `user` WHERE `id` = ?").ExpectExec().WithArgs(2).WillReturnResult(sqlmock.NewResult(0, 1))
	err = Rows(&rows).WithConn(conn).Where("id", ksql.Ge, 1).All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	m := rows[0]
	assert.Nil(t, err)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 18, m.Age)
	assert.Equal(t, "kovey", m.Name)
	assert.Equal(t, "2025-04-03 11:11:11", m.CreateTime)
	assert.Equal(t, 1, *m.Sex)
	err = m.Delete(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	m = rows[1]
	assert.Nil(t, err)
	assert.Equal(t, 2, m.Id)
	assert.Equal(t, 19, m.Age)
	assert.Equal(t, "kovey22", m.Name)
	assert.Equal(t, "2025-04-03 12:11:11", m.CreateTime)
	assert.Equal(t, 0, *m.Sex)
	err = m.Delete(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestModelInsert(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := db.Open(testDb, "mysql")
	assert.Nil(t, err)
	m := newTestmModel()
	m.WithConn(conn)
	mock.ExpectPrepare("INSERT INTO `user` (`age`, `name`, `create_time`, `sex`) VALUES (?, ?, ?, ?)").ExpectExec().WithArgs(18, "kovey", "2025-04-03 11:11:11", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	m.Age = 18
	m.Name = "kovey"
	m.CreateTime = "2025-04-03 11:11:11"
	sex := 1
	m.Sex = &sex
	err = m.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 18, m.Age)
	assert.Equal(t, "kovey", m.Name)
	assert.Equal(t, "2025-04-03 11:11:11", m.CreateTime)
	assert.Equal(t, 1, *m.Sex)
	err = m.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestModelInsertNoAutoInc(t *testing.T) {
	testDb, mock, err := sqlmock.NewWithDSN("root:123456@tcp(127.0.0.1:3306)/test_dev?charset=utf8mb4&parseTime=true", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()

	conn, err := db.Open(testDb, "mysql")
	assert.Nil(t, err)
	m := newTestmModel()
	m.WithConn(conn)
	m.NoAutoInc()
	assert.True(t, m.Empty())
	mock.ExpectPrepare("INSERT INTO `user` (`id`, `age`, `name`, `create_time`, `sex`) VALUES (?, ?, ?, ?, ?)").ExpectExec().WithArgs(1, 18, "kovey", "2025-04-03 11:11:11", 1).WillReturnResult(sqlmock.NewResult(0, 1))
	m.Id = 1
	m.Age = 18
	m.Name = "kovey"
	m.CreateTime = "2025-04-03 11:11:11"
	sex := 1
	m.Sex = &sex
	err = m.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "id", m.PrimaryId())
	assert.True(t, !m.Empty())
	assert.Equal(t, 1, m.Id)
	assert.Equal(t, 18, m.Age)
	assert.Equal(t, "kovey", m.Name)
	assert.Equal(t, "2025-04-03 11:11:11", m.CreateTime)
	assert.Equal(t, 1, *m.Sex)
	err = m.Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, mock.ExpectationsWereMet())
}
