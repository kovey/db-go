package mysql

import (
	"database/sql"
	"strings"

	"github.com/kovey/db-go/ksql/schema"
	ksql "github.com/kovey/db-go/v3"
)

type Index struct {
	*base
	Table         string
	Non_unique    int
	Key_name      string
	Seq_in_index  int
	Column_name   string
	Collation     string
	Cardinality   int
	Sub_part      sql.NullString
	Packed        sql.NullString
	Null          string
	Index_type    string
	Index_comment string
	Comm          string
	Visible       string
	Expression    sql.NullString
}

func (i *Index) Name() string {
	return i.Key_name
}

func (i *Index) NonUnique() int {
	return i.Non_unique
}

func (i *Index) Seq() int {
	return i.Seq_in_index
}

func (i *Index) Column() string {
	return i.Column_name
}

func (i *Index) Comment() string {
	return i.Index_comment
}

func (i *Index) IndexComment() string {
	return i.Index_comment
}

func (i *Index) Values() []any {
	return []any{&i.Table, &i.Non_unique, &i.Key_name, &i.Seq_in_index, &i.Column_name, &i.Collation, &i.Cardinality, &i.Sub_part, &i.Packed, &i.Null, &i.Index_type, &i.Index_comment, &i.Comm, &i.Visible, &i.Expression}
}

func (i *Index) Clone() ksql.RowInterface {
	return &Index{base: &base{empty: true}}
}

func (i *Index) Columns() []string {
	return []string{"Non_unique", "Key_name", "Seq_in_index", "Column_name", "Index_comment", "Comment"}
}

func (i *Index) Type() ksql.IndexType {
	if strings.ToUpper(i.Key_name) == "PRIMARY" {
		return ksql.Index_Type_Primary
	}

	if i.NonUnique() == 1 {
		return ksql.Index_Type_Normal
	}

	return ksql.Index_Type_Unique
}

func (i *Index) HasChanged(other schema.IndexMetaInterface) bool {
	return i.Name() != other.Name() || i.NonUnique() != other.NonUnique() || i.Seq() != other.Seq() || i.Column() != other.Column()
}

type TableIndex struct {
	name    string
	indexes []schema.IndexMetaInterface
	t       ksql.IndexType
}

func (t *TableIndex) Add(index schema.IndexMetaInterface) schema.IndexInfoInterface {
	if index.Name() != t.name {
		return t
	}

	t.t = index.Type()
	t.indexes = append(t.indexes, index)
	return t
}

func (t *TableIndex) Name() string {
	return t.name
}

func (t *TableIndex) Metas() []schema.IndexMetaInterface {
	return t.indexes
}

func (t *TableIndex) Type() ksql.IndexType {
	return t.t
}

func (t *TableIndex) HasChanged(other schema.IndexInfoInterface) bool {
	return t.Name() != other.Name() || t.Type() != other.Type() || len(t.Metas()) != len(other.Metas())
}

func (t *TableIndex) Columns() []string {
	var columns = make([]string, len(t.indexes))
	for i, index := range t.indexes {
		columns[i] = index.Column()
	}

	return columns
}
