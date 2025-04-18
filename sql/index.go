package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/table"
)

type Index struct {
	*base
	typ        ksql.IndexType
	index      string
	alg        ksql.IndexAlg
	table      string
	columns    *table.IndexColumns
	option     *table.IndexOption
	algOption  ksql.IndexAlgOption
	lockOption ksql.IndexLockOption
}

func NewIndex() *Index {
	i := &Index{base: newBase(), columns: &table.IndexColumns{}, option: table.NewIndexOption()}
	i.opChain.Append(keywordCreate, i._type, i._index, i._alg, i._table, i._indexOption, i._algorithmOption, i._lockOption)
	return i
}

func (i *Index) _type(builder *strings.Builder) {
	if i.typ == 0x0 {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(i.typ.String())
}

func (i *Index) _index(builder *strings.Builder) {
	builder.WriteString(" INDEX ")
	Backtick(i.index, builder)
}

func (i *Index) _alg(builder *strings.Builder) {
	if i.alg == "" {
		return
	}

	builder.WriteString(" USING ")
	builder.WriteString(string(i.alg))
}

func (i *Index) _table(builder *strings.Builder) {
	builder.WriteString(" ON ")
	Backtick(i.table, builder)
	i.columns.Build(builder)
}

func (i *Index) _indexOption(builder *strings.Builder) {
	i.option.Build(builder)
}

func (i *Index) _algorithmOption(builder *strings.Builder) {
	if i.algOption == "" {
		return
	}

	builder.WriteString(" ALGORITHM = ")
	builder.WriteString(string(i.algOption))
}

func (i *Index) _lockOption(builder *strings.Builder) {
	if i.lockOption == "" {
		return
	}

	builder.WriteString(" LOCK = ")
	builder.WriteString(string(i.lockOption))
}

func (i *Index) Type(typ ksql.IndexType) ksql.IndexInterface {
	if typ != ksql.Index_Type_FullText && typ != ksql.Index_Type_Unique && typ != ksql.Index_Type_Spatial {
		return i
	}

	i.typ = typ
	return i
}

func (i *Index) Index(name string) ksql.IndexInterface {
	i.index = name
	return i
}

func (i *Index) Algorithm(alg ksql.IndexAlg) ksql.IndexInterface {
	i.alg = alg
	return i
}

func (i *Index) On(table string) ksql.IndexInterface {
	i.table = table
	return i
}

func (i *Index) Column(name string, length int, order ksql.Order) ksql.IndexInterface {
	i.columns.Append(&table.IndexColumn{Name: name, Length: length, Type: table.Index_Column_Type_Name, Order: order})
	return i
}

func (i *Index) Express(express string, order ksql.Order) ksql.IndexInterface {
	i.columns.Append(&table.IndexColumn{Name: express, Type: table.Index_Column_Type_Expr, Order: order})
	return i
}

func (i *Index) BlockSize(size string) ksql.IndexInterface {
	i.option.BlockSize(size)
	return i
}

func (i *Index) WithParser(parserName string) ksql.IndexInterface {
	i.option.WithParser(parserName)
	return i
}

func (i *Index) Comment(comment string) ksql.IndexInterface {
	i.option.Comment(comment)
	return i
}

func (i *Index) Visible() ksql.IndexInterface {
	i.option.Visible()
	return i
}

func (i *Index) Invisible() ksql.IndexInterface {
	i.option.Invisible()
	return i
}

func (i *Index) EngineAttribute(attr string) ksql.IndexInterface {
	i.option.EngineAttribute(attr)
	return i
}

func (i *Index) SecondaryEngineAttribute(attr string) ksql.IndexInterface {
	i.option.SecondaryEngineAttribute(attr)
	return i
}

func (i *Index) AlgorithmOption(option ksql.IndexAlgOption) ksql.IndexInterface {
	i.algOption = option
	return i
}

func (i *Index) LockOption(option ksql.IndexLockOption) ksql.IndexInterface {
	i.lockOption = option
	return i
}
