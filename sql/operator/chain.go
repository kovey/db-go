package operator

import "strings"

type Operator func(*strings.Builder)

type Chain struct {
	chains []Operator
	isCall bool
}

func NewChain() *Chain {
	return &Chain{isCall: false}
}

func (o *Chain) Append(op ...Operator) *Chain {
	o.chains = append(o.chains, op...)
	return o
}

func (o *Chain) Call(builder *strings.Builder) *Chain {
	if o.isCall {
		return o
	}

	o.isCall = true
	for _, op := range o.chains {
		op(builder)
	}

	return o
}

func (o *Chain) Reset() {
	o.isCall = false
}
