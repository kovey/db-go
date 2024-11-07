package mysql

import (
	"strings"

	"github.com/kovey/db-go/ksql/schema"
	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

type Table struct {
	*base
	TABLE_NAME      string
	ENGINE          string
	TABLE_COLLATION string
	TABLE_COMMENT   string
	columns         []schema.ColumnInfoInterface
	indexes         *db.Map[string, schema.IndexInfoInterface]
}

func NewTable(conn ksql.ConnectionInterface) *Table {
	return &Table{base: &base{conn: conn, empty: true}, indexes: db.NewMap[string, schema.IndexInfoInterface]()}
}

func (t *Table) SetColumns(columns []schema.ColumnInfoInterface) schema.TableInfoInterface {
	t.columns = columns
	return t
}

func (t *Table) SetIndexes(indexes []schema.IndexMetaInterface) schema.TableInfoInterface {
	for _, index := range indexes {
		if !t.indexes.Has(index.Name()) {
			t.indexes.Set(index.Name(), &TableIndex{name: index.Name()})
		}

		t.indexes.Get(index.Name()).Add(index)
	}

	return t
}

func (t *Table) Indexes() []schema.IndexInfoInterface {
	return t.indexes.Values()
}

func (t *Table) Name() string {
	return t.TABLE_NAME
}

func (t *Table) Engine() string {
	return t.ENGINE
}

func (t *Table) Charset() string {
	return strings.Split(t.TABLE_COLLATION, "_")[0]
}

func (t *Table) Collation() string {
	return t.TABLE_COLLATION
}

func (t *Table) Comment() string {
	return t.TABLE_COMMENT
}

func (t *Table) Fields() []schema.ColumnInfoInterface {
	return t.columns
}

func (t *Table) HasChanged(other schema.TableInfoInterface) bool {
	if t.Name() != other.Name() || t.Engine() != other.Engine() || t.Charset() != other.Charset() || t.Collation() != other.Collation() || t.Comment() != other.Comment() {
		return true
	}

	if len(t.Fields()) != len(other.Fields()) || len(t.Indexes()) != len(other.Indexes()) {
		return true
	}

	for _, column := range t.columns {
		if !other.HasColumn(column) {
			return true
		}

		if other.GetColumn(column.Name()).HasChanged(column) {
			return true
		}
	}

	indexes := other.Indexes()
	for _, index := range indexes {
		if !t.HasIndex(index) {
			return true
		}

		if !t.indexes.Has(index.Name()) {
			return true
		}

		if t.indexes.Get(index.Name()).HasChanged(index) {
			return true
		}
	}

	for _, index := range t.Indexes() {
		if !other.HasIndex(index) {
			return true
		}

		oIndex := other.GetIndex(index.Name())
		if oIndex == nil {
			return true
		}

		if oIndex.HasChanged(index) {
			return true
		}
	}

	return false
}

func (t *Table) HasColumn(column schema.ColumnInfoInterface) bool {
	for _, col := range t.columns {
		if col.Name() == column.Name() {
			return true
		}
	}

	return false
}

func (t *Table) GetColumn(column string) schema.ColumnInfoInterface {
	for _, col := range t.columns {
		if col.Name() == column {
			return col
		}
	}
	return nil
}

func (t *Table) HasIndex(index schema.IndexInfoInterface) bool {
	return t.indexes.Has(index.Name())
}

func (t *Table) Values() []any {
	return []any{&t.TABLE_NAME, &t.ENGINE, &t.TABLE_COLLATION, &t.TABLE_COMMENT}
}

func (t *Table) Clone() ksql.RowInterface {
	return NewTable(nil)
}

func (t *Table) Columns() []string {
	return []string{"TABLE_NAME", "ENGINE", "TABLE_COLLATION", "TABLE_COMMENT"}
}

func (t *Table) GetIndex(index string) schema.IndexInfoInterface {
	return t.indexes.Get(index)
}

func (t *Table) CheckChanges(other schema.TableInfoInterface) schema.ChangedInterface {
	var changed = &Changed{column: &ColumnChanged{}, index: &IndexChanged{}}
	for _, column := range t.columns {
		if !other.HasColumn(column) {
			changed.column.adds = append(changed.column.adds, column)
			continue
		}

		old := other.GetColumn(column.Name())
		if column.HasChanged(old) {
			changed.column.changes = append(changed.column.changes, &ColumnMetaChanged{o: old, n: column})
		}
	}

	t.indexes.Range(func(key string, val schema.IndexInfoInterface) {
		if !other.HasIndex(val) {
			changed.index.adds = append(changed.index.adds, val)
			return
		}

		old := other.GetIndex(val.Name())
		if val.HasChanged(old) {
			changed.index.deletes = append(changed.index.deletes, old)
			changed.index.adds = append(changed.index.adds, val)
		}
	})

	for _, index := range other.Indexes() {
		if !t.HasIndex(index) {
			changed.index.deletes = append(changed.index.deletes, index)
		}
	}

	return changed
}
