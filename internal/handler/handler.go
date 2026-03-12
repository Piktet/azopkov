package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

const (
	n = 5
)

func HandlerPost(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// проверяем, что в запросе есть заголовок Content-Type : text/plain
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//сохраняем тело запроса в переменную body
	LngURL := r.Body
	if LngURL == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//генерируем ответ с коротким URL
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(ShortURL(LngURL)))
}

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// проверяем, что в запросе есть заголовок Content-Type : text/plain
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ShtURL := r.PathValue("id")
	// возвращаем ответ с кодом 307 и полным адресом
	http.Redirect(w, r, LongUrl(ShtURL), http.StatusTemporaryRedirect)
}

// функция генерация строки для URL длиной n
func RndShortCode(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

// функция которая будет принимать на вход URL и возвращать короткий URL
func ShortURL(LngURL string) (string, error) {
}

// функция которая будет принимать на вход короткий URL и возвращать полный URL
func LongUrl(ShtURL string) string {
}
