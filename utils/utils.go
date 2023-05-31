package utils

import (
	"math/rand"
	"time"
)

// Returns a random integer between min (inclusive) and max (exclusive).
func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Intn(max-min)
}
