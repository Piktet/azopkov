package model

import "context"

const (
	HeaderContentType     = "Content-Type"
	HeaderContentEncoding = "Content-Encoding"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderLocation        = "Location"
	ContentTypeText       = "text/plain"
	ContentTypeJSON       = "application/json"
	ContentTypeHTML       = "text/html"
)

type Request struct {
	Full string `json:"url"`
}

type Response struct {
	Short string `json:"result"`
}

type StorageLoader interface {
	Load(ctx context.Context) (map[string]string, error)
	GetShortList(ctx context.Context, fullList []FullItem) (map[string]string, error)
	GetShort(ctx context.Context, full string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)
}

type ConnLoader interface {
	Create(context.Context) error
	Ping(context.Context) error
}

// Storage — интерфейс для работы с хранилищем URL.
//   - GetShort — получить короткий идентификатор
//   - GetFull — получить полный URL
type Storage interface {
	GetShortList(ctx context.Context, full []FullItem) ([]ShortItem, error)
	GetShort(ctx context.Context, full string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)

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
	Short string `json:"short"`
	Full  string `json:"full"`
}
