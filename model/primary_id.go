package model

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
