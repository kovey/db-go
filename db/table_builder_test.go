package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestTableBuilderAlter(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)

	mock.ExpectPrepare("ALTER TABLE `user` COMMENT = 'user table', ADD COLUMN `user_id` BIGINT(20) DEFAULT '0' KEY COMMENT 'bigint', ADD COLUMN `images` BINARY(10240) COMMENT 'binary', ADD COLUMN `v_blob` BLOB COMMENT 'blob', ADD COLUMN `v_char` CHAR(10) NOT NULL COMMENT 'char', ADD COLUMN `add_column` VARCHAR(20) DEFAULT '' COMMENT 'add column', ADD COLUMN `v_date` DATE, ADD COLUMN `v_date_time` DATETIME(19), ADD COLUMN `d_decimal` DECIMAL(20,2) COMMENT 'decimal', ADD COLUMN `d_double` DOUBLE(20,2) COMMENT 'DOUBLE', ADD COLUMN `e_enum` ENUM('A','B','C'), ADD COLUMN `f_float` FLOAT(10,2), ADD COLUMN `geo_metry` GEOMETRY NULL, ADD UNIQUE INDEX `idx_date` USING 'BTREE' (`v_date_time`, `v_date`), ADD COLUMN `i_int` INT(11) DEFAULT '0', ADD COLUMN `line_string` LINESTRING, ADD COLUMN `p_point` POINT, ADD COLUMN `polygon` POLYGON, ADD PRIMARY KEY (`id_primary`), ADD COLUMN `s_set` SET('1','2'), ADD COLUMN `small_int` SMALLINT(3), ADD COLUMN `name` VARCHAR(10), ADD COLUMN `note` TEXT NULL, ADD COLUMN `create_time` TIMESTAMP(19) DEFAULT CURRENT_TIMESTAMP, ADD COLUMN `status` TINYINT(1) DEFAULT '1', ADD UNIQUE INDEX `uni_name` (`name`), CHANGE COLUMN `user_name` `user_nickname` VARCHAR(100), DROP COLUMN `sex`, DROP COLUMN `other`, DROP INDEX `idx_index`").ExpectExec().WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	tb := NewTableBuilder().Alter().Table("user").WithConn(conn)
	tb.AddBigInt("user_id").Index().Default("0").Comment("bigint")
	tb.AddBinary("images", 10240).Comment("binary")
	tb.AddBlob("v_blob").Comment("blob")
	tb.AddChar("v_char", 10).Comment("char").NotNullable()
	tb.AddColumn("add_column", "varchar", 20, 0).Default("").Comment("add column")
	tb.AddDate("v_date")
	tb.AddDateTime("v_date_time")
	tb.AddDecimal("d_decimal", 20, 2).Comment("decimal")
	tb.AddDouble("d_double", 20, 2).Comment("DOUBLE")
	tb.AddEnum("e_enum", []string{"A", "B", "C"})
	tb.AddFloat("f_float", 10, 2)
	tb.AddGeoMetry("geo_metry").Nullable()
	tb.AddIndex("idx_date").Unique().Algorithm(ksql.Index_Alg_BTree).Columns("v_date_time", "v_date")
	tb.AddInt("i_int").Default("0")
	tb.AddLineString("line_string")
	tb.AddPoint("p_point")
	tb.AddPolygon("polygon")
	tb.AddPrimary("id_primary")
	tb.AddSet("s_set", []string{"1", "2"})
	tb.AddSmallInt("small_int")
	tb.AddString("name", 10)
	tb.AddText("note").Nullable()
	tb.AddTimestamp("create_time").UseCurrent()
	tb.AddTinyInt("status").Default("1")
	tb.AddUnique("uni_name", "name")
	tb.ChangeColumn("user_name", "user_nickname", "varchar", 100, 0)
	tb.Charset("utf8").Collate("123").Comment("user table")
	tb.DropColumn("sex")
	tb.DropColumnIfExists("other")
	tb.DropIndex("idx_index")
	tb.Engine("InnoDB")
	err = tb.Exec(context.Background())
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestTableBuilderCreate(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)

	mock.ExpectPrepare("CREATE TABLE `user` (`user_id` BIGINT(20) DEFAULT '0' KEY COMMENT 'bigint', `images` BINARY(10240) COMMENT 'binary', `v_blob` BLOB COMMENT 'blob', `v_char` CHAR(10) NOT NULL COMMENT 'char', `add_column` VARCHAR(20) DEFAULT '' COMMENT 'add column', `v_date` DATE, `v_date_time` DATETIME(19), `d_decimal` DECIMAL(20,2) COMMENT 'decimal', `d_double` DOUBLE(20,2) COMMENT 'DOUBLE', `e_enum` ENUM('A','B','C'), `f_float` FLOAT(10,2), `geo_metry` GEOMETRY NULL, `i_int` INT(11) DEFAULT '0', `line_string` LINESTRING, `p_point` POINT, `polygon` POLYGON, `s_set` SET('1','2'), `small_int` SMALLINT(3), `name` VARCHAR(10), `note` TEXT NULL, `create_time` TIMESTAMP(19) DEFAULT CURRENT_TIMESTAMP, `status` TINYINT(1) DEFAULT '1', UNIQUE INDEX `idx_date` USING 'BTREE' (`v_date_time`, `v_date`), PRIMARY KEY (`id_primary`), UNIQUE INDEX `uni_name` (`name`)) CHARACTER SET = utf8, COLLATE = 123, COMMENT = 'user table', ENGINE = InnoDB").ExpectExec().WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	tb := NewTableBuilder().Create().Table("user").WithConn(conn)
	tb.AddBigInt("user_id").Index().Default("0").Comment("bigint")
	tb.AddBinary("images", 10240).Comment("binary")
	tb.AddBlob("v_blob").Comment("blob")
	tb.AddChar("v_char", 10).Comment("char").NotNullable()
	tb.AddColumn("add_column", "varchar", 20, 0).Default("").Comment("add column")
	tb.AddDate("v_date")
	tb.AddDateTime("v_date_time")
	tb.AddDecimal("d_decimal", 20, 2).Comment("decimal")
	tb.AddDouble("d_double", 20, 2).Comment("DOUBLE")
	tb.AddEnum("e_enum", []string{"A", "B", "C"})
	tb.AddFloat("f_float", 10, 2)
	tb.AddGeoMetry("geo_metry").Nullable()
	tb.AddIndex("idx_date").Unique().Algorithm(ksql.Index_Alg_BTree).Columns("v_date_time", "v_date")
	tb.AddInt("i_int").Default("0")
	tb.AddLineString("line_string")
	tb.AddPoint("p_point")
	tb.AddPolygon("polygon")
	tb.AddPrimary("id_primary")
	tb.AddSet("s_set", []string{"1", "2"})
	tb.AddSmallInt("small_int")
	tb.AddString("name", 10)
	tb.AddText("note").Nullable()
	tb.AddTimestamp("create_time").UseCurrent()
	tb.AddTinyInt("status").Default("1")
	tb.AddUnique("uni_name", "name")
	tb.Charset("utf8").Collate("123").Comment("user table").Engine("InnoDB")
	err = tb.Exec(context.Background())
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestTableBuilderCreateFrom(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)

	mock.ExpectPrepare("CREATE TABLE IF NOT EXISTS `user` AS SELECT `id`, `name`, `age` FROM `user_ext` WHERE `id` > ?").ExpectExec().WithArgs(1000).WillReturnResult(sqlmock.NewResult(0, 1))
	sub := NewQuery().Table("user_ext").Columns("id", "name", "age").Where("id", ">", 1000)
	tb := NewTableBuilder().Create().Table("user").WithConn(conn)
	tb.IfNotExists().From(sub)
	err = tb.Exec(context.Background())
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestTableBuilderCreateLike(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)

	mock.ExpectPrepare("CREATE TABLE IF NOT EXISTS `user` (LIKE `user_tpl`)").ExpectExec().WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	tb := NewTableBuilder().Create().Table("user").WithConn(conn)
	tb.IfNotExists().Like("user_tpl")
	err = tb.Exec(context.Background())
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestTableBuilderHasColumn(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("SHOW COLUMNS FROM `user` LIKE 'id'").ExpectQuery().WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{"field", "type", "null", "key", "default", "extra"}).AddRow("id", "int", "NO", "PRI", nil, "auto_increment"))
	ta := NewTableBuilder().Table("user")
	has, err := ta.HasColumn(context.Background(), "id")
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.True(t, has)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestTableBuilderHasIndex(t *testing.T) {
	testDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer testDb.Close()
	conn, err := Open(testDb, "mysql")
	assert.Nil(t, err)
	database = conn

	mock.ExpectPrepare("SHOW INDEX FROM `user` WHERE Key_name = ?").ExpectQuery().WithArgs("PRIMARY").WillReturnRows(sqlmock.NewRows([]string{"Table", "Non_unique", "Key_name", "Seq_in_index", "Column_name", "Collation", "Cardinality", "Sub_part", "Packed", "Null", "Index_type", "Comment", "Index_comment", "Visible", "Expression"}).AddRow("user", 0, "PRIMARY", 1, "id", "A", 2, nil, nil, "", "BTREE", "", "", "YES", nil))
	ta := NewTableBuilder().Table("user")
	has, err := ta.HasIndex(context.Background(), "PRIMARY")
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.True(t, has)
	assert.Nil(t, mock.ExpectationsWereMet())
}
