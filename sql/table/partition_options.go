package table

import (
	"strconv"
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type PartitionOptionsHash struct {
	expr   string
	linear string
}

func (p *PartitionOptionsHash) Linear() ksql.PartitionOptionsHashInterface {
	p.linear = "LINEAR"
	return p
}

func (p *PartitionOptionsHash) Build(builder *strings.Builder) {
	if p.linear != "" {
		builder.WriteString(" ")
		builder.WriteString(p.linear)
	}

	builder.WriteString(" HASH(")
	builder.WriteString(p.expr)
	builder.WriteString(")")
}

type PartitionOptionsKey struct {
	columns   []string
	algorithm string
	linear    string
}

func (p *PartitionOptionsKey) Algorithm(alg string) ksql.PartitionOptionsKeyInterface {
	p.algorithm = alg
	return p
}

func (p *PartitionOptionsKey) Linear() ksql.PartitionOptionsKeyInterface {
	p.linear = "LINEAR"
	return p
}

func (p *PartitionOptionsKey) Build(builder *strings.Builder) {
	if p.linear != "" {
		builder.WriteString(" ")
		builder.WriteString(p.linear)
	}

	if p.algorithm != "" {
		builder.WriteString(" ALGORITHM = ")
		builder.WriteString(p.algorithm)
	}

	builder.WriteString(" (")
	for index, column := range p.columns {
		if index > 0 {
			builder.WriteString(", ")
		}

		operator.Column(column, builder)
	}

	builder.WriteString(")")
}

type PartitionOptionsRange struct {
	columns []string
	expr    string
}

func (p *PartitionOptionsRange) Expr(expr string) ksql.PartitionOptionsRangeInterface {
	if p.columns != nil {
		return p
	}

	p.expr = expr
	return p
}

func (p *PartitionOptionsRange) Columns(columns []string) ksql.PartitionOptionsRangeInterface {
	if p.expr != "" {
		return p
	}

	p.columns = columns
	return p
}

func (p *PartitionOptionsRange) Build(builder *strings.Builder) {
	builder.WriteString(" RANGE")
	if p.expr != "" {
		builder.WriteString("(")
		builder.WriteString(p.expr)
		builder.WriteString(")")
		return
	}

	builder.WriteString(" COLUMNS(")
	for index, column := range p.columns {
		if index > 0 {
			builder.WriteString(", ")
		}

		operator.Column(column, builder)
	}
	builder.WriteString(")")
}

type PartitionOptionsSub struct {
	hash *PartitionOptionsHash
	key  *PartitionOptionsKey
	num  uint
}

func (p *PartitionOptionsSub) Hash(expr string) ksql.PartitionOptionsSubInterface {
	p.hash = &PartitionOptionsHash{expr: expr}
	return p
}

func (p *PartitionOptionsSub) Key(columns []string) ksql.PartitionOptionsKeyInterface {
	p.key = &PartitionOptionsKey{columns: columns}
	return p.key
}

func (p *PartitionOptionsSub) Build(builder *strings.Builder) {
	if p.hash != nil {
		p.hash.Build(builder)
	}

	if p.key != nil {
		p.key.Build(builder)
	}

	builder.WriteString(" PARTITIONS ")
	builder.WriteString(strconv.Itoa(int(p.num)))
}

type PartitionDefinitionLessthan struct {
	expr      string
	valueList []string
	maxValue  string
}

func (p *PartitionDefinitionLessthan) Expr(expr string) ksql.PartitionDefinitionLessthanInterface {
	if p.valueList != nil || p.maxValue != "" {
		return p
	}

	p.expr = expr
	return p
}

func (p *PartitionDefinitionLessthan) ValueList(valueList []string) ksql.PartitionDefinitionLessthanInterface {
	if p.expr != "" || p.maxValue != "" {
		return p
	}

	p.valueList = valueList
	return p
}

func (p *PartitionDefinitionLessthan) Build(builder *strings.Builder) {
	builder.WriteString(" LESS THAN ")
	if p.expr != "" {
		builder.WriteString("(")
		builder.WriteString(p.expr)
		builder.WriteString(")")
		return
	}

	if len(p.valueList) > 0 {
		builder.WriteString("(")
		for index, value := range p.valueList {
			if index > 0 {
				builder.WriteString(", ")
			}

			if _, err := strconv.ParseInt(value, 10, 64); err == nil {
				builder.WriteString(value)
			} else if _, err := strconv.ParseFloat(value, 64); err == nil {
				builder.WriteString(value)
			} else {
				operator.Quote(value, builder)
			}
		}

		builder.WriteString(")")
		return
	}

	if p.maxValue != "" {
		builder.WriteString(p.maxValue)
	}
}

func (p *PartitionDefinitionLessthan) MaxValue() ksql.PartitionDefinitionLessthanInterface {
	if p.expr != "" || p.valueList != nil {
		return p
	}

	p.maxValue = "MAXVALUE"
	return p
}

type PartitionDefinitionIn struct {
	valueList []string
}

func (p *PartitionDefinitionIn) Build(builder *strings.Builder) {
	if len(p.valueList) == 0 {
		return
	}
	builder.WriteString(" IN (")
	for index, value := range p.valueList {
		if index > 0 {
			builder.WriteString(", ")
		}

		if _, err := strconv.ParseInt(value, 10, 64); err == nil {
			builder.WriteString(value)
		} else if _, err := strconv.ParseFloat(value, 64); err == nil {
			builder.WriteString(value)
		} else {
			operator.Quote(value, builder)
		}
	}

	builder.WriteString(")")
}

type PartitionDefinitionValues struct {
	lessThan *PartitionDefinitionLessthan
	in       *PartitionDefinitionIn
}

func (p *PartitionDefinitionValues) LessThan() ksql.PartitionDefinitionLessthanInterface {
	p.lessThan = &PartitionDefinitionLessthan{}
	return p.lessThan
}

func (p *PartitionDefinitionValues) In(valueList []string) ksql.PartitionDefinitionInInterface {
	p.in = &PartitionDefinitionIn{valueList: valueList}
	return p.in
}

func (p *PartitionDefinitionValues) Build(builder *strings.Builder) {
	builder.WriteString(" VALUES ")
	if p.lessThan != nil {
		p.lessThan.Build(builder)
		return
	}

	if p.in != nil {
		p.in.Build(builder)
	}
}

type PartitionDefinitionOption struct {
	key   ksql.PartitionDefinitionOptKey
	value string
	isStr bool
}

func (p *PartitionDefinitionOption) Build(bulder *strings.Builder) {
	bulder.WriteString(string(p.key))
	bulder.WriteString(" = ")
	if p.isStr {
		operator.Quote(p.value, bulder)
	} else {
		bulder.WriteString(p.value)
	}
}

type PartitionDefinitionSub struct {
	name    string
	options []*PartitionDefinitionOption
}

func (p *PartitionDefinitionSub) Option(key ksql.PartitionDefinitionOptKey, value string) ksql.PartitionDefinitionSubInterface {
	p.options = append(p.options, &PartitionDefinitionOption{key: key, value: value, isStr: key.IsStr()})
	return p
}

func (p *PartitionDefinitionSub) Build(bulder *strings.Builder) {
	bulder.WriteString(" SUBPARTITION ")
	operator.Backtick(p.name, bulder)
	if len(p.options) > 0 {
		bulder.WriteString(" ")
	}
	for index, option := range p.options {
		if index > 0 {
			bulder.WriteString(", ")
		}

		option.Build(bulder)
	}
}

type PartitionDefinition struct {
	name     string
	values   *PartitionDefinitionValues
	options  []*PartitionDefinitionOption
	subs     []*PartitionDefinitionSub
	onlyBody bool
}

func (p *PartitionDefinition) Values() ksql.PartitionDefinitionValuesInterface {
	p.values = &PartitionDefinitionValues{}
	return p.values
}

func (p *PartitionDefinition) Option(key ksql.PartitionDefinitionOptKey, value string) ksql.PartitionDefinitionInterface {
	p.options = append(p.options, &PartitionDefinitionOption{key: key, value: value, isStr: key.IsStr()})
	return p
}

func (p *PartitionDefinition) Sub(name string) ksql.PartitionDefinitionSubInterface {
	sub := &PartitionDefinitionSub{name: name}
	p.subs = append(p.subs, sub)
	return sub
}

func (p *PartitionDefinition) Build(builder *strings.Builder) {
	if !p.onlyBody {
		builder.WriteString(" PARTITION ")
	}
	operator.Backtick(p.name, builder)

	if p.values != nil {
		p.values.Build(builder)
	}

	for _, option := range p.options {
		builder.WriteString(" ")
		option.Build(builder)
	}

	for _, sub := range p.subs {
		builder.WriteString(" ")
		sub.Build(builder)
	}
}

type PartitionOptions struct {
	hash    *PartitionOptionsHash
	key     *PartitionOptionsKey
	rge     *PartitionOptionsRange
	list    *PartitionOptionsRange
	num     uint
	sub     *PartitionOptionsSub
	defines []*PartitionDefinition
}

func (p *PartitionOptions) Hash(expr string) ksql.PartitionOptionsHashInterface {
	p.hash = &PartitionOptionsHash{expr: expr}
	return p.hash
}

func (p *PartitionOptions) Key(columns []string) ksql.PartitionOptionsKeyInterface {
	p.key = &PartitionOptionsKey{columns: columns}
	return p.key
}

func (p *PartitionOptions) Range() ksql.PartitionOptionsRangeInterface {
	p.rge = &PartitionOptionsRange{}
	return p.rge
}

func (p *PartitionOptions) List() ksql.PartitionOptionsRangeInterface {
	p.list = &PartitionOptionsRange{}
	return p.list
}

func (p *PartitionOptions) Sub() ksql.PartitionOptionsSubInterface {
	p.sub = &PartitionOptionsSub{}
	return p.sub
}

func (p *PartitionOptions) Definition(name string) ksql.PartitionDefinitionInterface {
	define := &PartitionDefinition{name: name}
	p.defines = append(p.defines, define)
	return define
}

func (p *PartitionOptions) Build(builder *strings.Builder) {
	builder.WriteString(" PARTITION BY")
	if p.hash != nil {
		p.hash.Build(builder)
	}

	if p.key != nil {
		p.key.Build(builder)
	}

	if p.rge != nil {
		p.rge.Build(builder)
	}

	if p.list != nil {
		p.list.Build(builder)
	}

	if p.num > 0 {
		builder.WriteString(" PARTITIONS ")
		builder.WriteString(strconv.Itoa(int(p.num)))
	}

	if p.sub != nil {
		p.sub.Build(builder)
	}

	for _, define := range p.defines {
		builder.WriteString(" ")
		define.Build(builder)
	}
}
