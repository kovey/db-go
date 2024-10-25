package sql

import (
	"github.com/kovey/db-go/v3"
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
	adds        []*table.Column
	drops       []string
	indexes     []*table.Index
	dropIndexes []string
	changes     []*table.Column
	changeOlds  []string
	comment     string
	charset     string
	collate     string
	engine      string
	table       string
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

func (a *Alter) AddPrimary(column string) ksql.AlterInterface {
	return a.AddIndex("", ksql.Index_Type_Primary, column)
}

func (a *Alter) Rename(as string) ksql.AlterInterface {
	a.builder.WriteString(" RENAME AS ")
	Backtick(as, &a.builder)

	return a
}
