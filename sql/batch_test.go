package sql

import "testing"

func TestBatchPrepare(t *testing.T) {
	batch := NewBatch("product")
	in1 := NewInsert("product")
	in1.Set("name", "kovey").Set("sex", 1)

	in2 := NewInsert("product")
	in2.Set("name", "koveys").Set("sex", 2)

	batch.Add(in1).Add(in2)

	t.Logf("sql: %s", batch)
}
