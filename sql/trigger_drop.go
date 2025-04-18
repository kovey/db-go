package sql

import (
	ksql "github.com/kovey/db-go/v3"
)

type TriggerDrop struct {
	*drop
}

func NewTriggerDrop() *TriggerDrop {
	i := &TriggerDrop{drop: newDrop("TRIGGER")}
	return i
}

func (s *TriggerDrop) Trigger(event string) ksql.DropTriggerInterface {
	s.name = event
	return s
}

func (s *TriggerDrop) Schema(schema string) ksql.DropTriggerInterface {
	s.schema = schema
	return s
}

func (s *TriggerDrop) IfExists() ksql.DropTriggerInterface {
	s.ifExists = true
	return s
}
