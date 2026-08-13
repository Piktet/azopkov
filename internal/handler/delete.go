package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"go.uber.org/zap"
)

// HandlerDelete — обработчик DELETE-запроса на пути "/api/user/urls".
//
// Принимает:
//   - Метод: DELETE
//   - Content-Type: application/json
//   - Тело: ["short_url_1", "short_url_2", ...]
//
// Возвращает:
//   - Код 202 Accepted — запрос принят в обработку
func (p *StorageServer) HandlerDelete(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerDelete")
	if r.Method != http.MethodDelete {
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

	var request []string
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		logger.Log().Error("error decoding request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log().Info("request", zap.Int("count", len(request)))
	go func() {
		if err := p.DeleteList(context.Background(), request, getUser(r)); err != nil {
			logger.Log().Error("error getting short", zap.Error(err))
		}
	}()

	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(http.StatusAccepted)

}
