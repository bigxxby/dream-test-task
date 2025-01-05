package utils

import "math/rand"

func GenerateShortLink() string {
	// Example method to generate a random short link. You can replace it with a more sophisticated approach.
	return RandStringBytes(6)
}

func RandStringBytes(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
