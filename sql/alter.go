package sql

import (
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/table"
)

//ALTER TABLE `skillw-ios`.`activity_record`
//DROP COLUMN `others`,
//ADD COLUMN `atest` VARCHAR(45) NOT NULL DEFAULT '' AFTER `updated_at`,
//ADD COLUMN `btest` VARCHAR(45) NOT NULL DEFAULT '' AFTER `atest`;

// ALTER TABLE `skillw-ios`.`activity_record`
// ADD UNIQUE INDEX `idx_test` (`activity_id` ASC, `date_end` ASC) VISIBLE,
// ADD INDEX `idx_normal` (`param3` ASC) VISIBLE,
// ADD INDEX `idx_more` (`param4` ASC, `param5` ASC) VISIBLE,
// DROP INDEX `date` ;

// ALTER TABLE `skillw-ios`.`activity_record`
// CHANGE COLUMN `param6` `param7` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '餐宿' ,
// ADD UNIQUE INDEX `param7_UNIQUE` (`param7` ASC) VISIBLE;
// COMMENT = '活动记录'

// ALTER TABLE `skillw-ios`.`abtest_configs`
// CHARACTER SET = utf8 , COLLATE = utf8_latvian_ci,ENGINE = MEMORY  ;

type Alter struct {
	*base
	adds          []*table.Column
	dropsIfExists []string
	drops         []string
	indexes       []*table.Index
	dropIndexes   []string
	changes       []*table.Column
	changeOlds    []string
	comment       string
	charset       string
	collate       string
	engine        string
	table         string
}

func NewAlter() *Alter {
	u := &Alter{base: &base{hasPrepared: false}}
	u.keyword("ALTER TABLE ")
	return u
}

func (u *Alter) Table(table string) ksql.AlterInterface {
	u.table = table
	return u
}

func (a *Alter) ChangeColumn(oldColumn, newColumn, t string, length, scale int, sets ...string) ksql.ColumnInterface {
	c := table.ParseType(t, length, scale, sets...)
	if c == nil {
		return nil
	}

	col := table.NewColumn(newColumn, c)
	a.changes = append(a.changes, col)
	a.changeOlds = append(a.changeOlds, oldColumn)
	return col
}

func (a *Alter) AddIndex(name string, t ksql.IndexType, column ...string) ksql.AlterInterface {
	index := &table.Index{Name: name, Type: t}
	index.Columns(column...)
	a.indexes = append(a.indexes, index)
	return a
}

func (a *Alter) DropIndex(name string) ksql.AlterInterface {
	a.dropIndexes = append(a.dropIndexes, name)
	return a
}

func (a *Alter) AddDecimal(column string, length, scale int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Decimal, length, scale)
}

func (a *Alter) AddDouble(column string, length, scale int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Double, length, scale)
}

func (a *Alter) AddFloat(column string, length, scale int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Float, length, scale)
}

func (a *Alter) AddBinary(column string, length int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Binary, length, 0)
}

func (a *Alter) AddGeoMetry(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_GeoMetry, 0, 0)
}

func (a *Alter) AddPolygon(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Polygon, 0, 0)
}

func (a *Alter) AddPoint(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Point, 0, 0)
}

func (a *Alter) AddLineString(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_LineString, 0, 0)
}

func (a *Alter) AddBlob(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Blob, 0, 0)
}

func (a *Alter) AddText(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Text, 0, 0)
}

func (a *Alter) AddSet(column string, sets []string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Set, 0, 0, sets...)
}

func (a *Alter) AddEnum(column string, options []string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Enum, 0, 0, options...)
}

func (a *Alter) AddDate(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Date, 0, 0)
}

func (a *Alter) AddDateTime(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_DateTime, 0, 0)
}

func (a *Alter) AddTimestamp(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Timestamp, 0, 0)
}

func (a *Alter) AddSmallInt(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_SmallInt, 3, 0)
}
func (a *Alter) AddTinyInt(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_TinyInt, 1, 0)
}

func (a *Alter) AddBigInt(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_BigInt, 20, 0)
}

func (a *Alter) AddInt(column string) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Int, 11, 0)
}

func (a *Alter) AddString(column string, length int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_VarChar, length, 0)
}

func (a *Alter) AddChar(column string, length int) ksql.ColumnInterface {
	return a.AddColumn(column, table.Type_Char, length, 0)
}

func (u *Alter) AddColumn(column, t string, length, scale int, sets ...string) ksql.ColumnInterface {
	c := table.ParseType(t, length, scale, sets...)
	if c == nil {
		return nil
	}

	col := table.NewColumn(column, c)
	u.adds = append(u.adds, col)
	return col
}

func (u *Alter) DropColumn(column string) ksql.AlterInterface {
	u.drops = append(u.drops, column)
	return u
}

func (u *Alter) DropColumnIfExists(column string) ksql.AlterInterface {
	u.dropsIfExists = append(u.dropsIfExists, column)
	return u
}

func (u *Alter) Comment(comment string) ksql.AlterInterface {
	u.comment = comment
	return u
}

func (u *Alter) Charset(charset string) ksql.AlterInterface {
	u.charset = charset
	return u
}

func (u *Alter) Collate(collate string) ksql.AlterInterface {
	u.collate = collate
	return u
}

func (u *Alter) Engine(engine string) ksql.AlterInterface {
	u.engine = engine
	return u
}

func (u *Alter) Prepare() string {
	if u.hasPrepared {
		return u.base.Prepare()
	}

	u.hasPrepared = true
	Column(u.table, &u.builder)
	u.builder.WriteString(" ")
	canAdd := false

	for index, drop := range u.drops {
		canAdd = true
		if index > 0 {
			u.builder.WriteString(",")
		}

		u.builder.WriteString("DROP COLUMN ")
		Column(drop, &u.builder)
	}

	for index, drop := range u.dropsIfExists {
		canAdd = true
		if index > 0 {
			u.builder.WriteString(",")
		}

		u.builder.WriteString("DROP COLUMN IF EXISTS ")
		Column(drop, &u.builder)
	}

	for idx, drop := range u.dropIndexes {
		if canAdd || idx > 0 {
			u.builder.WriteString(",")
		}

		canAdd = true
		u.builder.WriteString("DROP INDEX ")
		Column(drop, &u.builder)
	}

	for index, add := range u.adds {
		if canAdd || index > 0 {
			u.builder.WriteString(",")
		}

		canAdd = true
		u.builder.WriteString("ADD COLUMN ")
		u.builder.WriteString(add.Express())
	}

	for index, change := range u.changes {
		if canAdd || index > 0 {
			u.builder.WriteString(",")
		}
		canAdd = true
		u.builder.WriteString("CHANGE COLUMN ")
		Column(u.changeOlds[index], &u.builder)
		u.builder.WriteString(" ")
		u.builder.WriteString(change.Express())
	}

	for index, add := range u.indexes {
		if canAdd || index > 0 {
			u.builder.WriteString(",")
		}

		canAdd = true
		u.builder.WriteString(add.AlterExpress())
	}

	u._write("CHARACTER SET", u.charset, &canAdd)
	u._write("COLLATE", u.collate, &canAdd)
	u._write("ENGINE", u.engine, &canAdd)
	u._write("COMMENT", u.comment, &canAdd)

	return u.base.Prepare()
}

func (a *Alter) _write(key, val string, canAdd *bool) {
	if val == "" {
		return
	}

	if *canAdd {
		a.builder.WriteString(",")
	}
	*canAdd = true
	a.builder.WriteString(key)
	a.builder.WriteString(" = ")
	Quote(val, &a.builder)
}

func (a *Alter) AddUnique(name string, columns ...string) ksql.AlterInterface {
	return a.AddIndex(name, ksql.Index_Type_Unique, columns...)
}

func (a *Alter) AddPrimary(column string) ksql.AlterInterface {
	return a.AddIndex("", ksql.Index_Type_Primary, column)
}

func (a *Alter) Rename(as string) ksql.AlterInterface {
	a.builder.WriteString(" RENAME AS ")
	Backtick(as, &a.builder)

	return a
}
