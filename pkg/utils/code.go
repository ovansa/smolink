package utils

import (
	"math/rand"
	"time"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateShortCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return string(b)
}
