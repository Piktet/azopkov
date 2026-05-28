package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

func CreateShort(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

var ErrConflict = errors.New("already exist")
