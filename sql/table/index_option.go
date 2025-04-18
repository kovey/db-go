package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type IndexOption struct {
	blockSize     string
	alg           ksql.IndexAlg
	parserName    string
	comment       string
	visible       string
	engineAttr    string
	secEngineAttr string
	opChain       *operator.Chain
}

func NewIndexOption() *IndexOption {
	i := &IndexOption{opChain: operator.NewChain()}
	i.opChain.Append(i._blockSize, i._alg, i._other)
	return i
}

func (i *IndexOption) _blockSize(builder *strings.Builder) {
	if i.blockSize == "" {
		return
	}

	builder.WriteString(" KEY_BLOCK_SIZE =")
	operator.BuildPureString(i.blockSize, builder)
}

func (i *IndexOption) _alg(builder *strings.Builder) {
	if i.alg == "" {
		return
	}

	builder.WriteString(" USING ")
	builder.WriteString(string(i.alg))
}

func (i *IndexOption) _other(builder *strings.Builder) {
	if i.parserName != "" {
		builder.WriteString(" WITH PARSER")
		operator.BuildPureString(i.parserName, builder)
	}

	if i.comment != "" {
		builder.WriteString(" COMMENT")
		operator.BuildQuoteString(i.comment, builder)
	}

	operator.BuildPureString(i.visible, builder)

	if i.engineAttr != "" {
		builder.WriteString(" ENGINE_ATTRIBUTE = ")
		operator.Quote(i.engineAttr, builder)
	}

	if i.secEngineAttr != "" {
		builder.WriteString(" SECONDARY_ENGINE_ATTRIBUTE = ")
		operator.Quote(i.secEngineAttr, builder)
	}
}

func (i *IndexOption) Build(builder *strings.Builder) {
	i.opChain.Call(builder)
}

func (i *IndexOption) BlockSize(size string) *IndexOption {
	i.blockSize = size
	return i
}

func (i *IndexOption) Algorithm(alg ksql.IndexAlg) *IndexOption {
	i.alg = alg
	return i
}

func (i *IndexOption) WithParser(parserName string) *IndexOption {
	i.parserName = parserName
	return i
}

func (i *IndexOption) Comment(comment string) *IndexOption {
	i.comment = comment
	return i
}

func (i *IndexOption) Visible() *IndexOption {
	i.visible = "VISIBLE"
	return i
}

func (i *IndexOption) Invisible() *IndexOption {
	i.visible = "INVISIBLE"
	return i
}

func (i *IndexOption) EngineAttribute(attr string) *IndexOption {
	i.engineAttr = attr
	return i
}

func (i *IndexOption) SecondaryEngineAttribute(attr string) *IndexOption {
	i.secEngineAttr = attr
	return i
}
