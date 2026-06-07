package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWithAuth(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    http.HandlerFunc
		want int
	}{
		{
			name: "negative",
			h:    func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadRequest) },
			want: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := WithAuth(test.h)
			assert.ObjectsAreEqual(test.h, got)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			got.(http.HandlerFunc)(w, r)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, test.want, result.StatusCode, "status code")

		})
	}
}

func Test_newCookie(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    http.HandlerFunc
		want int
	}{{
		name: "positive",
		want: http.StatusBadRequest,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			newCookie(w, r)
		})
	}
}

func TestGetUser(t *testing.T) {

	user := uuid.New().String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
		User: user,
	})

	auth, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Error(err)
		return
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		auth    string
		user    string
		wantErr bool
	}{
		{
			name:    "negative",
			auth:    "auth",
			user:    "auth",
			wantErr: true,
		},
		{
			name:    "positive",
			auth:    auth,
			user:    user,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetUser(test.auth)
			// TODO: update the condition below to compare got with tt.want.
			if test.wantErr {
				assert.NotEqual(t, got, test.user)
			} else {
				assert.Equal(t, got, test.user)
			}
		})
	}
}
