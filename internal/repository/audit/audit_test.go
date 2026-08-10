package audit

import (
	"context"
	"testing"

	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewAuditEvent(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{
			name: "positive",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewAuditEvent()
			// TODO: update the condition below to compare got with tt.want.
			assert.NotEmpty(t, got)
		})
	}
}

func TestAuditEvent_Register(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		o    Observer
	}{
		{
			name: "positive_file",
			o:    NewFileObserver("file"),
		},
		{
			name: "positive_address",
			o:    NewAddressObserver("address"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewAuditEvent()
			if !assert.NotEmpty(t, p) {
				return
			}
			p.Register(test.o)
		})
	}
}

func TestAuditEvent_Send(t *testing.T) {
	p := NewAuditEvent()
	p.Register(NewFileObserver("file"))
	p.Register(NewAddressObserver("address"))

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		data    *model.AuditData
		wantErr bool
	}{
		{
			name:    "negative",
			data:    &model.AuditData{},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewAuditEvent()
			if !assert.NotEmpty(t, p) {
				return
			}
			gotErr := p.Send(context.Background(), test.data)
			if test.wantErr {
				assert.Error(t, gotErr)
				return
			}
			assert.NoError(t, gotErr)
		})
	}
}