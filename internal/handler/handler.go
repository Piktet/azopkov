package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/Piktet/azopkov.git/internal/service"
	"github.com/Piktet/azopkov.git/pkg/utils"
	"github.com/go-chi/chi/v5"
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

type ConnServer struct {
	model.ConnLoader
}

func NewConn(x model.ConnLoader) *ConnServer {
	return &ConnServer{ConnLoader: x}
}

// HTTP-сервер для сокращения URL.
type StorageServer struct {
	model.Storage          // соответствие short <-> full
	u     *url.URL 		   // URL (например, http://localhost:8080)
	audit model.Audit
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

func (p *StorageServer) SetLoader(loader model.Storage) {
	p.Storage = loader
}

func (p *StorageServer) SetAudit(audit model.Audit) {
	p.audit = audit
}

// format преобразует путь (например, "/EwHXdJfB") в полный URL.
// Используется для возврата клиенту сокращённого URL в виде строки.
//
// format("/xEwHXdJfByz") → "http://localhost:8080/EwHXdJfB"
func (p *StorageServer) format(path string) string {
	p.u.Path = path     // Устанавливаем путь
	return p.u.String() // Возвращаем строковое представление URL
}

func (p *StorageServer) sendAudit(ctx context.Context, action, user, address string) {
	if p.audit == nil {
		return
	}

	p.audit.Send(ctx, &model.AuditData{
		Created: time.Now().Unix(),
		Action:  action,
		User:    user,
		Address: address,
	})
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

	logger.Log().Info("HandlerPostFull")
	if r.Method != http.MethodPost {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeText {
		logger.Log().Error("error context type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log().Error("error getting request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	full := strings.TrimSpace(string(body))
	if _, err := url.ParseRequestURI(full); err != nil {
		logger.Log().Error("error parsing request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := getUser(r)
	short, shorterr := p.GetShort(r.Context(), full, user)
	if shorterr != nil && !errors.Is(shorterr, utils.ErrConflict) {
		if errors.Is(shorterr, model.ErrorDeleted) {
			logger.Log().Error("error getting short", zap.Error(err))
			w.WriteHeader(http.StatusGone)
			return
		}
		logger.Log().Error("error getting short", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set(model.HeaderContentType, model.ContentTypeText)
	if shorterr != nil {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(p.format(short)))

	p.sendAudit(r.Context(), model.ActionShorten, user, full)
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

	logger.Log().Info("HandlerGetFull")
	if r.Method != http.MethodGet {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user := getUser(r)
	full, err := p.GetFull(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrorDeleted) {
			logger.Log().Error("error getting full (is deleted)", zap.Error(err))
			w.WriteHeader(http.StatusGone)
			return
		}
		logger.Log().Error("error getting full", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set(model.HeaderLocation, full)
	w.WriteHeader(http.StatusTemporaryRedirect)

	p.sendAudit(r.Context(), model.ActionFollow, user, full)
}

// Эндпоинт с методом POST и путём /.
// Сервер принимает в теле запроса JSON URL как application/json
// и возвращает ответ с кодом 201 и сокращённым JSON URL как application/json.
func (p *StorageServer) HandlerPostFullJSON(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerPostFullJSON")

	if r.Method != http.MethodPost {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeJSON {
		logger.Log().Error("error context type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	var request model.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		logger.Log().Error("error decoding request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	full := strings.TrimSpace(string(request.Full))
	if _, err := url.ParseRequestURI(full); err != nil {
		logger.Log().Error("error parsing request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := getUser(r)
	short, shorterr := p.GetShort(r.Context(), full, user)
	if shorterr != nil && !errors.Is(shorterr, utils.ErrConflict) {
		logger.Log().Error("error getting short", zap.Error(shorterr))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := model.Response{
		Short: p.format(short),
	}

	enc, err := json.Marshal(response)
	if err != nil {
		logger.Log().Error("error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	if shorterr != nil {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(enc)
	p.sendAudit(r.Context(), model.ActionShorten, user, full)
}

func (p *ConnServer) HandlerGetPing(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerGetPing")
	if r.Method != http.MethodGet {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := p.Ping(r.Context()); err != nil {
		logger.Log().Error("error ping", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(model.HeaderContentType, model.ContentTypeText)
	w.WriteHeader(http.StatusOK)
}

// Эндпоинт /api/shorten/batch, принимающий в теле запроса множество URL для сокращения в формате json
func (p *StorageServer) HandlerPostBatch(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerPostBatch")
	if r.Method != http.MethodPost {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeJSON {
		logger.Log().Error("error context type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	var request []model.FullItem
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		logger.Log().Error("error decoding request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log().Info("request", zap.Int("count", len(request)))
	response, err := p.GetShortList(r.Context(), request, getUser(r))
	if err != nil {
		logger.Log().Error("error getting short", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log().Info("response", zap.Int("count", len(response)))

	for k, v := range response {
		response[k].Short = p.format(v.Short)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		logger.Log().Error("Error marshal JSON response = ", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Log().Info("response", zap.String("data", string(jsonResponse)))

	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)

}
