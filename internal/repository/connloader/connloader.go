package connloader

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type ConnLoader struct {
	conn *sql.DB
	addr string
}

func New(addr string) *ConnLoader {
	return &ConnLoader{addr: addr}

}

func (p *ConnLoader) Create(ctx context.Context) error {

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

func (p *ConnLoader) Ping(ctx context.Context) error {

	if p.conn == nil {
		if err := p.Create(ctx); err != nil {
			return err
		}
	}
	return p.conn.PingContext(ctx)
}

func (p *ConnLoader) Load(ctx context.Context) (map[string]string, error) {

	if err := p.Ping(ctx); err != nil {
		return nil, err
	}

	if _, err := p.conn.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS t_data (
    		id SERIAL PRIMARY KEY,
    		s_full VARCHAR(1000) NOT NULL,
    		s_short VARCHAR(100) NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_data_full ON t_data(s_full);
		CREATE INDEX IF NOT EXISTS idx_data_short ON t_data(s_short); `); err != nil {

		return nil, err
	}

	rows, err := p.conn.QueryContext(ctx, "select s_full, s_short from t_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var full, short *string

	list := make(map[string]string)
	for rows.Next() {
		if err := rows.Scan(&full, &short); err != nil {
			return nil, err
		}
		if full != nil && short != nil {
			list[*short] = *full
		}
		fmt.Println(full, short)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (p *ConnLoader) Store(ctx context.Context, full, short string) error {

	if err := p.Ping(ctx); err != nil {
		return err
	}

	_, err := p.conn.ExecContext(ctx, "insert into t_data(s_full, s_short) values($1, $2)", full, short)
	return err
}
