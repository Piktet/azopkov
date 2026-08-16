// Package utils вспомогательные функции
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

// CreateShort генерирует случайную строку заданной длины, используя base64-кодирование.
// Возвращает ошибку, если не удалось прочитать данные из криптографического генератора случайных чисел.
func CreateShort(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

// ErrConflict возвращается при попытке создать сокращённый URL, который уже существует в хранилище.
var ErrConflict = errors.New("already exist")
