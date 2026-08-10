package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_new(t *testing.T) {

	ctx, fnTimeoutCancel := context.WithTimeoutCause(context.Background(), time.Second*2, errors.New("stop timeout test"))
	defer fnTimeoutCancel()
	tests := []struct {
		name    string // description of this test case
		wantErr bool
	}{
		{
			name:    "positive",
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := new(context.WithCancelCause(ctx))
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_run(t *testing.T) {

	ctx, fnCancel := context.WithTimeoutCause(context.Background(), time.Second*2, errors.New("test stop with timeout 2 seconds"))
	defer fnCancel()
	tests := []struct {
		name    string // description of this test case
		wantErr bool
	}{
		{
			name:    "positive",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			listner := &http.Server{
				Addr: ":8000",
			}
			run(ctx, listner)
		})
	}
}
