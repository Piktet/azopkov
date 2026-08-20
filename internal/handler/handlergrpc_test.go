package handler

import (
	"context"
	"testing"

	"github.com/Piktet/azopkov.git/internal/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_getUserFromMd(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{
			name:    "negative",
			ctx:     context.Background(),
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := getUserFromMd(test.ctx)
			if test.wantErr {
				assert.Error(t, gotErr)
				return
			}
			if !assert.NoError(t, gotErr) {
				return
			}
			assert.Equal(t, got, test.want)
		})
	}
}

func TestStorageServer_ListUserURLs(t *testing.T) {
	server := New(addr)
	tests := []struct {
		name    string
		want    *proto.UserURLsResponse
		ctx     context.Context
		wantErr bool
	}{
		{
			name:    "negative",
			wantErr: true,
			ctx:     context.Background(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := server.ListUserURLs(test.ctx, &emptypb.Empty{})
			if test.wantErr {
				assert.Error(t, gotErr)
				return
			}
			assert.NoError(t, gotErr)
			assert.Equal(t, got, test.want)
		})
	}
}

func TestStorageServer_ExpandURL(t *testing.T) {
	server := New(addr)
	data := "ya.ru"
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		address string
		ctx     context.Context
		x       *proto.URLExpandRequest
		want    *proto.URLExpandResponse
		wantErr bool
	}{
		{
			name:    "negative",
			wantErr: true,
			ctx:     context.Background(),
			x:       proto.URLExpandRequest_builder{Id: &data}.Build(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := server.ExpandURL(test.ctx, test.x)
			if test.wantErr {
				assert.Error(t, gotErr)
				return
			}
			assert.NoError(t, gotErr)
			assert.Equal(t, got, test.want)
		})
	}
}

func TestStorageServer_ShortenURL(t *testing.T) {
	server := New(addr)
	data := "ya.ru"
	tests := []struct {
		name    string // description of this test case
		address string
		ctx     context.Context
		x       *proto.URLShortenRequest
		want    *proto.URLShortenResponse
		wantErr bool
	}{
		{
			name:    "negative",
			wantErr: true,
			ctx:     context.Background(),
			x:       proto.URLShortenRequest_builder{Url: &data}.Build(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := server.ShortenURL(test.ctx, test.x)
			if test.wantErr {
				assert.Error(t, gotErr)
				return
			}
			assert.NoError(t, gotErr)
			assert.Equal(t, got, test.want)
		})
	}
}
