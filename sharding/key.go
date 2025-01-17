package sharding

import (
	"bytes"
	"encoding/binary"
)

type String interface {
	String() string
}

func str2uint64(data string) uint64 {
	bits := []byte(data)
	length := int64(len(bits))
	mod := length & 63
	if mod != 0 {
		tmp := make([]byte, 64-mod)
		bits = append(bits, tmp...)
		length += 64 - mod
	}

	i := int64(0)
	dt := make([]byte, 8)
	for i < length {
		for j := int64(0); j < 8; j++ {
			dt[j] = byte(int64(dt[j]) + int64(bits[i+j])&255)
		}
		i += 8
	}

	var res uint64
	binary.Read(bytes.NewBuffer(dt), binary.BigEndian, &res)
	return res
}

func node(key any, total int) int {
	switch tmp := key.(type) {
	case string:
		return int(str2uint64(tmp) % uint64(total))
	case int:
		return tmp % total
	case int8:
		return int(tmp) % total
	case int16:
		return int(tmp) % total
	case int32:
		return int(tmp) % total
	case int64:
		return int(tmp) % total
	case uint:
		return int(tmp % uint(total))
	case uint8:
		return int(tmp) % int(total)
	case uint16:
		return int(tmp) % total
	case uint32:
		return int(tmp % uint32(total))
	case uint64:
		return int(tmp % uint64(total))
	}

	if tmp, ok := key.(String); ok {
		return int(str2uint64(tmp.String()) % uint64(total))
	}

	return 0
}
