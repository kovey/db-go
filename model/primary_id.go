package model

import (
	"reflect"
)

type PrimaryIdType int

const (
	Int PrimaryIdType = 1
	Str PrimaryIdType = 2
)

type PrimaryId struct {
	Type     PrimaryIdType
	Name     string
	IntValue int
	StrValue string
}

func NewPrimaryId(name string, t PrimaryIdType) *PrimaryId {
	return &PrimaryId{Name: name, Type: t}
}

func (p *PrimaryId) Parse(value reflect.Value) {
	switch p.Type {
	case Int:
		p.IntValue = int(value.Int())
		break
	case Str:
		p.StrValue = value.String()
		break
	}
}

func (p *PrimaryId) Value() any {
	switch p.Type {
	case Int:
		return p.IntValue
	case Str:
		return p.StrValue
	default:
		return nil
	}
}

func (p *PrimaryId) Null() bool {
	return p.IntValue == 0 && p.StrValue == ""
}
