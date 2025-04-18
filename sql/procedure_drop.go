package sql

import (
	ksql "github.com/kovey/db-go/v3"
)

type ProcedureDrop struct {
	*drop
}

func NewProcedureDrop() *ProcedureDrop {
	return &ProcedureDrop{drop: newDrop("PROCEDURE")}
}

func (s *ProcedureDrop) Procedure(event string) ksql.DropProcedureInterface {
	s.name = event
	return s
}

func (s *ProcedureDrop) IfExists() ksql.DropProcedureInterface {
	s.ifExists = true
	return s
}
