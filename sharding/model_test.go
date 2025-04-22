package sharding

import (
	"context"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/model"
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
