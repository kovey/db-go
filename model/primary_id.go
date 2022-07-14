package model

import "reflect"

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

func (p PrimaryId) Parse(value reflect.Value) {
	switch p.Type {
	case Int:
		p.IntValue = int(value.Int())
	case Str:
		p.StrValue = value.String()
	}
}

func (p PrimaryId) Value() interface{} {
	switch p.Type {
	case Int:
		return p.IntValue
	case Str:
		return p.StrValue
	default:
		return nil
	}
}

func (p PrimaryId) Null() bool {
	return p.IntValue == 0 && p.StrValue == ""
}
