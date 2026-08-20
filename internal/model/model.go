// Package model константы и структуры
package model

import (
	"context"
)

// HTTP-заголовки, используемые в API.
const (
	HeaderContentType     = "Content-Type"
	HeaderContentEncoding = "Content-Encoding"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderLocation        = "Location"
	HeaderAuth            = "Authorisation"
	HeaderRealIP          = "X-Real-IP"
	ContentTypeText       = "text/plain"
	ContentTypeJSON       = "application/json"
	ContentTypeHTML       = "text/html"

	CookieAuth = "Auth"
	CookieUser = "User"

	ContextValueUser = "User"
	ContextValueAuth = "authorization"
)

// Request представляет тело входного запроса на сокращение URL.
// generate:reset
type Request struct {
	Full string `json:"url"`
}

// Response представляет тело выходного ответа с сокращённым URL.
// generate:reset
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
	GetStat(context.Context) (int, int, error)
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
	GetStat(ctx context.Context) (int, int, error)

	Load(ctx context.Context, loader StorageLoader) error
}

// ShortItem представляет элемент ответа с коротким URL и корреляционным ID.
// generate:reset
type ShortItem struct {
	Corr  string `json:"correlation_id"`
	Short string `json:"short_url"`
}

// FullItem представляет элемент входного запроса с полным URL и корреляционным ID.
// generate:reset
type FullItem struct {
	Corr string `json:"correlation_id"`
	Full string `json:"original_url"`
}

// StoreItem представляет сохранённую запись сопоставления короткого и полного URL.
// generate:reset
type StoreItem struct {
	Short string `json:"short_url"`
	Full  string `json:"original_url"`
}

// AuditData содержит данные для записи в лог аудита.
// generate:reset
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

// Ответ GET /api/internal/stats
//
//	{
//	 "urls": <int>, // количество сокращённых URL в сервисе
//	 "users": <int> // количество пользователей в сервисе
//	}
type StatData struct {
	AddressCount int `json:"urls"`
	UserCount    int `json:"users"`
}
