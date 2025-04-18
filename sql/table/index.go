package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type Index struct {
	name       string
	alg        ksql.IndexAlg
	columns    *IndexColumns
	option     *IndexOption
	constraint string
	symbol     string
	unique     string
	primary    string
	foreign    string
	reference  *ColumnReference
	opChain    *operator.Chain
	typ        ksql.IndexType
	subType    ksql.IndexSubType
}

func NewIndex(name string) *Index {
	i := &Index{name: name, opChain: operator.NewChain(), typ: ksql.Index_Type_Normal, subType: ksql.Index_Sub_Type_Index, columns: &IndexColumns{}}
	i.opChain.Append(i._keyword, i._key, i._name, i._keyType, i.columns.Build, i._indexOption, i._reference)
	return i
}

func (i *Index) _keyword(builder *strings.Builder) {
	if i.constraint != "" {
		builder.WriteString(" CONSTRAINT")
		operator.BuildPureString(i.symbol, builder)
		return
	}

	switch i.typ {
	case ksql.Index_Type_FullText, ksql.Index_Type_Spatial:
		builder.WriteString(" ")
		builder.WriteString(i.typ.String())
	}
}

func (i *Index) _key(builder *strings.Builder) {
	switch i.typ {
	case ksql.Index_Type_FullText, ksql.Index_Type_Spatial, ksql.Index_Type_Normal:
		builder.WriteString(" ")
		builder.WriteString(string(i.subType))
	default:
		if i.primary != "" {
			builder.WriteString(" PRIMARY ")
			builder.WriteString(string(ksql.Index_Sub_Type_Key))
			return
		}

		if i.unique != "" {
			builder.WriteString(" UNIQUE ")
			builder.WriteString(string(i.subType))
			return
		}

		if i.foreign != "" {
			builder.WriteString(" FOREIGN ")
			builder.WriteString(string(ksql.Index_Sub_Type_Key))
		}
	}
}

func (i *Index) _name(builder *strings.Builder) {
	if i.primary != "" {
		return
	}

	operator.BuildBacktickString(i.name, builder)
}

func (i *Index) _keyType(builder *strings.Builder) {
	if i.alg == "" || i.foreign != "" {
		return
	}

	switch i.typ {
	case ksql.Index_Type_FullText, ksql.Index_Type_Spatial:
		return
	}

	builder.WriteString(" USING")
	operator.BuildQuoteString(string(i.alg), builder)
}

func (i *Index) _indexOption(builder *strings.Builder) {
	if i.option == nil || i.foreign != "" {
		return
	}

	i.option.Build(builder)
}

func (i *Index) _reference(builder *strings.Builder) {
	if i.foreign == "" || i.reference == nil {
		return
	}

	i.reference.Build(builder)
}

func (i *Index) Build(builder *strings.Builder) {
	i.opChain.Call(builder)
}

func (i *Index) Type(typ ksql.IndexType) ksql.TableIndexInterface {
	switch typ {
	case ksql.Index_Type_FullText, ksql.Index_Type_Spatial:
		i.typ = typ
	}

	return i
}

func (i *Index) SubType(subType ksql.IndexSubType) ksql.TableIndexInterface {
	i.subType = subType
	return i
}

func (i *Index) Algorithm(algorithm ksql.IndexAlg) ksql.TableIndexInterface {
	i.alg = algorithm
	return i
}

func (i *Index) Column(name string, length int, order ksql.Order) ksql.TableIndexInterface {
	i.columns.Append(&IndexColumn{Name: name, Length: length, Type: Index_Column_Type_Name, Order: order})
	return i
}

func (i *Index) Express(express string, order ksql.Order) ksql.TableIndexInterface {
	i.columns.Append(&IndexColumn{Name: express, Type: Index_Column_Type_Expr, Order: order})
	return i
}

func (i *Index) Columns(columns ...string) ksql.TableIndexInterface {
	for _, column := range columns {
		i.columns.Append(&IndexColumn{Name: column, Type: Index_Column_Type_Pure_Name, Order: ksql.Order_None})
	}

	return i
}

func (i *Index) BlockSize(size string) ksql.TableIndexInterface {
	i.option.BlockSize(size)
	return i
}

func (i *Index) WithParser(parserName string) ksql.TableIndexInterface {
	i.option.WithParser(parserName)
	return i
}

func (i *Index) Comment(comment string) ksql.TableIndexInterface {
	i.option.Comment(comment)
	return i
}

func (i *Index) Visible() ksql.TableIndexInterface {
	i.option.Visible()
	return i
}

func (i *Index) Invisible() ksql.TableIndexInterface {
	i.option.Invisible()
	return i
}

func (i *Index) EngineAttribute(attr string) ksql.TableIndexInterface {
	i.option.EngineAttribute(attr)
	return i
}

func (i *Index) SecondaryEngineAttribute(attr string) ksql.TableIndexInterface {
	i.option.SecondaryEngineAttribute(attr)
	return i
}

func (i *Index) Constraint(symbol string) ksql.TableIndexInterface {
	i.constraint = "CONSTRAINT"
	i.symbol = symbol
	return i
}

func (i *Index) Primary() ksql.TableIndexInterface {
	i.primary = "PRIMARY"
	i.typ = ksql.Index_Type_Primary
	return i
}

func (i *Index) Unique() ksql.TableIndexInterface {
	i.unique = "UNIQUE"
	i.typ = ksql.Index_Type_Unique
	return i
}

func (i *Index) Foreign() ksql.TableIndexInterface {
	i.foreign = "FOREIGN"
	i.typ = ksql.Index_Type_Foreign
	return i
}

func (i *Index) Reference(table string) ksql.ColumnReferenceInterface {
	i.reference = NewColumnReference(table)
	return i.reference
}
