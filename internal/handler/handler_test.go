package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/stretchr/testify/assert"
)

const addr = "localhost:8080"

func TestStorageServer_HandlerPostFull(t *testing.T) {

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
	router := chi.NewRouter()
	router.Post("/", server.HandlerPostFull)
	router.Get("/{id}", server.HandlerGetFull)
	router.Post(`/api/shorten`, server.HandlerPostFullJSON)

	go http.ListenAndServe(addr, router)

	haveMethod := http.MethodPost
	haveBody := "http://ya.ru"
	haveContentType := model.ContentTypeText
	wantContentType := model.ContentTypeText

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
			r.Header.Add(model.HeaderContentType, test.have.contentType)
			w := httptest.NewRecorder()
			server.HandlerPostFull(w, r)

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

func TestStorageServer_HandlerGetFull(t *testing.T) {

	type have struct {
		method  string
		request string
	}
	type want struct {
		code     int
		location string
	}

	server := New(addr)
	router := chi.NewRouter()
	router.Post("/", server.HandlerPostFull)
	router.Get("/{id}", server.HandlerGetFull)
	router.Post(`/api/shorten`, server.HandlerPostFullJSON)

	go http.ListenAndServe(addr, router)

	haveMethod := http.MethodGet
	full := "http://yandex.ru"
	short, err := server.GetShort(full)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string
		have have
		want want
	}{
		{
			name: "positive",
			have: have{
				method:  haveMethod,
				request: "/" + short,
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: full,
			},
		},
		{
			name: "negative method",
			have: have{
				method:  http.MethodPost,
				request: "/" + short,
			},
			want: want{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
		{
			name: "negative id",
			have: have{
				method:  haveMethod,
				request: "/srftgnj",
			},
			want: want{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.have.method, test.have.request, nil)
			w := httptest.NewRecorder()

			ctx := chi.NewRouteContext()
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
			ctx.URLParams.Add("id", test.have.request)

			server.HandlerGetFull(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.code, result.StatusCode)
			assert.Equal(t, test.want.location, result.Header.Get("Location"))
		})
	}

}

func TestStorageServer_HandlerPostFullJSON(t *testing.T) {

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
	router := chi.NewRouter()
	router.Post("/", server.HandlerPostFull)
	router.Get("/{id}", server.HandlerGetFull)
	router.Post(`/api/shorten`, server.HandlerPostFullJSON)

	go http.ListenAndServe(addr, router)

	haveMethod := http.MethodPost
	haveBody := `{ "url": "http://ya.ru" }`
	haveContentType := model.ContentTypeJSON
	wantContentType := model.ContentTypeJSON

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
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
				contentType: "plain/text",
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
			r := httptest.NewRequest(test.have.method, `/api/shorten`, body)
			r.Header.Add(model.HeaderContentType, test.have.contentType)
			w := httptest.NewRecorder()
			server.HandlerPostFullJSON(w, r)

			result := w.Result()
			haveShort, _ := io.ReadAll(result.Body)
			defer result.Body.Close()

			assert.Equal(t, test.want.code, result.StatusCode)

			if result.StatusCode == http.StatusCreated {
				assert.Equal(t, test.want.contentType, result.Header.Get(model.HeaderContentType))

				var response model.Response
				err := json.Unmarshal(haveShort, &response)
				if assert.NoError(t, err) {
					_, err := url.ParseRequestURI(string(response.Short))
					assert.NoError(t, err)

				}
			}
		})
	}
}
