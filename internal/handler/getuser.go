package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Piktet/azopkov.git/internal/auth"
	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"go.uber.org/zap"
)

func getUser(r *http.Request) string {
	token := r.Header.Get(model.HeaderAuth)
	return auth.GetUser(token)
}

// Эндпоинт /api/shorten/batch, принимающий в теле запроса множество URL для сокращения в формате json
func (p *StorageServer) HandlerGetUser(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerGetUser")
	if r.Method != http.MethodGet {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := p.GetUserList(r.Context(), getUser(r))
	if err != nil {
		logger.Log().Error("error getting short", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(response) == 0 {
		logger.Log().Error("error getting short", zap.Error(err))
		w.WriteHeader(http.StatusNoContent)
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
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
