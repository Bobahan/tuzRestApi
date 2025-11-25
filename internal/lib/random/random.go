package random

import (
	"math/rand"
	"time"
)

func NewRandomString(size int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)

	for i := range b {
		b[i] = chars[random.Intn(len(chars))]
	}

	return string(b)
}
