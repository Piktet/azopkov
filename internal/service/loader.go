package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Piktet/azopkov.git/internal/config"
	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/Piktet/azopkov.git/pkg/utils"
)

type memLoader struct{}

func (p *memLoader) Load(ctx context.Context) (map[string]string, error) {
	return map[string]string{}, nil
}

func (p *memLoader) GetShortList(ctx context.Context, fullList []model.FullItem, user string) (map[string]string, error) {
	resp := make(map[string]string)
	for _, v := range fullList {
		short, err := utils.CreateShort(config.ShortLen)
		if err != nil {
			return nil, err
		}
		resp[v.Full] = short
	}
	return resp, nil
}

func (p *memLoader) GetShort(ctx context.Context, full string, user string) (string, error) {
	short, err := utils.CreateShort(config.ShortLen)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (p *memLoader) GetFull(ctx context.Context, short string) (string, error) {
	return "", fmt.Errorf("full not found for short %s", short)
}

func (p *memLoader) GetUserList(ctx context.Context, user string) ([]model.StoreItem, error) {
	return nil, errors.New("unsupported")
}

func (p *memLoader) DeleteList(context.Context, []string, string) error {
	return errors.New("unsupport")
}

// GetStat получение статистики пользователей и адресов.
func (p *memLoader) GetStat(context.Context) (int, int, error) {
	return 0, 0, errors.New("unsupport")
}
