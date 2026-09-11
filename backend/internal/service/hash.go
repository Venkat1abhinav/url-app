package service

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func ShortenURL(id int64) (string, error) {
	if id < 0 {
		return "", errors.New("id cannot be negative")
	}

	random, err := randomCode(4)
	if err != nil {
		return "", err
	}

	encodedID := Base62(uint64(id), 4)

	return random + encodedID, nil
}

func randomCode(length int) (string, error) {
	result := make([]byte, length)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}

		result[i] = alphabet[n.Int64()]
	}

	return string(result), nil
}

func Base62(n uint64, length int) string {
	result := make([]byte, length)

	for i := range result {
		result[i] = alphabet[0]
	}

	for i := length - 1; i >= 0 && n > 0; i-- {
		result[i] = alphabet[n%62]
		n /= 62
	}

	if n > 0 {
		return ""
	}

	return string(result)
}
