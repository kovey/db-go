//go:build integration
// +build integration

// Package integration contains end-to-end tests that require a real MySQL instance.
// These tests are skipped during regular `go test` runs.
//
// Set MYSQL_DSN to point at a test database:
//
//	go test -tags=integration -count=1 ./integration/
//
// In CI, a MySQL service container is used (see .github/workflows/go-test.yml).
package integration

import (
	"context"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/kovey/db-go/v3/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func dsn() string {
	if d := os.Getenv("MYSQL_DSN"); d != "" {
		return d
	}
	// default for local development
	return "root:test@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=true"
}

func setupDB(t *testing.T) {
	t.Helper()

	err := db.Init(db.Config{
		DriverName:     "mysql",
		DataSourceName: dsn(),
		MaxIdleTime:    30 * time.Second,
		MaxLifeTime:    60 * time.Second,
		MaxIdleConns:   2,
		MaxOpenConns:   5,
	})
	require.NoError(t, err)

	t.Cleanup(func() { db.Close() })
}

func setupTable(t *testing.T) {
	t.Helper()
	err := db.Create(context.Background(), "user", func(table ksql.TableInterface) {
		table.AddBigInt("id").AutoIncrement().Unsigned().Comment("Primary key")
		table.AddString("name", 63).Default("").Comment("User name")
		table.AddInt("age").Default("0").Comment("Age")
		table.AddTinyInt("status").Default("0").Comment("Status")
		table.AddPrimary("id")
		table.Engine("InnoDB").Charset("utf8mb4")
	})
	require.NoError(t, err)
	t.Cleanup(func() { db.DropTable(context.Background(), "user") })
}

// ─────────────────────────────────────────────
// Test model
// ─────────────────────────────────────────────

type testUser struct {
	*model.Model
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Status int    `json:"status"`
}

func newTestUser() *testUser {
	return &testUser{Model: model.NewModel("user", "id", model.Type_Int)}
}

func (u *testUser) Clone() ksql.RowInterface { return newTestUser() }

func (u *testUser) Columns() []string {
	return []string{"id", "name", "age", "status"}
}

func (u *testUser) Values() []any {
	return []any{&u.Id, &u.Name, &u.Age, &u.Status}
}

func (u *testUser) Save(ctx context.Context) error {
	return u.Model.SaveBy(ctx, u)
}

func (u *testUser) Delete(ctx context.Context) error {
	return u.Model.DeleteBy(ctx, u)
}

// ─────────────────────────────────────────────
// Tests
// ─────────────────────────────────────────────

func TestCRUD_InsertAndFind(t *testing.T) {
	setupDB(t)
	ctx := context.Background()
	setupTable(t)

	u := newTestUser()
	u.Name = "alice"
	u.Age = 30
	u.Status = 1
	err := u.Save(ctx)
	require.NoError(t, err)
	assert.Greater(t, u.Id, int64(0))

	u2 := newTestUser()
	err = db.Model(u2).Where("id", ksql.Eq, u.Id).First(ctx)
	require.NoError(t, err)
	assert.Equal(t, "alice", u2.Name)
	assert.Equal(t, 30, u2.Age)
	assert.Equal(t, 1, u2.Status)
}

func TestCRUD_Update(t *testing.T) {
	setupDB(t)
	ctx := context.Background()
	setupTable(t)

	u := newTestUser()
	u.Name = "bob"
	u.Age = 25
	require.NoError(t, u.Save(ctx))

	u2 := newTestUser()
	require.NoError(t, db.Model(u2).Where("id", ksql.Eq, u.Id).First(ctx))

	u2.Name = "bob updated"
	u2.Age = 26
	require.NoError(t, u2.Save(ctx))

	u3 := newTestUser()
	require.NoError(t, db.Model(u3).Where("id", ksql.Eq, u.Id).First(ctx))
	assert.Equal(t, "bob updated", u3.Name)
	assert.Equal(t, 26, u3.Age)
}

func TestCRUD_Delete(t *testing.T) {
	setupDB(t)
	ctx := context.Background()
	setupTable(t)

	u := newTestUser()
	u.Name = "charlie"
	require.NoError(t, u.Save(ctx))

	u2 := newTestUser()
	require.NoError(t, db.Model(u2).Where("id", ksql.Eq, u.Id).First(ctx))
	require.NoError(t, u2.Delete(ctx))

	u3 := newTestUser()
	err := db.Model(u3).Where("id", ksql.Eq, u.Id).First(ctx)
	assert.Error(t, err)
}

func TestCRUD_BatchQuery(t *testing.T) {
	setupDB(t)
	ctx := context.Background()
	setupTable(t)

	for _, name := range []string{"u1", "u2", "u3"} {
		u := newTestUser()
		u.Name = name
		u.Age = 20
		require.NoError(t, u.Save(ctx))
	}

	var users []*testUser
	err := db.Models(&users).Where("age", ksql.Eq, 20).All(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 3)
}

func TestTransaction_Commit(t *testing.T) {
	setupDB(t)
	ctx := context.Background()
	setupTable(t)

	err := db.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
		u := newTestUser()
		u.Name = "tx_user"
		return u.Save(ctx)
	})
	require.NoError(t, err)

	u := newTestUser()
	if err := db.Model(u).Where("name", ksql.Eq, "tx_user").First(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestTransaction_Rollback(t *testing.T) {
	setupDB(t)
	ctx := context.Background()
	setupTable(t)

	_ = db.Transaction(ctx, func(ctx context.Context, conn ksql.ConnectionInterface) error {
		u := newTestUser()
		u.Name = "rolled_back"
		u.Save(ctx)
		_, err := db.ExecRaw(ctx, db.Raw("SELECT * FROM nonexistent_table"))
		return err
	})

	var users []*testUser
	err := db.Models(&users).Where("name", ksql.Eq, "rolled_back").All(ctx)
	require.NoError(t, err)
	assert.Empty(t, users)
}

func TestHasTable_HasColumn(t *testing.T) {
	setupDB(t)
	ctx := context.Background()
	setupTable(t)

	has, err := db.HasTable(ctx, "user")
	require.NoError(t, err)
	assert.True(t, has)

	has, err = db.HasColumn(ctx, "user", "id")
	require.NoError(t, err)
	assert.True(t, has)

	has, err = db.HasColumn(ctx, "user", "nonexistent")
	require.NoError(t, err)
	assert.False(t, has)
}

func TestPagination(t *testing.T) {
	setupDB(t)
	ctx := context.Background()
	setupTable(t)

	for i := 0; i < 10; i++ {
		u := newTestUser()
		u.Name = "page_user"
		u.Age = 30
		require.NoError(t, u.Save(ctx))
	}

	pageInfo, err := db.Models(&[]*testUser{}).
		Where("age", ksql.Eq, 30).
		Pagination(ctx, 1, 3)
	require.NoError(t, err)

	assert.Equal(t, uint64(10), pageInfo.TotalCount())
	assert.Equal(t, uint64(4), pageInfo.TotalPage())
	assert.Len(t, pageInfo.List(), 3)
}

func TestCount(t *testing.T) {
	setupDB(t)
	ctx := context.Background()
	setupTable(t)

	for i := 0; i < 5; i++ {
		u := newTestUser()
		u.Name = "count_user"
		u.Status = 1
		require.NoError(t, u.Save(ctx))
	}

	count, err := db.Models(&[]*testUser{}).Where("status", ksql.Eq, 1).Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(5), count)
}
