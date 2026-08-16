// Package auth авторизация
package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/google/uuid"
)

const (
	expire = time.Hour * 12
	secret = "qwerty"
)

// Claims — JWT-токен с данными пользователя.
type Claims struct {
	jwt.RegisteredClaims
	User string
}

// WithAuth — HTTP-мидлварь для обработки аутентификации по JWT-куку.
// Если кука отсутствует, создаёт новый токен и записывает его в ответ.
// Если кука есть, проверяет токен и записывает имя пользователя во вторую куку.
func WithAuth(h http.Handler) http.Handler {
	authFn := func(w http.ResponseWriter, r *http.Request) {

		auth, err := r.Cookie(model.CookieAuth)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				newCookie(w, r)
				h.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Header.Set(model.HeaderAuth, auth.Value)

		user := GetUser(auth.Value)
		if user == "" {
			newCookie(w, r)
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:    model.CookieUser,
				Value:   user,
				Path:    "/",
				Expires: time.Now().Add(expire),
			})
		}
		h.ServeHTTP(w, r)

	}
	return http.HandlerFunc(authFn)
}

func newCookie(w http.ResponseWriter, r *http.Request) {
	user := uuid.New().String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
		User: user,
	})

	auth, err := token.SignedString([]byte(secret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     model.CookieAuth,
		Value:    auth,
		Path:     "/",
		Expires:  time.Now().Add(expire),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    model.CookieUser,
		Value:   user,
		Path:    "/",
		Expires: time.Now().Add(expire),
	})
	r.Header.Set(model.HeaderAuth, auth)
}

// GetUser извлекает имя пользователя из JWT-токена.
// Возвращает пустую строку, если токен невалиден.
func GetUser(auth string) string {
	claims := &Claims{}
	if token, err := jwt.ParseWithClaims(auth, claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}); err != nil || !token.Valid {
		return ""
	}

	return claims.User
}
