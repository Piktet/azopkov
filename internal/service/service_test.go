package service

import (
	"context"
	"testing"

	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/Piktet/azopkov.git/internal/repository/connloader"
	"github.com/Piktet/azopkov.git/internal/repository/fileloader"
	"github.com/stretchr/testify/assert"
)

const defaultFileName = "../../../storage.json"

func TestServer_GetFull(t *testing.T) {

	full := "http://yandex.ru"
	server := New()

	loader := connloader.New("postgres://user:user@localhost:5433/practicum?sslmode=disable")
	server.Load(context.TODO(), loader)

	short, err := server.GetShort(context.TODO(), full)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		short   string
		want    string
		wantErr bool
	}{
		{
			name:    "positive",
			short:   short,
			want:    full,
			wantErr: false,
		},
		{
			name:    "negative",
			short:   full,
			want:    "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := server.GetFull(context.TODO(), test.short)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, test.want, got, "full")
				}
			}

		})
	}
}

func TestServer_GetShort(t *testing.T) {
	full := "http://yandex.ru"
	server := New()

	loader := connloader.New("postgres://user:user@localhost:5433/practicum?sslmode=disable")
	server.Load(context.TODO(), loader)

	short, err := server.GetShort(context.TODO(), full)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		full string
		want string
	}{
		{
			name: "positive",
			full: full,
			want: short,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := server.GetShort(context.TODO(), test.full)
			assert.NoError(t, err)
			assert.Equal(t, test.want, got, "full")
		})
	}
}

func TestServer_GetShortList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		fullList []model.FullItem
		wantErr  bool
	}{
		{
			fullList: []model.FullItem{{Corr: "aaa", Full: "full"}},
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New()
			got, err := p.GetShortList(context.Background(), test.fullList)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.NotEmpty(t, got)
				}
			}
		})
	}
}

func TestServer_Load(t *testing.T) {
	conn := connloader.New("")
	file := fileloader.New(defaultFileName)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		loader  model.StorageLoader
		wantErr bool
	}{
		{
			loader:  file,
			wantErr: false,
		},
		{
			loader:  &memLoader{},
			wantErr: false,
		},
		{
			loader:  conn,
			wantErr: true,
		},
		{
			loader:  nil,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New()
			err := p.Load(context.Background(), test.loader)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_checkFull(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		full    string
		wantErr bool
	}{
		{full: "http://ya.ru",
			wantErr: false,
		},
		{
			full:    "aaa",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkFull(test.full)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
