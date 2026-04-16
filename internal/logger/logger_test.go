package logger

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/Piktet/azopkov.git/internal/handler"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		level   string
		wantErr bool
	}{
		{
			level:   "Info",
			wantErr: false,
		},
		{
			level:   "Aaa",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := InitLogger(test.level)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWithLogging(t *testing.T) {

	server := handler.New("localhost:8080")

	router := chi.NewRouter()

	router.Post(`/`, WithLogging(server.HandlerPostFull))
	router.Get(`/{id}`, WithLogging(server.HandlerGetFull))

	body := strings.NewReader("http://yandex.ru")
	contentType := "text/plain"

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    http.HandlerFunc
		want int
	}{
		{
			name: "positiveGetFull",
			h:    server.HandlerGetFull,
			want: http.StatusBadRequest,
		},
		{
			name: "positivePostFull",
			h:    server.HandlerPostFull,
			want: http.StatusCreated,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := WithLogging(test.h)
			assert.ObjectsAreEqual(test.h, got)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", body)
			r.Header.Add("Content-Type", contentType)
			got(w, r)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, test.want, result.StatusCode, "status code")

		})
	}
}

func Test_loggingResponseWriter_Write(t *testing.T) {
	//server := handler.New(config.DefaultConfig.GetBaseAddress())

	data := []byte("qwerty")
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		b       []byte
		want    int
		wantErr bool
	}{
		{
			b:       data,
			want:    len(data),
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			//r := httptest.NewRequest(http.MethodGet, "/", nil)

			lw := loggingResponseWriter{
				ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
				responseData:   &responseData{},
			}

			// server.HandlerGetFull(&lw, r)
			len, err := lw.Write(test.b)
			if assert.NoError(t, err) {
				assert.Equal(t, test.want, lw.responseData.size)
				assert.Equal(t, test.want, len)
			}
		})
	}
}

func Test_loggingResponseWriter_WriteHeader(t *testing.T) {
	//server := handler.New(config.DefaultConfig.GetBaseAddress())
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		status int
		want   int
	}{
		{
			name:   "positive",
			status: 200,
			want:   200,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			//r := httptest.NewRequest(http.MethodGet, "/", nil)

			lw := loggingResponseWriter{
				ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
				responseData:   &responseData{},
			}

			// server.HandlerGetFull(&lw, r)
			lw.WriteHeader(test.status)
			assert.Equal(t, test.want, lw.responseData.status)

		})
	}
}
