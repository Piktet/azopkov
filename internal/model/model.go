// Package model константы и структуры
package model

import (
	"context"
	"errors"
)

// HTTP-заголовки, используемые в API.
const (
	HeaderContentType     = "Content-Type"
	HeaderContentEncoding = "Content-Encoding"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderLocation        = "Location"
	HeaderAuth            = "Authorisation"
	ContentTypeText       = "text/plain"
	ContentTypeJSON       = "application/json"
	ContentTypeHTML       = "text/html"

	CookieAuth = "Auth"
	CookieUser = "User"

	ContextValueUser = "User"
)

// ErrorDeleted возвращается при попытке получить запись, которая была удалена.
var ErrorDeleted = errors.New("item deleted")

// Request представляет тело входного запроса на сокращение URL.
type Request struct {
	Full string `json:"url"`
}

// Response представляет тело выходного ответа с сокращённым URL.
type Response struct {
	Short string `json:"result"`
}

// StorageLoader — интерфейс для загрузки данных из внешнего хранилища (файл, БД и т.д.).
// Реализации позволяют прочитать существующие сопоставления URL и создать новые.
type StorageLoader interface {
	Load(ctx context.Context) (map[string]string, error)
	GetShortList(ctx context.Context, fullList []FullItem, user string) (map[string]string, error)
	GetShort(ctx context.Context, full string, user string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)
	GetUserList(ctx context.Context, user string) ([]StoreItem, error)
	DeleteList(context.Context, []string, string) error
}

// ConnLoader — интерфейс для работы с подключением к базе данных.
type ConnLoader interface {
	Create(context.Context) error
	Ping(context.Context) error
}

// Storage — интерфейс для работы с хранилищем URL.
//   - GetShort — получить короткий идентификатор
//   - GetFull — получить полный URL
type Storage interface {
	GetUserList(ctx context.Context, user string) ([]StoreItem, error)
	GetShortList(ctx context.Context, full []FullItem, user string) ([]ShortItem, error)
	GetShort(ctx context.Context, full string, user string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)
	DeleteList(ctx context.Context, short []string, user string) error

	Load(ctx context.Context, loader StorageLoader) error
}

// ShortItem представляет элемент ответа с коротким URL и корреляционным ID.
type ShortItem struct {
	Corr  string `json:"correlation_id"`
	Short string `json:"short_url"`
}

// FullItem представляет элемент входного запроса с полным URL и корреляционным ID.
type FullItem struct {
	Corr string `json:"correlation_id"`
	Full string `json:"original_url"`
}

// StoreItem представляет сохранённую запись сопоставления короткого и полного URL.
type StoreItem struct {
	Short string `json:"short_url"`
	Full  string `json:"original_url"`
}

// AuditData содержит данные для записи в лог аудита.
type AuditData struct {
	Created int64  `json:"ts"`
	Action  string `json:"action"`
	User    string `json:"user_id"`
	Address string `json:"url"`
}

// Константы действий аудита.
const (
	ActionShorten = "shorten"
	ActionFollow  = "follow"
)

// Audit — интерфейс для отправки данных в систему аудита.
type Audit interface {
	Send(context.Context, *AuditData) error
}
