package sql

import (
	ksql "github.com/kovey/db-go/v3"
)

type SpatialReferenceSystemDrop struct {
	*drop
}

func NewSpatialReferenceSystemDrop() *SpatialReferenceSystemDrop {
	return &SpatialReferenceSystemDrop{drop: newDrop("SPATIAL REFERENCE SYSTEM")}
}

func (s *SpatialReferenceSystemDrop) Srid(event string) ksql.DropSpatialReferenceSystemInterface {
	s.name = event
	return s
}

func (s *SpatialReferenceSystemDrop) IfExists() ksql.DropSpatialReferenceSystemInterface {
	s.ifExists = true
	return s
}
