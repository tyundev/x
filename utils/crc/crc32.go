package crc

import (
	"hash/crc32"
)

func GetCrc32(value string) int64 {
	crc32q := crc32.MakeTable(0xD5828281)
	return int64(crc32.Checksum([]byte(value), crc32q))
}
