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

	type have struct {
		method  string
		request string
	}
	type want struct {
		code     int
		location string
	}

	server := New(addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HandlerPostFull)
	mux.HandleFunc("/{id}", server.HandlerGetFull)

	go http.ListenAndServe(addr, mux)

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
				request: "/srftgnj/",
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
			mux.ServeHTTP(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want.code, result.StatusCode)
			assert.Equal(t, test.want.location, result.Header.Get("Location"))
		})
	}

}
