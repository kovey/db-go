package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type tablespaceOption struct {
	prefix string
	key    string
	value  string
	isStr  bool
}

func (t *tablespaceOption) Build(builder *strings.Builder) {
	operator.BuildPureString(t.prefix, builder)
	operator.BuildPureString(t.key, builder)
	if t.isStr {
		operator.BuildQuoteString(t.value, builder)
	} else {
		operator.BuildPureString(t.value, builder)
	}
}

type Tablespace struct {
	*base
	name     string
	options  []*tablespaceOption
	isCreate bool
	undo     string
}

func NewTablespace() *Tablespace {
	t := &Tablespace{base: newBase(), isCreate: true}
	t.opChain.Append(t._keyword, t._option)
	return t
}

func (t *Tablespace) _keyword(builder *strings.Builder) {
	if t.isCreate {
		builder.WriteString("CREATE")
	} else {
		builder.WriteString("ALTER")
	}

	if t.undo != "" {
		builder.WriteString(" ")
		builder.WriteString(t.undo)
	}

	builder.WriteString(" TABLESPACE")
}

func (t *Tablespace) _option(builder *strings.Builder) {
	for _, option := range t.options {
		option.Build(builder)
	}
}

func (t *Tablespace) Alter() ksql.TablespaceInterface {
	t.isCreate = false
	return t
}

func (t *Tablespace) Undo() ksql.TablespaceInterface {
	t.undo = "UNDO"
	return t
}

func (t *Tablespace) Tablespace(tablespace string) ksql.TablespaceInterface {
	t.name = tablespace
	return t
}

func (t *Tablespace) Option(key, value string) ksql.TablespaceInterface {
	t.options = append(t.options, &tablespaceOption{key: key, value: value})
	return t
}

func (t *Tablespace) OptionWith(prefix, key, value string) ksql.TablespaceInterface {
	t.options = append(t.options, &tablespaceOption{prefix: prefix, key: key, value: value})
	return t
}

func (t *Tablespace) OptionStr(key, value string) ksql.TablespaceInterface {
	t.options = append(t.options, &tablespaceOption{key: key, value: value, isStr: true})
	return t
}

func (t *Tablespace) OptionStrWith(prefix, key, value string) ksql.TablespaceInterface {
	t.options = append(t.options, &tablespaceOption{prefix: prefix, key: key, value: value, isStr: true})
	return t
}

func (t *Tablespace) OptionOnlyKey(key string) ksql.TablespaceInterface {
	t.options = append(t.options, &tablespaceOption{key: key})
	return t
}
