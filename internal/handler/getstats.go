package handler

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"go.uber.org/zap"
)

// HandlerGetStat статистика по объектам
func (p *StorageServer) HandlerGetStat(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerGetStat")
	if r.Method != http.MethodGet {
		logger.Log().Error("error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	real := r.Header.Get(model.HeaderRealIP)

	if real == "" {
		logger.Log().Error("error header " + model.HeaderRealIP)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if p.subnet == "" {
		logger.Log().Error("error config value " + p.subnet)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !checkTrustedSubnet(real, p.subnet) {
		w.WriteHeader(http.StatusForbidden)
		return

	}

	countUser, countAddress, err := p.GetStat(r.Context())
	if err != nil {
		logger.Log().Error("error getting stat", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(model.StatData{
		AddressCount: countAddress,
		UserCount:    countUser,
	})
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

func checkTrustedSubnet(real, subnet string) bool {
	if real == "" || subnet == "" {
		return false
	}
	x := net.ParseIP(real)

	_, s, err := net.ParseCIDR(subnet)
	if err != nil {
		logger.Log().Error("Error parse CIDR", zap.Error(err))
		return false
	}
	return s.Contains(x)
}
