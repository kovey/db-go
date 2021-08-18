package sql

import "testing"

func TestFormatString(t *testing.T) {
	str := formatString("test")
	t.Logf("str: %s", str)
	if str != "'test'" {
		t.Fatal("formatString test fail")
	}
}

func TestFormaValue(t *testing.T) {
	str := formatValue("name")
	t.Logf("str: %s", str)

	if str != "`name`" {
		t.Fatal("formatValue test fail")
	}

	str = formatValue("test.name")
	t.Logf("str: %s", str)
	if str != "`test`.`name`" {
		t.Fatal("formatValue test fail")
	}
}

func TestFormatOrder(t *testing.T) {
	str := formatOrder("create_time DESC")
	t.Logf("str: %s", str)
	if str != "`create_time` DESC" {
		t.Fatal("formatOrder test fail")
	}

	str = formatOrder("test.create_time ASC")
	t.Logf("str: %s", str)

	if str != "`test`.`create_time` ASC" {
		t.Fatal("formatOrder test fail")
	}
}
