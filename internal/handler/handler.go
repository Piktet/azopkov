package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/service"
	"go.uber.org/zap"
)

const (
	headerContentType = "Content-Type"
	headerLocation    = "Location"
	contentTypeText   = "text/plain"
	contextTypeJSON   = "application/json"
)

type Request struct {
	Full string `json:"url"`
}

type Response struct {
	Short string `json:"result"`
}

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

// Эндпоинт с методом POST и путём /.
// Сервер принимает в теле запроса JSON URL как application/json
// и возвращает ответ с кодом 201 и сокращённым JSON URL как application/json.
func (p *StorageServer) HandlerPostFullJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log().Debug("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(headerContentType)
	if contentType != contextTypeJSON {
		logger.Log().Debug("error contect type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	var request Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		logger.Log().Debug("error decoding request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	full := strings.TrimSpace(string(request.Full))
	if _, err := url.ParseRequestURI(full); err != nil {
		logger.Log().Debug("error parsing request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	short, err := p.GetShort(full)
	if err != nil {
		logger.Log().Debug("error getting short", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := Response{
		Short: p.format(short),
	}

	enc, err := json.Marshal(response)
	if err != nil {
		logger.Log().Debug("error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	w.Header().Set(headerContentType, contextTypeJSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(enc)
}
