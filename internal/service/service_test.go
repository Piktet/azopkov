package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_GetFull(t *testing.T) {

	full := "http://yandex.ru"
	server := New()
	short, err := server.GetShort(full)
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

			got, err := server.GetFull(test.short)

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
	short, err := server.GetShort(full)
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
			got, err := server.GetShort(test.full)
			assert.NoError(t, err)
			assert.Equal(t, test.want, got, "full")
		})
	}
}
