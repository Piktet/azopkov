package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetBaseAddress(t *testing.T) {

	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: defaultBaseAddress,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetBaseAddress())
		})
	}
}

func TestConfig_GetServerAddress(t *testing.T) {

	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: defaultServerAddress,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetServerAddress())
		})
	}
}

func TestConfig_GetLogLevel(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: defaultLogLevel,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetLogLevel())
		})
	}
}

func TestConfig_GetFileName(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetFileName())
		})
	}
}

func TestConfig_GetConnAddress(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetConnAddress())
		})
	}
}

func TestConfig_GetAuditFile(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetAuditFile())
		})
	}
}

func TestConfig_GetAuditAddress(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetAuditAddress())
		})
	}
}

func TestConfig_IsEnableHTTPS(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want bool
	}{
		{
			name: "positive",
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.IsEnableHTTPS())
		})
	}
}

func TestConfig_getString(t *testing.T) {
	cfg := New()
	tests := []struct {
		name  string
		value configType
		want  string
	}{
		{
			name:  "positive",
			value: configLogLevel,
			want:  defaultLogLevel,
		},
		{
			name:  "negative_name",
			value: "empty",
			want:  "",
		},
		{
			name:  "negative_bool",
			value: configEnableHTTPS,
			want:  "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, cfg.getString(test.value), test.want)
		})
	}
}

func TestConfig_getBool(t *testing.T) {
	cfg := New()
	tests := []struct {
		name  string
		value configType
		want  bool
	}{
		{
			name:  "positive",
			value: configEnableHTTPS,
			want:  false,
		},
		{
			name:  "negative_name",
			value: "empty",
			want:  false,
		},
		{
			name:  "negative_string",
			value: configLogLevel,
			want:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, cfg.getBool(test.value), test.want)
		})
	}
}
