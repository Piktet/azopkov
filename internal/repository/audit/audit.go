// Package audit передача сообщений аудита
package audit

import (
	"context"
	"encoding/json"

	"github.com/Piktet/azopkov.git/internal/model"
	"golang.org/x/sync/errgroup"
)

type Observer interface {
	GetName() string
	Send(context.Context, []byte) error
}

type AuditEvent struct {
	observer map[string]Observer
}

func NewAuditEvent() *AuditEvent {
	return &AuditEvent{
		observer: make(map[string]Observer),
	}
}

func (p *AuditEvent) Register(o Observer) {
	if p.observer == nil {
		p.observer = make(map[string]Observer)
	}

	p.observer[o.GetName()] = o
}

func (p *AuditEvent) Send(ctx context.Context, data *model.AuditData) error {

	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	wg, sendCtx := errgroup.WithContext(ctx)
	for _, v := range p.observer {
		wg.Go(func() error {
			return v.Send(sendCtx, resp)
		})
	}

	return wg.Wait()
}
