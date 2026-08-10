package model

import (
	"context"
	"errors"
)

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

var ErrorDeleted = errors.New("item deleted")

type Request struct {
	Full string `json:"url"`
}

type Response struct {
	Short string `json:"result"`
}

type StorageLoader interface {
	Load(ctx context.Context) (map[string]string, error)
	GetShortList(ctx context.Context, fullList []FullItem, user string) (map[string]string, error)
	GetShort(ctx context.Context, full string, user string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)
	GetUserList(ctx context.Context, user string) ([]StoreItem, error)
	DeleteList(context.Context, []string, string) error
}

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

type ShortItem struct {
	Corr  string `json:"correlation_id"`
	Short string `json:"short_url"`
}

type FullItem struct {
	Corr string `json:"correlation_id"`
	Full string `json:"original_url"`
}

type StoreItem struct {
	Short string `json:"short_url"`
	Full  string `json:"original_url"`
}

type AuditData struct {
	Created int64  `json:"ts"`
	Action  string `json:"action"`
	User    string `json:"user_id"`
	Address string `json:"url"`
}

const (
	ActionShorten = "shorten"
	ActionFollow  = "follow"
)

type Audit interface {
	Send(context.Context, *AuditData) error
}
