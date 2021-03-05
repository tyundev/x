package math

import (
	"github.com/google/uuid"
)

func RandID(prefix string) string {
	prefix += "_" + uuid.New().String()
	return prefix
}
