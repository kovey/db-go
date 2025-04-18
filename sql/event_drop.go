package sql

import (
	ksql "github.com/kovey/db-go/v3"
)

type EventDrop struct {
	*drop
}

func NewEventDrop() *EventDrop {
	return &EventDrop{drop: newDrop("EVENT")}
}

func (s *EventDrop) Event(event string) ksql.DropEventInterface {
	s.name = event
	return s
}

func (s *EventDrop) IfExists() ksql.DropEventInterface {
	s.ifExists = true
	return s
}
