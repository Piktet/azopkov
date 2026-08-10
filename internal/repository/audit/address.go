package audit

import (
	"bytes"
	"context"
	"net/http"
	
	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"go.uber.org/zap"
)

const observerAddressName = "url"

type AddressObserver struct {
	address string
}

func NewAddressObserver(address string) *AddressObserver {
	return &AddressObserver{address: address}
}

func (p *AddressObserver) GetName() string {
	return observerAddressName
}

func (p *AddressObserver) Send(ctx context.Context, data []byte) error {

	r := bytes.NewReader(data)
	resp, err := http.Post(p.address, model.ContentTypeJSON, r)
	if err != nil {
		logger.Log().Error("send audit url", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	return nil
}