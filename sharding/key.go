package sharding

import "hash/crc32"

func getHashKey(key string) uint32 {
	if len(key) >= 64 {
		return crc32.ChecksumIEEE([]byte(key))
	}

	var buf [64]byte
	copy(buf[:], key)
	return crc32.ChecksumIEEE(buf[:len(key)])
}
