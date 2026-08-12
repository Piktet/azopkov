package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_runSrv(t *testing.T) {

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
			runSrv(ctx, listner)
		})
	}
}
