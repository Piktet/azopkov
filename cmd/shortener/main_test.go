package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Piktet/azopkov.git/internal/config"
	"github.com/Piktet/azopkov.git/internal/handler"
	"github.com/Piktet/azopkov.git/internal/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_runSrv(t *testing.T) {

	ctx, fnTimeoutCancel := context.WithTimeoutCause(context.Background(), time.Second*10, errors.New("stop timeout test"))
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
			err := runSrv(context.WithCancelCause(ctx))
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_runHTTP(t *testing.T) {

	cfg, err := config.New()
	if err != nil {
		t.Error(err)
		return
	}

	srv := handler.New(cfg.GetBaseAddress())

	ctx, fnCancel := context.WithTimeoutCause(context.Background(), time.Second*2, errors.New("test stop with timeout 2 seconds"))
	defer fnCancel()
	tests := []struct {
		name        string // description of this test case
		enableHTTPS bool
		wantErr     bool
	}{
		{
			name:        "positive_http",
			enableHTTPS: false,
			wantErr:     true,
		},
		{
			name:        "positive_https",
			enableHTTPS: true,
			wantErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			grpcServer := grpc.NewServer()
			proto.RegisterShortenerServiceServer(grpcServer, srv)

			tlsConfig := &tls.Config{}

			if cfg.IsEnableHTTPS() {
				if cert, err := makeCertificate(); err == nil {
					tlsConfig.Certificates = cert
				}
			}

			httpServer := &http.Server{
				Addr:         cfg.GetServerAddress(),
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
				IdleTimeout:  120 * time.Second,
				TLSConfig:    tlsConfig,
			}

			go configureStop(ctx, grpcServer, httpServer)

			runHTTP(cfg.GetServerAddress(), httpServer)
		})
	}
}

func Test_runGrpc(t *testing.T) {

	cfg, err := config.New()
	if err != nil {
		t.Error(err)
		return
	}
	srv := handler.New(cfg.GetBaseAddress())

	ctx, fnCancel := context.WithTimeoutCause(context.Background(), time.Second*2, errors.New("test stop with timeout 2 seconds"))
	defer fnCancel()
	tests := []struct {
		name        string // description of this test case
		enableHTTPS bool
		wantErr     bool
	}{
		{
			name:        "positive_grpc",
			enableHTTPS: false,
			wantErr:     true,
		},
		{
			name:        "negative_grpc",
			enableHTTPS: true,
			wantErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			grpcServer := grpc.NewServer()
			proto.RegisterShortenerServiceServer(grpcServer, srv)

			tlsConfig := &tls.Config{}

			if cfg.IsEnableHTTPS() {
				if cert, err := makeCertificate(); err == nil {
					tlsConfig.Certificates = cert
				}
			}
			httpServer := &http.Server{
				Addr:         cfg.GetServerAddress(),
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
				IdleTimeout:  120 * time.Second,
				TLSConfig:    tlsConfig,
			}

			go configureStop(ctx, grpcServer, httpServer)

			runHTTP(cfg.GetServerAddress(), httpServer)
		})
	}
}
