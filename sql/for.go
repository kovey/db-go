package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/sql/operator"
)

type For struct {
	typ             string
	of              []string
	noWait          string
	lockInShareMode string
}

func (f *For) Empty() bool {
	return f.typ == "" && len(f.of) == 0 && f.noWait == "" && f.lockInShareMode == ""
}

func (f *For) Update() ksql.ForInterface {
	f.typ = "UPDATE"
	return f
}

func (f *For) Share() ksql.ForInterface {
	f.typ = "SHARE"
	return f
}

func (f *For) Of(tables ...string) ksql.ForInterface {
	f.of = append(f.of, tables...)
	return f
}

func (f *For) NoWait() ksql.ForInterface {
	f.noWait = "NOWAIT"
	f.lockInShareMode = ""
	return f
}

func (f *For) SkipLocked() ksql.ForInterface {
	f.noWait = "SKIP LOCKED"
	f.lockInShareMode = ""
	return f
}

func (f *For) LockInShareMode() ksql.ForInterface {
	f.lockInShareMode = "LOCK IN SHARE MODE"
	f.noWait = ""
	f.of = nil
	return f
}

func (f *For) Build(builder *strings.Builder) {
	builder.WriteString(" FOR")
	operator.BuildPureString(f.typ, builder)
	if len(f.of) > 0 {
		builder.WriteString(" ")
		for index, o := range f.of {
			if index > 0 {
				builder.WriteString(",")
			}
			operator.BuildColumnString(o, builder)
		}
	}
	operator.BuildPureString(f.noWait, builder)
	operator.BuildPureString(f.lockInShareMode, builder)
}
