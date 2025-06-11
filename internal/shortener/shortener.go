package shortener

import (
	"math/rand/v2"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const keyLength = 7

// Generate Short Key generates a random short key of a specific length.
// This function is now safe to be called by multiple goroutines at the same time.
func GenerateShortKey() string {
	b := make([]byte, keyLength)
	for i := range b {
		b[i] = charset[rand.N(uint64(len(charset)))]
	}
	return string(b)
}
