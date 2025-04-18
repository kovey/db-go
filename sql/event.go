package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

func keywordEvent(builder *strings.Builder) {
	builder.WriteString(" EVENT")
}

func keywordOnSchedule(builder *strings.Builder) {
	builder.WriteString(" ON SCHEDULE")
}

type intervalInfo struct {
	value string
	unit  ksql.IntervalUnit
}

type Event struct {
	*base
	event          string
	definer        string
	ifNotExists    bool
	comment        string
	does           []string
	at             string
	atIntervals    []*intervalInfo
	every          *intervalInfo
	starts         string
	startsInterval []*intervalInfo
	ends           string
	endsInterval   []*intervalInfo
	status         ksql.EventStatus
	completion     string
	newName        string
	isCreate       bool
}

func NewEvent() *Event {
	e := &Event{base: newBase(), isCreate: true}
	e.opChain.Append(e._keyword, e._definer, keywordEvent, e._event, e._keywordOnSchedule, e._at, e._every, e._completion, e._rename, e._status, e._comment, e._do)
	return e
}

func (e *Event) _keyword(builder *strings.Builder) {
	if e.isCreate {
		keywordCreate(builder)
		return
	}

	keywordAlter(builder)
}

func (e *Event) _keywordOnSchedule(builder *strings.Builder) {
	if e.isCreate {
		keywordOnSchedule(builder)
		return
	}

	if e.at != "" || e.every != nil {
		keywordOnSchedule(builder)
	}
}

func (e *Event) _definer(builder *strings.Builder) {
	if e.definer == "" {
		return
	}

	builder.WriteString(" DEFINER = ")
	builder.WriteString(e.definer)
}

func (e *Event) _event(builder *strings.Builder) {
	builder.WriteString(" ")
	if e.ifNotExists {
		builder.WriteString("IF NOT EXISTS ")
	}
	Backtick(e.event, builder)
}

func (e *Event) _at(builder *strings.Builder) {
	if e.at == "" {
		return
	}

	builder.WriteString(" AT ")
	Quote(e.at, builder)
	e._intervals(e.atIntervals, builder)
}

func (e *Event) _every(builder *strings.Builder) {
	if e.every == nil {
		return
	}

	builder.WriteString(" EVERY ")
	builder.WriteString("+ INTERVAL ")
	if e.every.unit.IsNumber() {
		builder.WriteString(e.every.value)
	} else {
		Quote(e.every.value, builder)
	}
	builder.WriteString(" ")
	builder.WriteString(string(e.every.unit))

	if e.starts != "" {
		builder.WriteString(" STARTS ")
		Quote(e.starts, builder)
		e._intervals(e.startsInterval, builder)
	}

	if e.ends != "" {
		builder.WriteString(" ENDS ")
		Quote(e.ends, builder)
		e._intervals(e.endsInterval, builder)
	}
}

func (e *Event) _intervals(intervals []*intervalInfo, builder *strings.Builder) {
	for _, interval := range intervals {
		builder.WriteString(" ")
		builder.WriteString("+ INTERVAL ")
		if interval.unit.IsNumber() {
			builder.WriteString(interval.value)
		} else {
			Quote(interval.value, builder)
		}
		builder.WriteString(" ")
		builder.WriteString(string(interval.unit))
	}
}

func (e *Event) _comment(builder *strings.Builder) {
	if e.comment == "" {
		return
	}

	builder.WriteString(" COMMENT ")
	Quote(e.comment, builder)
}

func (e *Event) _completion(builder *strings.Builder) {
	if e.completion == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(e.completion)
}

func (e *Event) _rename(builder *strings.Builder) {
	if e.newName == "" {
		return
	}

	builder.WriteString(" RENAME TO ")
	Backtick(e.newName, builder)
}

func (e *Event) _status(builder *strings.Builder) {
	if e.status == "" {
		return
	}

	builder.WriteString(" ")
	builder.WriteString(string(e.status))
}

func (e *Event) _do(builder *strings.Builder) {
	if len(e.does) == 0 {
		return
	}

	builder.WriteString(" DO BEGIN ")
	for _, do := range e.does {
		builder.WriteString(do)
		builder.WriteString("; ")
	}
	builder.WriteString("END")
}

func (e *Event) Status(status ksql.EventStatus) ksql.EventInterface {
	e.status = status
	return e
}

func (e *Event) Definer(name string) ksql.EventInterface {
	e.definer = name
	return e
}
func (e *Event) IfNotExists() ksql.EventInterface {
	e.ifNotExists = true
	return e
}

func (e *Event) Event(name string) ksql.EventInterface {
	e.event = name
	return e
}

func (e *Event) Comment(comment string) ksql.EventInterface {
	e.comment = comment
	return e
}

func (e *Event) Do(sql ksql.SqlInterface) ksql.EventInterface {
	if sql == nil {
		return e
	}

	e.does = append(e.does, sql.Prepare())
	e.binds = append(e.binds, sql.Binds()...)
	return e
}

func (e *Event) DoRaw(sql ksql.ExpressInterface) ksql.EventInterface {
	if sql == nil {
		return e
	}
	e.does = append(e.does, sql.Statement())
	e.binds = append(e.binds, sql.Binds()...)
	return e
}

func (e *Event) At(timestamp string) ksql.EventInterface {
	if e.every != nil {
		return e
	}

	e.at = timestamp
	return e
}

func (e *Event) AtInterval(interval string, unit ksql.IntervalUnit) ksql.EventInterface {
	e.atIntervals = append(e.atIntervals, &intervalInfo{value: interval, unit: unit})
	return e
}

func (e *Event) Every(interval string, unit ksql.IntervalUnit) ksql.EventInterface {
	if e.at != "" {
		return e
	}

	e.every = &intervalInfo{value: interval, unit: unit}
	return e
}

func (e *Event) Starts(timestamp string) ksql.EventInterface {
	e.starts = timestamp
	return e
}

func (e *Event) StartsInterval(interval string, unit ksql.IntervalUnit) ksql.EventInterface {
	e.startsInterval = append(e.startsInterval, &intervalInfo{value: interval, unit: unit})
	return e
}

func (e *Event) Ends(timestamp string) ksql.EventInterface {
	e.ends = timestamp
	return e
}

func (e *Event) EndsInterval(interval string, unit ksql.IntervalUnit) ksql.EventInterface {
	e.endsInterval = append(e.endsInterval, &intervalInfo{value: interval, unit: unit})
	return e
}

func (e *Event) OnCompletion() ksql.EventInterface {
	e.completion = "ON COMPLETION PRESERVE"
	return e
}

func (e *Event) OnCompletionNot() ksql.EventInterface {
	e.completion = "ON COMPLETION NOT PRESERVE"
	return e
}

func (e *Event) Alter() ksql.EventInterface {
	e.isCreate = false
	return e
}

func (e *Event) Rename(name string) ksql.EventInterface {
	if e.isCreate {
		return e
	}

	e.newName = name
	return e
}
