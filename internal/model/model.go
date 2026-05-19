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
	Load(ctx context.Context) (map[string]string, error)        // return map [short string] full string
	Store(ctx context.Context, full string, short string) error // store (full, short)
}

type Storage interface {
	GetShort(ctx context.Context, full string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)
	Load(ctx context.Context, loader StorageLoader) error
}

type ConnLoader interface {
	Create(context.Context) error
	Ping(context.Context) error
}
