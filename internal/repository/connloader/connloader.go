package connloader

import (
	"context"
	"database/sql"

	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/Piktet/azopkov.git/internal/repository/db"
	_ "github.com/jackc/pgx/v5/stdlib"

	"go.uber.org/zap"
)

type ConnLoader struct {
	conn *sql.DB
	addr string
}

func New(addr string) *ConnLoader {
	return &ConnLoader{addr: addr}

}

func (p *ConnLoader) Create(ctx context.Context) error {

	db, err := sql.Open("pgx", p.addr)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return err
	}
	if err := db.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return err
	}

	db.SetMaxOpenConns(2) // Установить максимальное количество открытых соединений к базе данных
	db.SetMaxIdleConns(2) // Установить максимальное количество неактивных соединений в пуле

	p.conn = db
	return nil
}

func (p *ConnLoader) Ping(ctx context.Context) error {

	if p.conn == nil {
		if err := p.Create(ctx); err != nil {
			logger.Log().Error("error", zap.Error(err))
			return err
		}
	}
	return p.conn.PingContext(ctx)
}

func (p *ConnLoader) Load(ctx context.Context) (map[string]string, error) {

	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	if err := db.Create(ctx, p.conn); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	return db.LoadList(ctx, p.conn)
}

func (p *ConnLoader) GetShortList(ctx context.Context, fullList []model.FullItem) (map[string]string, error) {

	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	// defer tx.Rollback()

	shortList := make(map[string]string)

	for _, full := range fullList {

		short, err := db.GetShort(ctx, tx, full.Full)
		if err != nil {
			logger.Log().Error("error", zap.Error(err))
			tx.Rollback()
			return nil, err
		}
		shortList[full.Full] = short
	}
	return shortList, tx.Commit()
}

func (p *ConnLoader) GetShort(ctx context.Context, full string) (string, error) {

	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}
	//defer tx.Rollback()

	short, err := db.GetShort(ctx, tx, full)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		tx.Rollback()
		return "", err
	}

	return short, tx.Commit()
}

func (p *ConnLoader) GetFull(ctx context.Context, short string) (string, error) {

	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}

	return db.GetFull(ctx, p.conn, short)
}
