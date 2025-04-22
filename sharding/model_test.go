package sharding

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/model"
	"github.com/stretchr/testify/assert"
)

type test_model struct {
	*Model
	Id     int
	UserId int64
	Age    int
	Name   string
}

func newTestModel() *test_model {
	return &test_model{Model: NewModel("user", "id", model.Type_Int)}
}

func (t *test_model) Clone() ksql.RowInterface {
	return newTestModel()
}

func (t *test_model) Columns() []string {
	return []string{"id", "user_id", "age", "name"}
}

func (t *test_model) Values() []any {
	return []any{&t.Id, &t.UserId, &t.Age, &t.Name}
}

func (t *test_model) Save(ctx context.Context) error {
	return t.Model.Save(ctx, t)
}

func (t *test_model) Delete(ctx context.Context) error {
	return t.Model.Delete(ctx, t)
}

func TestModel(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()
	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)

	tm1 := newTestModel()
	tm2 := newTestModel()
	mock1.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_0` WHERE `id` = ?").ExpectQuery().WithArgs(2).WillReturnRows(sqlmock.NewRows(tm1.Columns()).AddRow(1, 2, 18, "kovey"))
	mock2.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_1` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(tm2.Columns()).AddRow(1, 1, 19, "kovey1"))

	err = Row(2, tm1).Where("id", ksql.Eq, 2).First(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, tm1.Id)
	assert.Equal(t, int64(2), tm1.UserId)
	assert.Equal(t, 18, tm1.Age)
	assert.Equal(t, "kovey", tm1.Name)
	err = Row(1, tm2).Where("id", ksql.Eq, 1).First(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, tm2.Id)
	assert.Equal(t, int64(1), tm2.UserId)
	assert.Equal(t, 19, tm2.Age)
	assert.Equal(t, "kovey1", tm2.Name)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestModelSave(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()
	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)

	tm1 := newTestModel()
	tm2 := newTestModel()
	mock1.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_0` WHERE `id` = ?").ExpectQuery().WithArgs(2).WillReturnRows(sqlmock.NewRows(tm1.Columns()).AddRow(1, 2, 18, "kovey"))
	mock1.ExpectPrepare("UPDATE `user_0` SET `age` = ?, `name` = ? WHERE `id` = ?").ExpectExec().WithArgs(19, "kovey save", 1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock2.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_1` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(tm2.Columns()).AddRow(1, 1, 19, "kovey1"))
	mock2.ExpectPrepare("UPDATE `user_1` SET `age` = ?, `name` = ? WHERE `id` = ?").ExpectExec().WithArgs(20, "kovey save", 1).WillReturnResult(sqlmock.NewResult(0, 1))

	err = Row(2, tm1).Where("id", ksql.Eq, 2).First(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, tm1.Id)
	assert.Equal(t, int64(2), tm1.UserId)
	assert.Equal(t, 18, tm1.Age)
	assert.Equal(t, "kovey", tm1.Name)
	err = Row(1, tm2).Where("id", ksql.Eq, 1).First(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, tm2.Id)
	assert.Equal(t, int64(1), tm2.UserId)
	assert.Equal(t, 19, tm2.Age)
	assert.Equal(t, "kovey1", tm2.Name)

	tm1.Age = 19
	tm1.Name = "kovey save"
	err = tm1.Save(context.Background())
	assert.Nil(t, err)
	err = tm1.Save(context.Background())
	assert.Nil(t, err)

	tm2.Age = 20
	tm2.Name = "kovey save"
	err = tm2.Save(context.Background())
	assert.Nil(t, err)
	err = tm2.Save(context.Background())
	assert.Nil(t, err)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}

func TestModelRows(t *testing.T) {
	testDb1, mock1, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	testDb2, mock2, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb1.Close()
	defer testDb2.Close()
	err = InitBy("mysql", []*sql.DB{testDb1, testDb2})
	assert.Nil(t, err)

	tm1 := newTestModel()
	tm2 := newTestModel()
	mock1.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_0` WHERE `id` = ?").ExpectQuery().WithArgs(2).WillReturnRows(sqlmock.NewRows(tm1.Columns()).AddRow(1, 2, 18, "kovey"))
	mock2.ExpectPrepare("SELECT `id`, `user_id`, `age`, `name` FROM `user_1` WHERE `id` = ?").ExpectQuery().WithArgs(1).WillReturnRows(sqlmock.NewRows(tm2.Columns()).AddRow(1, 1, 19, "kovey1"))

	var tm1s []*test_model
	var tm2s []*test_model
	err = Rows(2, &tm1s).Where("id", ksql.Eq, 2).All(context.Background())
	assert.Nil(t, err)
	tm1 = tm1s[0]
	assert.Equal(t, 1, tm1.Id)
	assert.Equal(t, int64(2), tm1.UserId)
	assert.Equal(t, 18, tm1.Age)
	assert.Equal(t, "kovey", tm1.Name)
	err = Rows(1, &tm2s).Where("id", ksql.Eq, 1).All(context.Background())
	assert.Nil(t, err)
	tm2 = tm2s[0]
	assert.Equal(t, 1, tm2.Id)
	assert.Equal(t, int64(1), tm2.UserId)
	assert.Equal(t, 19, tm2.Age)
	assert.Equal(t, "kovey1", tm2.Name)
	assert.Nil(t, mock1.ExpectationsWereMet())
	assert.Nil(t, mock2.ExpectationsWereMet())
}
