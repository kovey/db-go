package sql

import (
	"strings"

	"github.com/kovey/db-go/v3/sql/operator"
)

type drop struct {
	*base
	name     string
	ifExists bool
	keyword  string
	isMulti  bool
	names    []string
	schema   string
}

func newDrop(keyword string) *drop {
	d := &drop{base: newBase(), keyword: keyword}
	d.opChain.Append(d._build)
	return d
}

func (d *drop) _build(builder *strings.Builder) {
	builder.WriteString("DROP")
	operator.BuildPureString(d.keyword, builder)
	if d.ifExists {
		builder.WriteString(" IF EXISTS")
	}

	if !d.isMulti {
		d._buildName(d.name, builder)
		return
	}

	for index, name := range d.names {
		if index > 0 {
			builder.WriteString(",")
		}

		d._buildName(name, builder)
	}
}

func (d *drop) _buildName(name string, builder *strings.Builder) {
	if d.schema == "" {
		operator.BuildColumnString(name, builder)
		return
	}

	operator.BuildColumnString(d.schema, builder)
	builder.WriteString(".")
	operator.Column(name, builder)
}
