package utils

import (
	"math/rand"
)

func RandomBoolean() bool {
	return rand.Float32() < 0.5
}
