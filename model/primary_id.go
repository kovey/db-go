package model

import (
	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

type PrimaryIdType int

const (
	Int          PrimaryIdType = 1
	Str          PrimaryIdType = 2
	namespace                  = "ko.db.model"
	primary_name               = "Primary"
)

func init() {
	pool.DefaultNoCtx(namespace, primary_name, func() any {
		return &PrimaryId{IsAutoInc: true, ObjNoCtx: object.NewObjNoCtx(namespace, primary_name)}
	})
}

type PrimaryId struct {
	*object.ObjNoCtx
	Type      PrimaryIdType
	Name      string
	IntValue  int
	StrValue  string
	IsAutoInc bool
}

func NewPrimaryId(name string, t PrimaryIdType) *PrimaryId {
	return &PrimaryId{Name: name, Type: t, IsAutoInc: true}
}

func NewPrimaryIdBy(ctx object.CtxInterface, name string, t PrimaryIdType) *PrimaryId {
	obj := ctx.GetNoCtx(namespace, primary_name).(*PrimaryId)
	obj.Name = name
	obj.Type = t
	return obj
}

func (p *PrimaryId) Reset() {
	p.Name = ""
	p.Type = Int
}

func (p *PrimaryId) Parse(val any) {
	switch p.Type {
	case Int:
		switch tmp := val.(type) {
		case int:
			p.IntValue = tmp
		case int8:
			p.IntValue = int(tmp)
		case int16:
			p.IntValue = int(tmp)
		case int32:
			p.IntValue = int(tmp)
		case int64:
			p.IntValue = int(tmp)
		case uint:
			p.IntValue = int(tmp)
		case uint8:
			p.IntValue = int(tmp)
		case uint16:
			p.IntValue = int(tmp)
		case uint32:
			p.IntValue = int(tmp)
		case uint64:
			p.IntValue = int(tmp)
		default:
			panic("val is not integer")
		}
	case Str:
		p.StrValue = val.(string)
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
