package db

type Int interface {
	int | int8 | int16 | int32 | int64
}

type UInt interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type Number interface {
	Int | UInt
}

type FindType interface {
	Number | string
}
