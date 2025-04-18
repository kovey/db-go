package sql

import (
	"strconv"
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type srsAttr struct {
	key   ksql.SrsAttr
	value string
	orgId uint32
}
type SpatialReferenceSystem struct {
	*base
	isRelace    bool
	ifNotExists bool
	srid        uint32
	attrs       []*srsAttr
}

func NewSpatialReferenceSystem() *SpatialReferenceSystem {
	s := &SpatialReferenceSystem{base: newBase()}
	s.opChain.Append(keywordCreate, s._keyword, s._srid, s._sridAttrs)
	return s
}

func (s *SpatialReferenceSystem) _keyword(builder *strings.Builder) {
	if s.isRelace {
		builder.WriteString(" OR REPLACE SPATIAL REFERENCE SYSTEM")
		return
	}

	builder.WriteString(" SPATIAL REFERENCE SYSTEM")
	if s.ifNotExists {
		builder.WriteString(" IF NOT EXISTS")
	}
}

func (s *SpatialReferenceSystem) _srid(builder *strings.Builder) {
	builder.WriteString(" ")
	builder.WriteString(strconv.Itoa(int(s.srid)))
}

func (s *SpatialReferenceSystem) _sridAttrs(builder *strings.Builder) {
	builder.WriteString(" ")
	for index, attr := range s.attrs {
		if index > 0 {
			builder.WriteString(" ")
		}

		builder.WriteString(string(attr.key))
		builder.WriteString(" ")
		Quote(attr.value, builder)
		if attr.orgId > 0 {
			builder.WriteString(" IDENTIFIED BY ")
			builder.WriteString(strconv.Itoa(int(attr.orgId)))
		}
	}
}

func (s *SpatialReferenceSystem) Replace() ksql.SpatialReferenceSystemInterface {
	if s.ifNotExists {
		return s
	}

	s.isRelace = true
	return s
}

func (s *SpatialReferenceSystem) IfNotExists() ksql.SpatialReferenceSystemInterface {
	if s.isRelace {
		return s
	}

	s.ifNotExists = true
	return s
}

func (s *SpatialReferenceSystem) Srid(srid uint32) ksql.SpatialReferenceSystemInterface {
	s.srid = srid
	return s
}

func (s *SpatialReferenceSystem) Atrribute(key ksql.SrsAttr, value string) ksql.SpatialReferenceSystemInterface {
	s.attrs = append(s.attrs, &srsAttr{key: key, value: value})
	return s
}

func (s *SpatialReferenceSystem) Organization(value string, identified uint32) ksql.SpatialReferenceSystemInterface {
	s.attrs = append(s.attrs, &srsAttr{key: "ORGANIZATION", value: value, orgId: identified})
	return s
}
