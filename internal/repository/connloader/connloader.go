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

type Connector interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)

	PingContext(context.Context) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type ConnLoader struct {
	conn Connector
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

	db.SetMaxOpenConns(100) // Установить максимальное количество открытых соединений к базе данных
	db.SetMaxIdleConns(100) // Установить максимальное количество неактивных соединений в пуле

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

func (p *ConnLoader) GetShortList(ctx context.Context, fullList []model.FullItem, user string) (map[string]string, error) {

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

		short, err := db.GetShort(ctx, tx, full.Full, user)
		if err != nil {
			logger.Log().Error("error", zap.Error(err))
			tx.Rollback()
			return nil, err
		}
		shortList[full.Full] = short
	}
	return shortList, tx.Commit()
}

func (p *ConnLoader) GetShort(ctx context.Context, full string, user string) (string, error) {

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

	short, err := db.GetShort(ctx, tx, full, user)
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

func (p *ConnLoader) GetUserList(ctx context.Context, user string) ([]model.StoreItem, error) {
	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	return db.GetUser(ctx, p.conn, user)
}

func (p *ConnLoader) deleteList(ctx context.Context, short []string, user string) error {
	if err := p.Ping(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return err
	}

	return db.Delete(ctx, p.conn, short, user)
}

func (p *ConnLoader) DeleteList(ctx context.Context, short []string, user string) error {

	return p.deleteList(ctx, short, user)

	/*
		chShort := make(chan string, len(short))
		defer close(chShort)

		go p.deleteList(ctx, chShort, user)

		var wg sync.WaitGroup
		for _, v := range short {
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case chShort <- v:
					return
				case <-ctx.Done():
					return
				}
			}()
		}
		wg.Wait()
		return nil
	*/
}
