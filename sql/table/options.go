package table

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type Option struct {
	prefix      ksql.TableOptPrefix
	key         ksql.TableOptKey
	value       string
	isStr       bool
	unions      []string
	spaceOption *TablespaceOption
}

func (o *Option) Build(builder *strings.Builder) {
	if o.spaceOption != nil {
		o.spaceOption.Build(builder)
		return
	}

	if o.prefix != "" {
		builder.WriteString(string(o.prefix))
		builder.WriteString(" ")
	}

	builder.WriteString(string(o.key))
	if o.value != "" {
		builder.WriteString(" = ")
	}

	if len(o.unions) > 0 {
		builder.WriteString(" = ")
		builder.WriteString("(")
		for index, union := range o.unions {
			if index > 0 {
				builder.WriteString(", ")
			}

			operator.Column(union, builder)
		}
		builder.WriteString(")")
		return
	}

	if o.isStr {
		operator.Quote(o.value, builder)
	} else {
		builder.WriteString(o.value)
	}
}

type TablespaceOption struct {
	name    string
	storage ksql.ColumnStorage
}

func (t *TablespaceOption) Build(builder *strings.Builder) {
	if t.name != "" {
		builder.WriteString(" TABLESPACE ")
		operator.Backtick(t.name, builder)
	}

	if t.storage != "" {
		builder.WriteString(" STORAGE ")
		builder.WriteString(string(t.storage))
	}
}

type Options struct {
	options      []*Option
	spaceOptions []*TablespaceOption
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Append(key ksql.TableOptKey, value string) ksql.TableOptionsInterface {
	o.options = append(o.options, &Option{key: key, value: value, isStr: key.IsStr()})
	return o
}

func (o *Options) AppendWithPrefix(prefix ksql.TableOptPrefix, key ksql.TableOptKey, value string) ksql.TableOptionsInterface {
	o.options = append(o.options, &Option{prefix: prefix, key: key, value: value, isStr: key.IsStr()})
	return o
}

func (o *Options) AppendWith(key ksql.TableOptKey, value ksql.TableOptVal) ksql.TableOptionsInterface {
	o.options = append(o.options, &Option{key: key, value: string(value)})
	return o
}

func (o *Options) InsertMethod(value ksql.InsertMethod) ksql.TableOptionsInterface {
	o.options = append(o.options, &Option{key: "INSERT_METHOD", value: string(value)})
	return o
}

func (o *Options) StartTransation() ksql.TableOptionsInterface {
	o.options = append(o.options, &Option{key: "START TRANSACTION"})
	return o
}

func (o *Options) RowFormat(value ksql.RowFormat) ksql.TableOptionsInterface {
	o.options = append(o.options, &Option{key: "ROW_FORMAT", value: string(value)})
	return o
}

func (o *Options) Build(builder *strings.Builder) {
	index := 0
	for _, option := range o.options {
		if index > 0 {
			builder.WriteString(", ")
		}

		option.Build(builder)
		index++
	}

	for _, option := range o.spaceOptions {
		if index > 0 {
			builder.WriteString(",")
		}

		option.Build(builder)
		index++
	}
}

func (o *Options) Union(tables ...string) ksql.TableOptionsInterface {
	o.options = append(o.options, &Option{unions: tables, key: "UNION"})
	return o
}

func (o *Options) SpaceOption(name string, storage ksql.ColumnStorage) ksql.TableOptionsInterface {
	o.spaceOptions = append(o.spaceOptions, &TablespaceOption{name: name, storage: storage})
	return o
}

func (o *Options) Empty() bool {
	return len(o.options) == 0 && len(o.spaceOptions) == 0
}
