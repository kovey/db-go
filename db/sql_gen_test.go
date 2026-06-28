package db

import (
	"strings"
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql"
	"github.com/stretchr/testify/assert"
)

// buildWhere is a test helper to get the SQL output of a Where clause.
func buildWhere(w ksql.WhereInterface) string {
	var b strings.Builder
	w.Build(&b)
	return b.String()
}

// buildHaving is a test helper to get the SQL output of a Having clause.
func buildHaving(h ksql.HavingInterface) string {
	var b strings.Builder
	h.Build(&b)
	return b.String()
}

// test_user is a minimal RowInterface implementation for tests.
type test_user struct {
	*Row
	Id         int64
	Age        int
	Name       string
	CreateTime string
	Balance    float64
}

func newTestUser() *test_user {
	return &test_user{Row: &Row{}}
}

func (t *test_user) Clone() ksql.RowInterface { return newTestUser() }
func (t *test_user) Columns() []string        { return []string{"id", "age", "name", "create_time", "balance"} }
func (t *test_user) Values() []any            { return []any{&t.Id, &t.Age, &t.Name, &t.CreateTime, &t.Balance} }

// ─────────────────────────────────────────────
// Layer 1 — SQL generation assertions
// No sqlmock, no database — just verify the SQL strings and binds.
// ─────────────────────────────────────────────

func TestNewQuery_Select_Basic(t *testing.T) {
	q := sql.NewQuery().Table("user").Columns("id", "name", "age")
	assert.Equal(t, "SELECT `id`, `name`, `age` FROM `user`", q.Prepare())
	assert.Nil(t, q.Binds())
}

func TestNewQuery_Select_AllColumns(t *testing.T) {
	q := sql.NewQuery().Table("user")
	assert.Equal(t, "SELECT * FROM `user`", q.Prepare())
}

func TestNewQuery_Select_WhereEq(t *testing.T) {
	q := sql.NewQuery().Table("user").Columns("id", "name").Where("status", ksql.Eq, 0)
	assert.Equal(t, "SELECT `id`, `name` FROM `user` WHERE `status` = ?", q.Prepare())
	assert.Equal(t, []any{0}, q.Binds())
}

func TestNewQuery_Select_WhereOps(t *testing.T) {
	tests := []struct {
		name string
		op   ksql.Op
		sql  string
	}{
		{"eq", ksql.Eq, "`id` = ?"},
		{"neq", ksql.Neq, "`id` <> ?"},
		{"gt", ksql.Gt, "`id` > ?"},
		{"ge", ksql.Ge, "`id` >= ?"},
		{"lt", ksql.Lt, "`id` < ?"},
		{"le", ksql.Le, "`id` <= ?"},
		{"like", ksql.Like, "`id` LIKE ?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := sql.NewQuery().Table("t").Columns("id").Where("id", tt.op, 1)
			assert.Contains(t, q.Prepare(), tt.sql)
		})
	}
}

func TestNewQuery_Select_WhereIsNull(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").WhereIsNull("deleted_at").WhereIsNotNull("email")
	assert.Contains(t, q.Prepare(), "`deleted_at` IS NULL")
	assert.Contains(t, q.Prepare(), "`email` IS NOT NULL")
}

func TestNewQuery_Select_WhereIn_WhereNotIn(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").
		WhereIn("role", []any{"a", "b"}).
		WhereNotIn("status", []any{1, 2, 3})

	assert.Contains(t, q.Prepare(), "`role` IN (?, ?)")
	assert.Contains(t, q.Prepare(), "`status` NOT IN (?, ?, ?)")
	assert.Equal(t, []any{"a", "b", 1, 2, 3}, q.Binds())
}

func TestNewQuery_Select_Between(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").
		Between("age", 10, 20).
		NotBetween("score", 1000, 2000)

	assert.Contains(t, q.Prepare(), "`age` BETWEEN ? AND ?")
	assert.Contains(t, q.Prepare(), "`score` NOT BETWEEN ? AND ?")
	assert.Equal(t, []any{10, 20, 1000, 2000}, q.Binds())
}

func TestNewQuery_Select_WhereInBy(t *testing.T) {
	sub := sql.NewQuery().Table("orders").Columns("user_id").Where("total", ksql.Gt, 100)
	q := sql.NewQuery().Table("user").Columns("id").
		WhereInBy("id", sub).
		WhereNotInBy("id", sub)

	assert.Contains(t, q.Prepare(), "`id` IN (SELECT `user_id` FROM `orders` WHERE `total` > ?)")
	assert.Contains(t, q.Prepare(), "`id` NOT IN (SELECT `user_id` FROM `orders` WHERE `total` > ?)")
	assert.Equal(t, []any{100, 100}, q.Binds())
}

func TestNewQuery_Select_OrWhere(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").
		Where("status", ksql.Eq, 1).
		OrWhere(func(w ksql.WhereInterface) {
			w.Between("score", 10, 20).Where("vip", ksql.Eq, 1)
		})

	assert.Contains(t, q.Prepare(), "OR (")
	assert.Equal(t, []any{1, 10, 20, 1}, q.Binds())
}

func TestNewQuery_Select_AndWhere(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").
		Where("status", ksql.Eq, 1).
		AndWhere(func(w ksql.WhereInterface) {
			w.Where("age", ksql.Ge, 18).Where("age", ksql.Le, 65)
		})

	assert.Contains(t, q.Prepare(), "AND (")
	assert.Equal(t, []any{1, 18, 65}, q.Binds())
}

func TestNewQuery_Select_WhereExpress(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").
		WhereExpress(Raw("JSON_EXTRACT(meta, '$.vip') = ?", true))

	assert.Contains(t, q.Prepare(), "JSON_EXTRACT(meta, '$.vip') = ?")
	assert.Equal(t, []any{true}, q.Binds())
}

func TestNewQuery_Select_Join(t *testing.T) {
	q := sql.NewQuery().Table("user").As("u").Columns("u.id", "p.bio")
	q.Join("profile").As("p").On("p.user_id", "=", "u.id")
	q.LeftJoin("settings").As("s").On("s.user_id", "=", "u.id")
	q.RightJoin("log").As("l").On("l.user_id", "=", "u.id")

	sqlStr := q.Prepare()
	assert.Contains(t, sqlStr, "INNER JOIN `profile` AS `p`")
	assert.Contains(t, sqlStr, "LEFT JOIN `settings` AS `s`")
	assert.Contains(t, sqlStr, "RIGHT JOIN `log` AS `l`")
}

func TestNewQuery_Select_JoinExpress(t *testing.T) {
	q := sql.NewQuery().Table("user").As("u").Columns("u.id")
	q.JoinExpress(Raw("JOIN extra AS e ON e.id = u.id AND e.status = ?", 1))

	assert.Contains(t, q.Prepare(), "JOIN extra AS e ON e.id = u.id AND e.status = ?")
	assert.Equal(t, []any{1}, q.Binds())
}

func TestNewQuery_Select_GroupBy(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("name", "COUNT(1) as cnt").
		Group("name").GroupWithRollUp()

	assert.Contains(t, q.Prepare(), "GROUP BY `name`")
	assert.Contains(t, q.Prepare(), "WITH ROLLUP")
}

func TestNewQuery_Select_OrderBy(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").
		Order("name").OrderDesc("score")

	assert.Contains(t, q.Prepare(), "ORDER BY `name` ASC, `score` DESC")
}

func TestNewQuery_Select_LimitOffset(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").Limit(10).Offset(20)

	assert.Contains(t, q.Prepare(), "LIMIT ?")
	assert.Contains(t, q.Prepare(), "OFFSET ?")
	assert.Equal(t, []any{10, 20}, q.Binds())
}

func TestNewQuery_Select_Distinct(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("role").Distinct()
	assert.Contains(t, q.Prepare(), "SELECT DISTINCT")
	assert.NotContains(t, q.Prepare(), "SELECT DISTINCT *")
}

func TestNewQuery_Select_Func(t *testing.T) {
	q := sql.NewQuery().Table("t").Func("SUM", "balance", "total")
	assert.Contains(t, q.Prepare(), "SUM(`balance`) AS `total`")
}

func TestNewQuery_Select_DistinctFunc(t *testing.T) {
	q := sql.NewQuery().Table("t").FuncDistinct("COUNT", "id", "cnt")
	assert.Contains(t, q.Prepare(), "DISTINCT")
	assert.Contains(t, q.Prepare(), "`cnt`")
}

func TestNewQuery_Select_ColumnWithAs(t *testing.T) {
	q := sql.NewQuery().Table("t").Column("nickname", "name")
	assert.Contains(t, q.Prepare(), "`nickname` AS `name`")
}

func TestNewQuery_Select_ColumnsExpress(t *testing.T) {
	q := sql.NewQuery().Table("t").ColumnsExpress(Raw("COUNT(1) as cnt"))
	assert.Contains(t, q.Prepare(), "COUNT(1) as cnt")
	assert.Nil(t, q.Binds())
}

func TestNewQuery_Select_SubqueryInFrom(t *testing.T) {
	sub := sql.NewQuery().Table("orders").Columns("user_id", "COUNT(1) as cnt").Group("user_id")
	q := sql.NewQuery().TableBy(sub, "o").Columns("o.user_id", "o.cnt")

	sqlStr := q.Prepare()
	assert.Contains(t, sqlStr, "FROM (SELECT")
	assert.Contains(t, sqlStr, "`o`.`user_id`")
}

func TestNewQuery_Select_ForUpdate(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id")
	q.For().Update()
	assert.Contains(t, q.Prepare(), "FOR UPDATE")
}

func TestNewQuery_Select_ForShare(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id")
	q.For().Share()
	assert.Contains(t, q.Prepare(), "FOR SHARE")
}

func TestNewQuery_Select_ForUpdateNowait(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id")
	q.For().Update().NoWait()
	assert.Contains(t, q.Prepare(), "FOR UPDATE NOWAIT")
}

func TestNewQuery_Select_ForUpdateSkipLocked(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id")
	q.For().Update().SkipLocked()
	assert.Contains(t, q.Prepare(), "FOR UPDATE SKIP LOCKED")
}

func TestNewQuery_Select_Modifiers(t *testing.T) {
	tests := []struct {
		name   string
		modFn  func(q ksql.QueryInterface)
		expect string
	}{
		{"high_priority", func(q ksql.QueryInterface) { q.HighPriority() }, "HIGH_PRIORITY"},
		{"straight_join", func(q ksql.QueryInterface) { q.StraightJoin() }, "STRAIGHT_JOIN"},
		{"sql_small_result", func(q ksql.QueryInterface) { q.SqlSmallResult() }, "SQL_SMALL_RESULT"},
		{"sql_big_result", func(q ksql.QueryInterface) { q.SqlBigResult() }, "SQL_BIG_RESULT"},
		{"sql_buffer_result", func(q ksql.QueryInterface) { q.SqlBufferResult() }, "SQL_BUFFER_RESULT"},
		{"sql_no_cache", func(q ksql.QueryInterface) { q.SqlNoCache() }, "SQL_NO_CACHE"},
		{"sql_calc_found_rows", func(q ksql.QueryInterface) { q.SqlCalcFoundRows() }, "SQL_CALC_FOUND_ROWS"},
		{"all", func(q ksql.QueryInterface) { q.All() }, "ALL"},
		{"distinct_row", func(q ksql.QueryInterface) { q.DistinctRow() }, "DISTINCTROW"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := sql.NewQuery().Table("t").Columns("id")
			tt.modFn(q)
			assert.Contains(t, q.Prepare(), tt.expect)
		})
	}
}

func TestNewQuery_Select_Partitions(t *testing.T) {
	q := sql.NewQuery().Partitions("p0", "p1").Table("t").Columns("id")
	assert.Contains(t, q.Prepare(), "PARTITION")
	assert.Contains(t, q.Prepare(), "`p0`")
	assert.Contains(t, q.Prepare(), "`p1`")
}

func TestNewQuery_Select_Having(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("role", "COUNT(1) as cnt").Group("role").
		Having("cnt", ksql.Gt, 10).
		HavingIsNull("deleted").
		HavingIsNotNull("name").
		HavingIn("status", []any{1, 2}).
		HavingNotIn("type", []any{0}).
		HavingBetween("score", 60, 100).
		HavingNotBetween("age", 0, 18)

	sqlStr := q.Prepare()
	assert.Contains(t, sqlStr, "HAVING")
	assert.Contains(t, sqlStr, "`cnt` > ?")
	assert.Contains(t, sqlStr, "`deleted` IS NULL")
	assert.Contains(t, sqlStr, "`name` IS NOT NULL")
	assert.Contains(t, sqlStr, "`status` IN (?, ?)")
	assert.Contains(t, sqlStr, "`type` NOT IN (?)")
	assert.Contains(t, sqlStr, "`score` BETWEEN ? AND ?")
	assert.Contains(t, sqlStr, "`age` NOT BETWEEN ? AND ?")
}

func TestNewQuery_Select_OrHaving(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").Group("id").
		Having("a", ksql.Gt, 1).
		OrHaving(func(h ksql.HavingInterface) {
			h.Having("b", ksql.Gt, 2)
		})
	assert.Contains(t, q.Prepare(), "OR (")
	assert.Contains(t, q.Prepare(), "`b` > ?")
}

func TestNewQuery_Select_AndHaving(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").Group("id").
		Having("a", ksql.Gt, 1).
		AndHaving(func(h ksql.HavingInterface) {
			h.Having("b", ksql.Gt, 2)
		})
	assert.Contains(t, q.Prepare(), "AND (")
}

func TestNewQuery_Select_IntoVar(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").IntoVar("v_id")
	assert.Contains(t, q.Prepare(), "INTO")
	assert.Contains(t, q.Prepare(), "`v_id`")
}

func TestNewQuery_Select_Window(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").
		Window("w AS (PARTITION BY dept ORDER BY salary DESC)", "w")
	assert.Contains(t, q.Prepare(), "WINDOW")
	assert.Contains(t, q.Prepare(), "`w`")
}

func TestNewQuery_Pagination(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").Pagination(2, 15)

	assert.Contains(t, q.Prepare(), "LIMIT ?")
	assert.Contains(t, q.Prepare(), "OFFSET ?")
	assert.Equal(t, []any{15, 15}, q.Binds()) // pageSize=15, offset=(2-1)*15=15
}

func TestNewQuery_Clone(t *testing.T) {
	q := sql.NewQuery().Table("t").Columns("id").Where("id", ksql.Eq, 1)
	q2 := q.Clone()

	// Clone produces a valid non-nil query
	assert.NotNil(t, q2)

	// Clone preserves WHERE clause
	assert.Contains(t, q2.Prepare(), "WHERE `id` = ?")

	// Original still works
	assert.Contains(t, q.Prepare(), "`id`")
}

// ─────────────────────────────────────────────
// INSERT
// ─────────────────────────────────────────────

func TestInsert_Basic(t *testing.T) {
	ins := NewInsert().Table("user").Add("name", "alice").Add("age", 18)
	assert.Equal(t, "INSERT INTO `user` (`name`, `age`) VALUES (?, ?)", ins.Prepare())
	assert.Equal(t, []any{"alice", 18}, ins.Binds())
}

func TestInsert_LowPriority(t *testing.T) {
	ins := NewInsert().LowPriority().Table("user").Add("name", "alice").Add("age", 18)
	assert.Contains(t, ins.Prepare(), "LOW_PRIORITY")
}

func TestInsert_Delayed(t *testing.T) {
	ins := NewInsert().Delayed().Table("user").Add("name", "alice").Add("age", 18)
	assert.Contains(t, ins.Prepare(), "DELAYED")
}

func TestInsert_HighPriority(t *testing.T) {
	ins := NewInsert().HighPriority().Table("user").Add("name", "alice").Add("age", 18)
	assert.Contains(t, ins.Prepare(), "HIGH_PRIORITY")
}

func TestInsert_Ignore(t *testing.T) {
	ins := NewInsert().Ignore().Table("user").Add("name", "alice").Add("age", 18)
	assert.Contains(t, ins.Prepare(), "IGNORE")
}

func TestInsert_FromQuery(t *testing.T) {
	sub := NewQuery().Table("users").Columns("name", "age").Where("status", ksql.Eq, 1)
	ins := NewInsert().Table("user").Columns("name", "age").From(sub)
	assert.Contains(t, ins.Prepare(), "SELECT")
}

func TestInsert_FromTable(t *testing.T) {
	ins := NewInsert().Table("user").Columns("name", "age").FromTable("user_backup")
	assert.Contains(t, ins.Prepare(), "TABLE `user_backup`")
}

func TestInsert_OnDuplicateKeyUpdate(t *testing.T) {
	ins := NewInsert().Table("user").Add("name", "alice").Add("age", 18).
		OnDuplicateKeyUpdate("name", "alice").
		OnDuplicateKeyUpdateColumn("age", "VALUES(age)").
		OnDuplicateKeyUpdateExpress(Raw("score = score + 1"))

	assert.Contains(t, ins.Prepare(), "ON DUPLICATE KEY UPDATE")
}

func TestInsert_MultipleRows(t *testing.T) {
	ins := NewInsert().Table("user").Columns("name", "age").
		Values("alice", 18).Values("bob", 20)

	assert.Contains(t, ins.Prepare(), "VALUES (?, ?), (?, ?)")
	assert.Equal(t, []any{"alice", 18, "bob", 20}, ins.Binds())
}

func TestInsert_Partitions(t *testing.T) {
	ins := NewInsert().Partitions("p0").Table("user").Add("name", "alice").Add("age", 18)
	assert.Contains(t, ins.Prepare(), "INTO `user` `p0`")
}

func TestInsert_As(t *testing.T) {
	ins := NewInsert().Table("user").Add("name", "alice").Add("age", 18).
		As("new", "n_name")
	assert.Contains(t, ins.Prepare(), "AS")
	assert.Contains(t, ins.Prepare(), "`n_name`")
}

func TestInsert_Set(t *testing.T) {
	ins := NewInsert().Table("user").Set("name", "alice").SetColumn("age", "other_age")
	assert.Contains(t, ins.Prepare(), "SET")
	assert.Contains(t, ins.Prepare(), "`name`")
	assert.Contains(t, ins.Prepare(), "`age` = `other_age`")
}

// ─────────────────────────────────────────────
// UPDATE
// ─────────────────────────────────────────────

func TestUpdate_Basic(t *testing.T) {
	w := NewWhere().Where("id", ksql.Eq, 1)
	up := NewUpdate().Table("user").Set("name", "alice").Where(w)
	assert.Contains(t, up.Prepare(), "UPDATE `user` SET")
	assert.Contains(t, up.Prepare(), "`name` = ?")
	assert.Contains(t, up.Prepare(), "WHERE `id` = ?")
	assert.Equal(t, []any{"alice", 1}, up.Binds())
}

func TestUpdate_MultipleSets(t *testing.T) {
	w := NewWhere().Where("id", ksql.Eq, 1)
	up := NewUpdate().Table("user").Set("name", "alice").Set("age", 18).Where(w)
	assert.Contains(t, up.Prepare(), "`name` = ?")
	assert.Contains(t, up.Prepare(), "`age` = ?")
	assert.Equal(t, []any{"alice", 18, 1}, up.Binds())
}

func TestUpdate_LowPriority(t *testing.T) {
	up := NewUpdate().LowPriority().Table("user").Set("name", "alice")
	assert.Contains(t, up.Prepare(), "LOW_PRIORITY")
}

func TestUpdate_Ignore(t *testing.T) {
	up := NewUpdate().Ignore().Table("user").Set("name", "alice")
	assert.Contains(t, up.Prepare(), "IGNORE")
}

func TestUpdate_OrderByAsc(t *testing.T) {
	up := NewUpdate().Table("user").Set("name", "alice").OrderByAsc("id")
	assert.Contains(t, up.Prepare(), "ORDER BY `id` ASC")
}

func TestUpdate_OrderByDesc(t *testing.T) {
	up := NewUpdate().Table("user").Set("name", "alice").OrderByDesc("id")
	assert.Contains(t, up.Prepare(), "ORDER BY `id` DESC")
}

func TestUpdate_Limit(t *testing.T) {
	up := NewUpdate().Table("user").Set("name", "alice").Limit(5)
	assert.Contains(t, up.Prepare(), "LIMIT 5")
}

func TestUpdate_SetColumn(t *testing.T) {
	up := NewUpdate().Table("user").SetColumn("balance", "balance").SetColumn("score", "points")
	assert.Contains(t, up.Prepare(), "`balance` = `balance`")
	assert.Contains(t, up.Prepare(), "`score` = `points`")
}

func TestUpdate_SetExpress(t *testing.T) {
	up := NewUpdate().Table("user").SetExpress(Raw("score = score + ?", 1))
	assert.Contains(t, up.Prepare(), "score = score + ?")
}

func TestUpdate_IncColumn(t *testing.T) {
	up := NewUpdate().Table("user").IncColumn("score", 10).IncColumn("lives", -1)
	assert.Contains(t, up.Prepare(), "`score` = `score` + ?")
	assert.Contains(t, up.Prepare(), "`lives` = `lives` - ?")
	assert.Equal(t, []any{10, 1}, up.Binds())
}

// ─────────────────────────────────────────────
// DELETE
// ─────────────────────────────────────────────

func TestDelete_Basic(t *testing.T) {
	w := NewWhere().Where("id", ksql.Eq, 1)
	del := NewDelete().Table("user").Where(w)
	assert.Contains(t, del.Prepare(), "DELETE FROM `user`")
	assert.Contains(t, del.Prepare(), "WHERE `id` = ?")
	assert.Equal(t, []any{1}, del.Binds())
}

func TestDelete_LowPriority(t *testing.T) {
	del := NewDelete().LowPriority().Table("user").Where(NewWhere().Where("id", ksql.Eq, 1))
	assert.Contains(t, del.Prepare(), "LOW_PRIORITY")
}

func TestDelete_Quick(t *testing.T) {
	del := NewDelete().Quick().Table("user").Where(NewWhere().Where("id", ksql.Eq, 1))
	assert.Contains(t, del.Prepare(), "QUICK")
}

func TestDelete_Ignore(t *testing.T) {
	del := NewDelete().Ignore().Table("user").Where(NewWhere().Where("id", ksql.Eq, 1))
	assert.Contains(t, del.Prepare(), "IGNORE")
}

func TestDelete_As(t *testing.T) {
	del := NewDelete().Table("user").As("u").Where(NewWhere().Where("u.id", ksql.Eq, 1))
	assert.Contains(t, del.Prepare(), "AS `u`")
}

func TestDelete_Partitions(t *testing.T) {
	del := NewDelete().Partitions("p0").Table("user").Where(NewWhere().Where("id", ksql.Eq, 1))
	assert.Contains(t, del.Prepare(), "PARTITION")
	assert.Contains(t, del.Prepare(), "`p0`")
}

func TestDelete_OrderByAsc(t *testing.T) {
	del := NewDelete().Table("user").OrderByAsc("id").Where(NewWhere().Where("id", ksql.Eq, 1))
	assert.Contains(t, del.Prepare(), "ORDER BY `id` ASC")
}

func TestDelete_OrderByDesc(t *testing.T) {
	del := NewDelete().Table("user").OrderByDesc("id").Where(NewWhere().Where("id", ksql.Eq, 1))
	assert.Contains(t, del.Prepare(), "ORDER BY `id` DESC")
}

func TestDelete_Limit(t *testing.T) {
	del := NewDelete().Table("user").Limit(5).Where(NewWhere().Where("id", ksql.Eq, 1))
	assert.Contains(t, del.Prepare(), "LIMIT 5")
}

// ─────────────────────────────────────────────
// DDL
// ─────────────────────────────────────────────

func TestDropTable(t *testing.T) {
	dt := NewDropTable().Table("user")
	assert.Equal(t, "DROP TABLE `user`", dt.Prepare())
}

func TestDropTable_IfExists(t *testing.T) {
	dt := NewDropTable().Table("user").IfExists()
	assert.Equal(t, "DROP TABLE IF EXISTS `user`", dt.Prepare())
}

func TestDropTable_Temporary(t *testing.T) {
	dt := NewDropTable().Table("user").Temporary()
	assert.Equal(t, "DROP TEMPORARY TABLE `user`", dt.Prepare())
}

func TestDropTable_Restrict(t *testing.T) {
	dt := NewDropTable().Table("user").Restrict()
	assert.Contains(t, dt.Prepare(), "RESTRICT")
}

func TestDropTable_Cascade(t *testing.T) {
	dt := NewDropTable().Table("user").Cascade()
	assert.Contains(t, dt.Prepare(), "CASCADE")
}

// ─────────────────────────────────────────────
// WHERE builder
// ─────────────────────────────────────────────

func TestWhere_Empty(t *testing.T) {
	w := NewWhere()
	assert.True(t, w.Empty())
	w.Where("id", ksql.Eq, 1)
	assert.False(t, w.Empty())
}

func TestWhere_IsNull_IsNotNull(t *testing.T) {
	w := NewWhere()
	w.IsNull("a").IsNotNull("b")
	sqlStr := buildWhere(w)
	assert.Contains(t, sqlStr, "`a` IS NULL")
	assert.Contains(t, sqlStr, "`b` IS NOT NULL")
}

func TestWhere_Clone_NotNil(t *testing.T) {
	w := NewWhere()
	w.Where("id", ksql.Eq, 1)
	w2 := w.Clone()

	assert.NotNil(t, w2)
	assert.False(t, w2.Empty())
}

func TestWhere_Express(t *testing.T) {
	w := NewWhere().Express(Raw("a = ? AND b = ?", 1, 2))
	assert.Contains(t, buildWhere(w), "a = ? AND b = ?")
	assert.Equal(t, []any{1, 2}, w.Binds())
}

func TestWhere_InBy_NotInBy(t *testing.T) {
	sub := NewQuery().Table("t").Columns("id").Where("s", ksql.Eq, 1)
	w := NewWhere().InBy("a", sub).NotInBy("b", sub)

	sqlStr := buildWhere(w)
	assert.Contains(t, sqlStr, "`a` IN (SELECT")
	assert.Contains(t, sqlStr, "`b` NOT IN (SELECT")
}

func TestWhere_Between_NotBetween(t *testing.T) {
	w := NewWhere().Between("a", 1, 10).NotBetween("b", 100, 200)
	sqlStr := buildWhere(w)
	assert.Contains(t, sqlStr, "`a` BETWEEN ? AND ?")
	assert.Contains(t, sqlStr, "`b` NOT BETWEEN ? AND ?")
	assert.Equal(t, []any{1, 10, 100, 200}, w.Binds())
}

func TestWhere_NestedAndOr(t *testing.T) {
	w := NewWhere().Where("a", ksql.Eq, 1).
		AndWhere(func(o ksql.WhereInterface) {
			o.Where("b", ksql.Gt, 0).IsNull("c")
		}).
		OrWhere(func(o ksql.WhereInterface) {
			o.Where("d", ksql.Lt, 0)
		})

	sqlStr := buildWhere(w)
	assert.Contains(t, sqlStr, "AND (`b` > ? AND `c` IS NULL)")
	assert.Contains(t, sqlStr, "OR (`d` < ?)")
	assert.Equal(t, []any{1, 0, 0}, w.Binds())
}

func TestWhere_PanicOnUnsupportedOp(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unsupported op")
		}
	}()
	NewWhere().Where("x", "INVALID", 1)
}

// ─────────────────────────────────────────────
// HAVING builder
// ─────────────────────────────────────────────

func TestHaving_Basic(t *testing.T) {
	h := sql.NewHaving()
	h.Having("cnt", ksql.Gt, 10).IsNull("a").IsNotNull("b")
	sqlStr := buildHaving(h)
	assert.Contains(t, sqlStr, "`cnt` > ?")
	assert.Contains(t, sqlStr, "`a` IS NULL")
	assert.Contains(t, sqlStr, "`b` IS NOT NULL")
}

func TestHaving_In(t *testing.T) {
	h := sql.NewHaving().In("a", []any{1, 2}).NotIn("b", []any{3})
	sqlStr := buildHaving(h)
	assert.Contains(t, sqlStr, "`a` IN (?, ?)")
	assert.Contains(t, sqlStr, "`b` NOT IN (?)")
}

func TestHaving_Between(t *testing.T) {
	h := sql.NewHaving().Between("a", 1, 10).NotBetween("b", 100, 200)
	assert.Contains(t, buildHaving(h), "BETWEEN")
}

func TestHaving_Clone(t *testing.T) {
	h := sql.NewHaving().Having("a", ksql.Gt, 1)
	h2 := h.Clone()
	h2.Having("b", ksql.Gt, 2)
	assert.NotContains(t, buildHaving(h), "`b`")
	assert.Contains(t, buildHaving(h2), "`b`")
}

// ─────────────────────────────────────────────
// Raw SQL
// ─────────────────────────────────────────────

func TestRaw_Statement(t *testing.T) {
	r := Raw("SELECT * FROM user WHERE id = ?", 1)
	assert.Equal(t, "SELECT * FROM user WHERE id = ?", r.Statement())
	assert.Equal(t, []any{1}, r.Binds())
	assert.False(t, r.IsExec())  // SELECT is a query
	assert.True(t, Raw("INSERT INTO t VALUES(?)", 1).IsExec())
	assert.True(t, Raw("UPDATE t SET x = ?", 1).IsExec())
	assert.True(t, Raw("DELETE FROM t WHERE x = ?", 1).IsExec())
	assert.True(t, Raw("DROP TABLE t").IsExec())
	assert.True(t, Raw("CREATE TABLE t(x INT)").IsExec())
	assert.True(t, Raw("ALTER TABLE t ADD x INT").IsExec())
}

func TestRaw_TrimsWhitespace(t *testing.T) {
	r := Raw("  SELECT 1  \n")
	assert.Equal(t, "SELECT 1", r.Statement())
}

// ─────────────────────────────────────────────
// Data
// ─────────────────────────────────────────────

func TestData_SetGet(t *testing.T) {
	d := NewData()
	d.Set("name", "alice").Set("age", 18)

	assert.Equal(t, "alice", d.Get("name"))
	assert.Equal(t, 18, d.Get("age"))
	assert.Equal(t, []string{"name", "age"}, d.Keys())
	assert.False(t, d.Empty())
}

func TestData_Overwrite(t *testing.T) {
	d := NewData()
	d.Set("name", "alice")
	d.Set("name", "bob")
	assert.Equal(t, "bob", d.Get("name"))
	assert.Equal(t, []string{"name"}, d.Keys())
}

func TestData_Range(t *testing.T) {
	d := NewData()
	d.Set("a", 1).Set("b", 2)

	var keys []string
	var vals []any
	d.Range(func(k string, v any) {
		keys = append(keys, k)
		vals = append(vals, v)
	})
	assert.Equal(t, []string{"a", "b"}, keys)
}

func TestData_Changed(t *testing.T) {
	d := NewData()
	d.Set("name", "alice")

	assert.True(t, d.Changed("name", "bob"))
	assert.False(t, d.Changed("name", "alice"))
	assert.True(t, d.Changed("missing", "x"))
}

func TestData_From(t *testing.T) {
	src := NewData().Set("a", 1).Set("b", 2)
	dst := NewData()
	dst.From(src)
	assert.Equal(t, 1, dst.Get("a"))
	assert.Equal(t, 2, dst.Get("b"))
	assert.Equal(t, []string{"a", "b"}, dst.Keys())
}

// ─────────────────────────────────────────────
// PageInfo
// ─────────────────────────────────────────────

func TestPageInfo(t *testing.T) {
	list := []*test_user{{Id: 1}, {Id: 2}}
	p := NewPageInfo(list)
	assert.Equal(t, 2, len(p.List()))

	p.Set(9, 2)
	assert.Equal(t, uint64(9), p.TotalCount())
	assert.Equal(t, uint64(5), p.TotalPage()) // 9/2 = 4.5 → 5 pages
}

func TestPageInfo_ExactDivision(t *testing.T) {
	p := NewPageInfo([]*test_user{})
	p.Set(10, 5)
	assert.Equal(t, uint64(10), p.TotalCount())
	assert.Equal(t, uint64(2), p.TotalPage()) // 10/5 = 2 exactly
}

// ─────────────────────────────────────────────
// Map
// ─────────────────────────────────────────────

func TestMap(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1).Set("b", 2)

	assert.Equal(t, 1, m.Get("a"))
	assert.True(t, m.Has("a"))
	assert.False(t, m.Has("c"))
	assert.Equal(t, []string{"a", "b"}, m.Keys())
	assert.Equal(t, []int{1, 2}, m.Values())
	assert.Equal(t, 2, m.GetBy(1))

	var val int
	assert.Equal(t, val, m.GetBy(99)) // out of bounds returns zero value
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func TestToList(t *testing.T) {
	result := ToList([]int{1, 2, 3})
	assert.Equal(t, []any{1, 2, 3}, result)
}
