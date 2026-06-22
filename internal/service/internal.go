package service

import (
	"fmt"

	"github.com/Piktet/azopkov.git/internal/logger"
	"go.uber.org/zap"
)

func (p *Server) addItem(full, short string) {
	p.Lock()
	defer p.Unlock()

	logger.Log().Info("add item to mem", zap.String("full", full), zap.String("short", short))

	p.fullList[full] = short
	p.shortList[short] = full
}

func (p *Server) addList(list map[string]string) {
	p.Lock()
	defer p.Unlock()

	logger.Log().Info("add to mem list", zap.Int("count", len(list)))

	p.fullList = list

	for k, v := range p.shortList {
		p.shortList[v] = k
	}
}

func (p *Server) getShort(full string) (string, error) {
	p.RLock()
	defer p.RUnlock()

	if short, ok := p.fullList[full]; ok {
		return short, nil
	}

	return "", fmt.Errorf("short not found in mem for %s", full)

}

func (p *Server) getFull(short string) (string, error) {

	p.RLock()
	defer p.RUnlock()

	if full, ok := p.shortList[short]; ok {
		return full, nil
	}

	return "", fmt.Errorf("full not found in mem for %s", short)

}

func (p *Server) deleteShort(short string) {

	p.Lock()
	defer p.Unlock()

	delete(p.shortList, short)
	for k, v := range p.fullList {
		if v == short {
			delete(p.fullList, k)
			return
		}
	}
}
