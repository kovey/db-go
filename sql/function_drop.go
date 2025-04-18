package sql

import (
	ksql "github.com/kovey/db-go/v3"
)

type FunctionDrop struct {
	*drop
}

func NewFunctionDrop() *FunctionDrop {
	return &FunctionDrop{drop: newDrop("FUNCTION")}
}

func (s *FunctionDrop) Function(event string) ksql.DropFunctionInterface {
	s.name = event
	return s
}

func (s *FunctionDrop) IfExists() ksql.DropFunctionInterface {
	s.ifExists = true
	return s
}
