package handler

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Piktet/azopkov.git/internal/service"
)

// HTTP-сервер для сокращения URL.
type StorageServer struct {
	service.Storage          // соответствие short <-> full
	u               *url.URL // URL (например, http://localhost:8080)
}

// New новый экземпляр сервера в формате "host:port".
// По умолчанию "localhost".
// При ошибке - panic-а
func New(address string) *StorageServer {
	u, err := url.Parse(address)
	if err != nil {
		panic(err)
	}

	return &(StorageServer{Storage: service.New(), u: u})
}

// format преобразует путь (например, "/EwHXdJfB") в полный URL.
// Используется для возврата клиенту сокращённого URL в виде строки.
//
// format("/xEwHXdJfByz") → "http://localhost:8080/EwHXdJfB"
func (p *StorageServer) format(path string) string {
	p.u.Path = path     // Устанавливаем путь
	return p.u.String() // Возвращаем строковое представление URL
}

// HandlerPostFull — обработчик POST-запросов на пути "/".
//
// Принимает:
//   - Метод: POST
//   - Content-Type: text/plain
//   - Тело запроса: строка - URL
//
// Возвращает:
//   - Код 201 Created
//   - Тело: сокращённый URL (http://localhost:8080/EwHXdJfB)
//   - Content-Type: text/plain
//
// Ошибки:
//   - 400 Bad Request
func (p *StorageServer) HandlerPostFull(w http.ResponseWriter, r *http.Request) {
	// Проверка HTTP-метода
	// if r.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// Проверка типа содержимого
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Чтение тела
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // закрытие

	// Очистка от пробелов и проверка URL
	fullURL := strings.TrimSpace(string(body))
	if _, err := url.ParseRequestURI(fullURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Получение короткого идентификатора
	short, err := p.GetShort(fullURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Формирование ответа
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(p.format(short)))
}

// HandlerGetFull — обработчик GET-запросов на пути "/{id}".
//
// Принимает:
//   - Метод: GET
//   - Путь: /{id} (id- EwHXdJfB)
//
// Возвращает:
//   - Код 307 Temporary Redirect
//   - Location: URL
//
// Ошибки:
//   - 400 Bad Request
func (p *StorageServer) HandlerGetFull(w http.ResponseWriter, r *http.Request) {
	// Проверка метода
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Извлечение идентификатора из пути
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Поиск полного URL
	full, err := p.GetFull(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Установка заголовка
	w.Header().Set("Location", full)

	// Отправка временного редиректа (307)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
