package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
)

type StorageLoader interface {
	Load() (map[string]string, error) // return map [short string] full string
	Store(string, string) error       // store (full, short)
}

// Storage — интерфейс для работы с хранилищем URL.
//   - GetShort — получить короткий идентификатор
//   - GetFull — получить полный URL
type Storage interface {
	GetShort(full string) (string, error)
	GetFull(short string) (string, error)
	Load(loader StorageLoader) error
}

// Server — реализация интерфейса Storage.
//   - мапа хранения соответствий short → full URL.
//   - Мьютекс для синхронизации доступа к данным.
//   - Базовый URL
type Server struct {
	*sync.RWMutex // Позволяет использовать RWMutex через встраивание
	shortList     map[string]string
	fullList      map[string]string
	loader        StorageLoader
}

// shortLen — длина генерируемых коротких идентификаторов (в символах).
const shortLen = 6

// New новый экземпляр сервера
func New() *Server {
	return &Server{
		RWMutex:   &sync.RWMutex{},
		shortList: make(map[string]string),
		fullList:  make(map[string]string),
	}
}

func (p *Server) Load(loader StorageLoader) error {
	p.loader = loader
	list, err := loader.Load()
	if err != nil {
		return err
	}
	p.shortList = list

	for k, v := range p.shortList {
		p.fullList[v] = k
	}
	return nil
}

func (p *Server) store(full, short string) error {
	p.shortList[short] = full
	p.fullList[full] = short
	if p.loader != nil {
		if err := p.loader.Store(full, short); err != nil {
			return err
		}
	}
	return nil
}

// GetShort возвращает короткий идентификатор
// Если URL нет — генерирует новый
func (p *Server) GetShort(full string) (string, error) {
	p.Lock()
	defer p.Unlock()

	if short, ok := p.fullList[full]; ok {
		return short, nil
	}

	short, err := CreateShort(shortLen)
	if err != nil {
		return "", err
	}

	if err := p.store(full, short); err != nil {
		return "", err
	}
	p.shortList[short] = full
	return short, nil
}

// GetFull возвращает полный URL
func (p *Server) GetFull(short string) (string, error) {
	p.RLock()
	defer p.RUnlock()

	// Удаляем /
	short = strings.Trim(short, "/")

	// Проверяем наличие
	if full, ok := p.shortList[short]; ok {
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
