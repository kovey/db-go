package sql

import (
	"testing"

	ksql "github.com/kovey/db-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestSpatialReferenceSystemReplace(t *testing.T) {
	s := NewSpatialReferenceSystem().Srid(123).Replace().IfNotExists().Atrribute(ksql.Srs_Attr_Name, "kovey").Organization("aaaa", 11)
	s.Atrribute(ksql.Srs_Attr_Definition, "aaaaa").Atrribute(ksql.Srs_Attr_Description, "desc")
	assert.Equal(t, "CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 123 NAME 'kovey' ORGANIZATION 'aaaa' IDENTIFIED BY 11 DEFINITION 'aaaaa' DESCRIPTION 'desc'", s.Prepare())
	assert.Nil(t, s.Binds())
}

func TestSpatialReferenceSystemIfNotExists(t *testing.T) {
	s := NewSpatialReferenceSystem().Srid(123).IfNotExists().Replace().Atrribute(ksql.Srs_Attr_Name, "kovey").Organization("aaaa", 11)
	s.Atrribute(ksql.Srs_Attr_Definition, "aaaaa").Atrribute(ksql.Srs_Attr_Description, "desc")
	assert.Equal(t, "CREATE SPATIAL REFERENCE SYSTEM IF NOT EXISTS 123 NAME 'kovey' ORGANIZATION 'aaaa' IDENTIFIED BY 11 DEFINITION 'aaaaa' DESCRIPTION 'desc'", s.Prepare())
	assert.Nil(t, s.Binds())
}

func TestSpatialReferenceSystem(t *testing.T) {
	s := NewSpatialReferenceSystem().Srid(123).Atrribute(ksql.Srs_Attr_Name, "kovey").Organization("aaaa", 11)
	s.Atrribute(ksql.Srs_Attr_Definition, "aaaaa").Atrribute(ksql.Srs_Attr_Description, "desc")
	assert.Equal(t, "CREATE SPATIAL REFERENCE SYSTEM 123 NAME 'kovey' ORGANIZATION 'aaaa' IDENTIFIED BY 11 DEFINITION 'aaaaa' DESCRIPTION 'desc'", s.Prepare())
	assert.Nil(t, s.Binds())
}
