package audit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileObserver(t *testing.T) {
	tests := []struct {
		name      string // description of this test case
		auditName string
	}{
		{
			name:      "positive",
			auditName: "file",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewFileObserver(test.name)
			// TODO: update the condition below to compare got with tt.want.
			assert.NotEmpty(t, got)
		})
	}
}

func TestFileObserver_GetName(t *testing.T) {

	tests := []struct {
		name      string // description of this test case
		auditName string
		want      string
	}{
		{
			name:      "positive",
			auditName: "file",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewFileObserver(test.name)
			if !assert.NotEmpty(t, p) {
				return
			}
			got := p.GetName()
			// TODO: update the condition below to compare got with tt.want.
			assert.Equal(t, got, observerFileName)
		})
	}
}

func TestFileObserver_Send(t *testing.T) {
	tests := []struct {
		name      string // description of this test case
		auditName string
		data      []byte
		wantErr   bool
	}{
		{
			name:      "positive",
			auditName: "file",
			wantErr:   false,
			data:      []byte("data"),
		},
		{
			name:      "positive",
			auditName: "file",
			wantErr:   false,
			data:      []byte("data"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewFileObserver(test.name)
			assert.NotEmpty(t, p)
			gotErr := p.Send(context.Background(), test.data)

			if test.wantErr {
				assert.Error(t, gotErr)
				return
			}
			assert.NoError(t, gotErr)
		})
	}
}
