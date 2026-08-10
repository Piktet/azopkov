package audit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAddressObserver(t *testing.T) {
	tests := []struct {
		name      string // description of this test case
		auditName string
	}{
		{
			name:      "negative",
			auditName: "address",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewAddressObserver(test.name)
			// TODO: update the condition below to compare got with tt.want.
			assert.NotEmpty(t, got)
		})
	}
}

func TestAddressObserver_GetName(t *testing.T) {

	tests := []struct {
		name      string // description of this test case
		auditName string
		want      string
	}{
		{
			name:      "negative",
			auditName: "address",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewAddressObserver(test.name)
			if !assert.NotEmpty(t, p) {
				return
			}
			got := p.GetName()
			// TODO: update the condition below to compare got with tt.want.
			assert.Equal(t, got, observerAddressName)
		})
	}
}

func TestAddressObserver_Send(t *testing.T) {
	tests := []struct {
		name      string // description of this test case
		auditName string
		data      []byte
		wantErr   bool
	}{
		{
			name:      "negative",
			auditName: "address",
			wantErr:   true,
			data:      []byte("data"),
		},
		{
			name:      "negative",
			auditName: "address",
			wantErr:   true,
			data:      []byte("data"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewAddressObserver(test.name)
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
