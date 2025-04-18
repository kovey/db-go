package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type AddIndex struct {
	index *Index
}

func NewAddIndex(index string) *AddIndex {
	return &AddIndex{index: NewIndex(index)}
}

func (a *AddIndex) Index() ksql.TableIndexInterface {
	return a.index
}

func (a *AddIndex) Build(builder *strings.Builder) {
	builder.WriteString("ADD")
	a.index.Build(builder)
}
