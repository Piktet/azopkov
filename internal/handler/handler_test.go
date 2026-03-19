package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const addr = "localhost:8080"

func TestHandlerPostFull(t *testing.T) {

	type have struct {
		method      string
		contentType string
		body        string
	}
	type want struct {
		code        int
		contentType string
	}

	server := New(addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HandlerPostFull)
	mux.HandleFunc("/{id}", server.HandlerGetFull)

	go http.ListenAndServe(addr, mux)

	haveMethod := http.MethodPost
	haveBody := "http://ya.ru"
	haveContentType := "text/plain"
	wantContentType := "text/plain"

	tests := []struct {
		name string
		have have
		want want
	}{
		{
			name: "positive",
			have: have{
				method:      haveMethod,
				contentType: haveContentType,
				body:        haveBody,
			},
			want: want{
				code:        http.StatusCreated,
				contentType: wantContentType,
			},
		},
		{
			name: "negative method",
			have: have{
				method:      http.MethodGet,
				contentType: haveContentType,
				body:        haveBody,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: wantContentType,
			},
		},
		{
			name: "negative contentType",
			have: have{
				method:      haveMethod,
				contentType: "application/json",
				body:        haveBody,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: wantContentType,
			},
		},
		{
			name: "negative body",
			have: have{
				method:      haveMethod,
				contentType: haveContentType,
				body:        "http//ya.ru",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: wantContentType,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := strings.NewReader(test.have.body)
			r := httptest.NewRequest(test.have.method, `/`, body)
			r.Header.Add("Content-Type", test.have.contentType)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)

			result := w.Result()
			haveShort, _ := io.ReadAll(result.Body)
			defer result.Body.Close()

			assert.Equal(t, test.want.code, result.StatusCode)

			if result.StatusCode == http.StatusCreated {
				assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
				_, err := url.ParseRequestURI(string(haveShort))
				assert.NoError(t, err)
			}
		})
	}

}

//-----

func TestHandlerGetFull(t *testing.T) {
	const addr = "localhost:8080"

	server := New(addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HandlerPostFull)
	mux.HandleFunc("/{id}", server.HandlerGetFull)

	// Предварительно добавим один URL вручную через POST
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // блокируем редирект, чтобы проверить статус
		},
	}

	// Создаём тестовый сервер
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Вспомогательная функция для создания короткой ссылки
	createShort := func(fullURL string) (string, error) {
		resp, err := http.Post(ts.URL, "text/plain", strings.NewReader(fullURL))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return string(body), nil
	}

	tests := []struct {
		name           string
		setup          func() string // возвращает shortID
		id             string
		wantStatus     int
		wantLocation   string // ожидаемый заголовок Location при редиректе
		wantNoLocation bool   // если не должен быть Location (ошибки)
	}{
		{
			name: "positive redirect",
			setup: func() string {
				short, _ := createShort("https://ya.ru")
				return short
			},
			id:           "", // будет заполнено setup'ом
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: "https://ya.ru",
		},
		{
			name: "negative not found",
			setup: func() string {
				return "unknown1" // ID не существует
			},
			id:             "unknown1",
			wantStatus:     http.StatusNotFound,
			wantNoLocation: true,
		},
		{
			name: "negative empty id",
			setup: func() string {
				return ""
			},
			id:             "",
			wantStatus:     http.StatusBadRequest,
			wantNoLocation: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var shortID string
			if test.setup != nil {
				shortID = test.setup()
			}
			if test.id == "" && test.setup != nil {
				test.id = shortID
			}

			req, err := http.NewRequest(http.MethodGet, ts.URL+"/"+test.id, nil)
			assert.NoError(t, err)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, test.wantStatus, resp.StatusCode)

			if test.wantLocation != "" {
				location := resp.Header.Get("Location")
				assert.Equal(t, test.wantLocation, location)
			}
			if test.wantNoLocation {
				location := resp.Header.Get("Location")
				assert.Empty(t, location)
			}
		})
	}
}
