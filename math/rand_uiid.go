package math

import (
	"github.com/google/uuid"
	"github.com/rs/xid"
)

func RandID(prefix string) string {
	prefix += "_" + uuid.New().String()
	return prefix
}

func RandXID(prefix string) string {
	prefix += "_" + xid.New().String()
	return prefix
}
