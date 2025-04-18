package sql

import (
	ksql "github.com/kovey/db-go/v3"
)

type ServerDrop struct {
	*drop
}

func NewServerDrop() *ServerDrop {
	return &ServerDrop{drop: newDrop("SERVER")}
}

func (s *ServerDrop) Server(event string) ksql.DropServerInterface {
	s.name = event
	return s
}

func (s *ServerDrop) IfExists() ksql.DropServerInterface {
	s.ifExists = true
	return s
}
