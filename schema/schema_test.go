package schema

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	err = db.InitBy(testDb, "mysql")
	assert.Nil(t, err)

	mock.ExpectPrepare("ALTER SCHEMA `user_db` DEFAULT CHARACTER SET = utf8mb4 DEFAULT ENCRYPTION = 'Y'").ExpectExec().WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	err = Schema(context.Background(), "user_db", func(schema ksql.SchemaInterface) {
		schema.Charset("utf8mb4").Encryption("Y")
	})
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestSchemaCreate(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	err = db.InitBy(testDb, "mysql")
	assert.Nil(t, err)

	mock.ExpectPrepare("CREATE SCHEMA IF NOT EXISTS `user_db` DEFAULT CHARACTER SET = utf8mb4 ENCRYPTION = 'Y'").ExpectExec().WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	err = Create(context.Background(), "user_db", func(schema ksql.SchemaInterface) {
		schema.Charset("utf8mb4").Encryption("Y").IfNotExists()
	})
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
