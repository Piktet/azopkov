// Package handler обработчики запросов
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/Piktet/azopkov.git/internal/proto"
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

// Request представляет тело входного запроса на сокращение URL.
type Request struct {
	Full string `json:"url"`
}

// Response представляет тело выходного ответа с сокращённым URL.
type Response struct {
	Short string `json:"result"`
}

// ConnServer — HTTP-сервер для работы с подключением к базе данных (ping и т.д.).
type ConnServer struct {
	model.ConnLoader
}

// NewConn создаёт новый экземпляр ConnServer с заданным подключением к БД.
func NewConn(x model.ConnLoader) *ConnServer {
	return &ConnServer{ConnLoader: x}
}

// StorageServer — HTTP-сервер для сокращения URL.
//   - Storage — соответствие short <-> full URL.
//   - u — базовый URL (например, http://localhost:8080).
//   - audit — интерфейс для отправки логов аудита.
type StorageServer struct {
	proto.UnimplementedShortenerServiceServer
	model.Storage
	u      *url.URL
	audit  model.Audit
	subnet string
}

// New создаёт новый экземпляр StorageServer с заданным базовым адресом.
// По умолчанию "localhost".
// При ошибке парсинга адреса — panic.
func New(address string) *StorageServer {
	u, err := url.Parse(address)
	if err != nil {
		panic(err)
	}

	return &(StorageServer{Storage: service.New(), u: u})
}

// SetLoader устанавливает реализацию Storage для сервера.
func (p *StorageServer) SetLoader(loader model.Storage) {
	p.Storage = loader
}

// SetAudit устанавливает интерфейс аудита для сервера.
func (p *StorageServer) SetAudit(audit model.Audit) {
	p.audit = audit
}

// SetTrustedSubnet установка места отправки аудита.
func (p *StorageServer) SetTrustedSubnet(subnet string) {
	p.subnet = subnet
}

func (p *StorageServer) format(path string) string {
	p.u.Path = path
	return p.u.String()
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log().Error("error getting request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	full := strings.TrimSpace(string(body))
	if _, err = url.ParseRequestURI(full); err != nil {
		logger.Log().Error("error parsing request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := getUser(r)
	short, shorterr := p.GetShort(r.Context(), full, user)
	if shorterr != nil && !errors.Is(shorterr, utils.ErrConflict) {
		if errors.Is(shorterr, model.ErrorDeleted) {
			logger.Log().Error("error getting short", zap.Error(shorterr))
			w.WriteHeader(http.StatusGone)
			return
		}
		logger.Log().Error("error getting short", zap.Error(shorterr))
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
//   - Путь: /{id} (id - EwHXdJfB)
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

// HandlerPostFullJSON — обработчик POST-запросов на пути "/api/shorten" в формате JSON.
//
// Принимает:
//   - Метод: POST
//   - Content-Type: application/json
//   - Тело: {"url": "http://..."}
//
// Возвращает:
//   - Код 201 Created / 409 Conflict
//   - Тело: {"result": "http://..."}
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

// HandlerGetPing — обработчик GET-запроса на пути "/ping" для проверки подключения к БД.
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

// HandlerPostBatch — обработчик POST-запроса на пути "/api/shorten/batch" для пакетного сокращения URL.
//
// Принимает:
//   - Метод: POST
//   - Content-Type: application/json
//   - Тело: [{"correlation_id": "...", "original_url": "..."}, ...]
//
// Возвращает:
//   - Код 201 Created
//   - Тело: [{"correlation_id": "...", "short_url": "..."}, ...]
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
