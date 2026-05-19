package connloader

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type ConnLoader struct {
	conn *sql.DB
	addr string
}

func New(addr string) *ConnLoader {
	return &ConnLoader{addr: addr}

}

func (p *ConnLoader) Ping(ctx context.Context) error {
	if p.conn == nil {
		return errors.New("not inited")
	}
	return p.conn.PingContext(ctx)
}

func (p *ConnLoader) Load(ctx context.Context) error {

	db, err := sql.Open("postgres", p.addr)
	if err != nil {
		return err
	}
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	p.conn = db
	return nil
}
