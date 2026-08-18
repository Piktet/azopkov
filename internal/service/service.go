// Package service сервис скоращения адресов
package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/Piktet/azopkov.git/pkg/utils"

	"go.uber.org/zap"
)

// Server — реализация интерфейса Storage.
//   - мапа хранения соответствий short → full URL.
//   - Мьютекс для синхронизации доступа к данным.
//   - Базовый URL
type Server struct {
	*sync.RWMutex // Позволяет использовать RWMutex через встраивание
	shortList     map[string]string
	fullList      map[string]string
	loader        model.StorageLoader
}

// New новый экземпляр сервера
func New() *Server {
	return &Server{
		RWMutex:   &sync.RWMutex{},
		shortList: make(map[string]string),
		fullList:  make(map[string]string),
		loader:    &memLoader{},
	}
}

func (p *Server) Load(ctx context.Context, loader model.StorageLoader) error {
	if loader == nil {
		return errors.New("empty loader")
	}
	list, err := loader.Load(ctx)
	if err != nil {
		logger.Log().Error("Server.Load", zap.Error(err))
		return err
	}

	p.addList(list)

	p.loader = loader
	return nil
}

func checkFull(full string) error {
	if _, err := url.ParseRequestURI(full); err != nil {
		logger.Log().Error("checkFull", zap.Error(err))
		return err
	}
	return nil
}

func (p *Server) GetShortList(ctx context.Context, fullList []model.FullItem, user string) ([]model.ShortItem, error) {
	shortList := make([]model.ShortItem, 0, len(fullList))
	storeList := make([]model.FullItem, 0, len(fullList))
	for _, v := range fullList {

		if err := checkFull(v.Full); err != nil {
			return nil, err
		}

		if short, err := p.getShort(v.Full); err == nil {
			shortList = append(shortList, model.ShortItem{Corr: v.Corr, Short: short})
			continue
		}
		storeList = append(storeList, v)
	}

	newList, err := p.loader.GetShortList(ctx, storeList, user)

	if newList != nil {
		for _, v := range storeList {
			short, ok := newList[v.Full]
			if !ok {
				err = errors.Join(err, fmt.Errorf("short not created for full %s", v.Full))
			}
			p.addItem(v.Full, short)
			shortList = append(shortList, model.ShortItem{Corr: v.Corr, Short: short})
		}
	}

	return shortList, err
}

// GetShort возвращает короткий идентификатор
// Если URL нет — генерирует новый
func (p *Server) GetShort(ctx context.Context, full string, user string) (string, error) {

	logger.Log().Info("service.GetFull", zap.String("full", full))

	if err := checkFull(full); err != nil {
		return "", err
	}

	if short, err := p.getShort(full); err == nil {
		return short, utils.ErrConflict
	}

	// Значение не найдено в памяти. Берем его из хранилища и сохраняем в память
	short, err := p.loader.GetShort(ctx, full, user)
	if err != nil {
		return "", err
	}

	p.addItem(full, short)

	return short, nil
}

func (p *Server) GetFullList(ctx context.Context, shortList []model.ShortItem) ([]model.FullItem, error) {
	fullList := make([]model.FullItem, 0, len(shortList))
	for _, v := range shortList {
		full, err := p.GetFull(ctx, v.Short)
		if err != nil {
			return nil, err
		}
		fullList = append(fullList, model.FullItem{Corr: v.Corr, Full: full})
	}
	return fullList, nil
}

// GetFull возвращает полный URL
func (p *Server) GetFull(ctx context.Context, short string) (string, error) {

	logger.Log().Info("service.GetFull", zap.String("short", short))

	short = strings.Trim(short, "/")
	if full, err := p.getFull(short); err == nil {
		logger.Log().Info("service.GetFull from memory", zap.String("short", short), zap.String("full", full))
		return full, nil
	}

	// Значение не найдено в памяти. Берем его из хранилища.
	full, err := p.loader.GetFull(ctx, short)
	if err != nil {
		logger.Log().Info("service.GetFull get error", zap.Error(err))
		return "", err
	}

	logger.Log().Info("service.GetFull return full", zap.String("short", short), zap.String("full", full))
	return full, nil

}

func (p *Server) GetUserList(ctx context.Context, user string) ([]model.StoreItem, error) {
	return p.loader.GetUserList(ctx, user)

}

func (p *Server) DeleteList(ctx context.Context, short []string, user string) error {
	for _, v := range short {
		p.deleteShort(v)
	}
	return p.loader.DeleteList(ctx, short, user)
}

// GetStat получение статистики пользователей и подключений.
func (p *Server) GetStat(ctx context.Context) (int, int, error) {
	return p.loader.GetStat(ctx)
}
