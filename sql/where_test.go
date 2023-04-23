package sql

import "testing"

func TestPrepare(t *testing.T) {
	where := NewWhere()
	where.Eq("age", 18)
	where.Neq("name", "kovey")
	where.Like("nickname", "%golang%")
	where.Between("id", 1, 100)
	where.Gt("balance", 1000)
	where.Ge("count", 10)
	where.Lt("sum", 15)
	where.Le("people", 100)
	where.In("sex", []any{0, 1, 2})
	where.NotIn("lang", []any{"php", "java", "ruby", "rust"})
	where.IsNull("content")
	where.IsNotNull("title")
	where.Statement("last_id > 0")

	expected := "WHERE (`age` = ? AND `name` <> ? AND `nickname` LIKE ? AND `id` BETWEEN ? AND ? AND `balance` > ? AND `count` >= ? AND `sum` < ? AND `people` <= ? AND `sex` IN(?,?,?) AND `lang` NOT IN(?,?,?,?) AND `content` IS NULL AND `title` IS NOT NULL AND last_id > 0)"
	realData := where.Prepare()
	if expected != realData {
		t.Errorf("expected: %s realData: %s", expected, realData)
	}

	args := where.Args()
	t.Logf("args: %v", args)
	if len(args) != 16 {
		t.Fatal("args len is error")
	}

	orData := where.OrPrepare()
	t.Logf("or where: %s", orData)
	expected = "WHERE (`age` = ? OR `name` <> ? OR `nickname` LIKE ? OR `id` BETWEEN ? AND ? OR `balance` > ? OR `count` >= ? OR `sum` < ? OR `people` <= ? OR `sex` IN(?,?,?) OR `lang` NOT IN(?,?,?,?) OR `content` IS NULL OR `title` IS NOT NULL OR last_id > 0)"
	t.Logf("or where: %s", expected)
	if expected != orData {
		t.Fatal("or where is error")
	}

	t.Logf("sql: %s", where)
}
