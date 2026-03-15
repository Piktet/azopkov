package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"sync"
)

// Storage — интерфейс для работы с хранилищем URL.
//   - GetShort — получить короткий идентификатор
//   - GetFull — получить полный URL
type Storage interface {
	GetShort(full string) (string, error)
	GetFull(short string) (string, error)
}

// Server — реализация интерфейса Storage.
//   - мапа хранения соответствий short → full URL.
//   - Мьютекс для синхронизации доступа к данным.
//   - Базовый URL
type Server struct {
	*sync.RWMutex                   // Позволяет использовать RWMutex через встраивание
	list          map[string]string // Хранилище: короткий ID → полный URL
	u             url.URL           // Базовый URL
}

// shortLen — длина генерируемых коротких идентификаторов (в символах).
const shortLen = 6

// New новый экземпляр сервера
func New() *Server {
	return &Server{
		RWMutex: &sync.RWMutex{},
		list:    make(map[string]string),
	}
}

// GetShort возвращает короткий идентификатор
// Если URL нет — генерирует новый
func (p *Server) GetShort(full string) (string, error) {
	p.Lock()
	defer p.Unlock()

	// Проверяем URL
	for k, v := range p.list {
		if v == full {
			return k, nil // Возвращаем ID
		}
	}

	// Генерируем новый короткий ID
	short, err := CreateShort(shortLen)
	if err != nil {
		return "", err
	}

	// Сохраняем
	p.list[short] = full
	return short, nil
}

// GetFull возвращает полный URL
func (p *Server) GetFull(short string) (string, error) {
	p.RLock()
	defer p.RUnlock()

	// Удаляем /
	short = strings.Trim(short, "/")

	// Проверяем наличие
	if full, ok := p.list[short]; ok {
		return full, nil // Возвращаем найденный URL
	}

	// Если не найдено — ошибка
	return "", fmt.Errorf("path %s not found", short)
}

// CreateShort генерирует строку длиной n.
func CreateShort(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Кодируем в base64 для URL
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}
