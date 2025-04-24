package ksql

type CallInterface interface {
	SqlInterface
	Call(spName string) CallInterface
	Params(params ...string) CallInterface
}

type DoInterface interface {
	SqlInterface
	Do(expr ExpressInterface) DoInterface
}
