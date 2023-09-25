package sql

import (
	"fmt"
	"strings"
)

const (
	batchFormat      = "INSERT INTO %s (%s) VALUES %s"
	batchValueFormat = "(%s)"
)

type Batch struct {
	ins          []*Insert
	table        string
	argsCount    int
	args         []any
	placeholders []string
}

func NewBatch(table string) *Batch {
	return &Batch{table: table, ins: make([]*Insert, 0), argsCount: 0}
}

func (b *Batch) Add(insert *Insert) *Batch {
	b.ins = append(b.ins, insert)
	b.argsCount += len(insert.data)
	return b
}

func (b *Batch) Args() []any {
	return b.args
}

func (b *Batch) getFields() []string {
	if len(b.ins) < 1 {
		return []string{}
	}

	b.args = make([]any, b.argsCount)
	b.placeholders = make([]string, len(b.ins))

	first := b.ins[0]
	fields := first.getFields()

	index := 0
	bi := 0
	for _, in := range b.ins {
		for _, field := range first.fields {
			b.args[index] = in.data[field]
			index++
		}

		b.placeholders[bi] = b.formatValue(in.getPlaceholder())
		bi++
	}

	return fields
}

func (b *Batch) formatValue(placeholders []string) string {
	return fmt.Sprintf(batchValueFormat, strings.Join(placeholders, comma))
}

func (b *Batch) Prepare() string {
	return fmt.Sprintf(batchFormat, formatValue(b.table), strings.Join(b.getFields(), comma), strings.Join(b.placeholders, comma))
}

func (b *Batch) String() string {
	return String(b)
}

func (b *Batch) Inserts() []*Insert {
	return b.ins
}
