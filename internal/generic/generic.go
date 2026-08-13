package generic

import (
	"sync"
)

// ResetPool интерфейс c методом Reset.
type ResetPool interface {
	Reset()
}

// Pool структура для хранения данных.
type Pool[T any] struct {
	pool sync.Pool
}

// ResetPool сбросить данные.
func New[T ResetPool]() *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				var x T
				return x
			},
		},
	}
}

// Get получить данные.
func (p *Pool[T]) Get() T {
	if v := p.pool.Get(); v != nil {
		return v.(T)
	}
	var x T
	return x
}

// Put добавить данные.
func (p *Pool[T]) Put(obj T) {
	if x, ok := any(obj).(ResetPool); ok {
		x.Reset()
	}
	p.pool.Put(obj)
}
