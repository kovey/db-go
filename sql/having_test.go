package sql

import "testing"

func TestHavingPrepare(t *testing.T) {
	where := NewHaving()
	where.Eq("age", 18).Neq("name", "kovey").Like("nickname", "%golang%").Between("id", 1, 100).Gt("balance", 1000).Ge("count", 10).Lt("sum", 15).Le("people", 100)
	where.In("sex", []any{0, 1, 2}).NotIn("lang", []any{"php", "java", "ruby", "rust"}).IsNull("content").IsNotNull("title").Statement("last_id > 0")

	expected := "HAVING (`age` = ? AND `name` <> ? AND `nickname` LIKE ? AND `id` BETWEEN ? AND ? AND `balance` > ? AND `count` >= ? AND `sum` < ? AND `people` <= ? AND `sex` IN(?,?,?) AND `lang` NOT IN(?,?,?,?) AND `content` IS NULL AND `title` IS NOT NULL AND last_id > 0)"
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
	t.Logf("or having: %s", orData)
	expected = "HAVING (`age` = ? OR `name` <> ? OR `nickname` LIKE ? OR `id` BETWEEN ? AND ? OR `balance` > ? OR `count` >= ? OR `sum` < ? OR `people` <= ? OR `sex` IN(?,?,?) OR `lang` NOT IN(?,?,?,?) OR `content` IS NULL OR `title` IS NOT NULL OR last_id > 0)"
	t.Logf("or having: %s", expected)
	if expected != orData {
		t.Fatal("or where is error")
	}

	t.Logf("sql: %s", where)
}
