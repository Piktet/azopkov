package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Piktet/azopkov.git/internal/config"
	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/pkg/utils"

	"go.uber.org/zap"
)

type Connector interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func Create(ctx context.Context, conn Connector) error {

	_, err := conn.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS t_data (
    		id SERIAL PRIMARY KEY,
    		s_full VARCHAR(1000) NOT NULL,
    		s_short VARCHAR(100) NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_data_full ON t_data(s_full);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_data_short ON t_data(s_short);`)

	if err != nil {
		logger.Log().Error("error", zap.Error(err))
	}
	return err
}

func LoadList(ctx context.Context, conn Connector) (map[string]string, error) {

	rows, err := conn.QueryContext(ctx, "select s_full, s_short from t_data")
	if err != nil {
		logger.Log().Error("select list from db error", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var full, short *string

	list := make(map[string]string)
	for rows.Next() {
		if err := rows.Scan(&full, &short); err != nil {
			logger.Log().Error("select list from db error", zap.Error(err))
			return nil, err
		}
		if full != nil && short != nil {
			list[*short] = *full
		}
	}

	if err = rows.Err(); err != nil {
		logger.Log().Error("select list from db error", zap.Error(err))
		return nil, err
	}

	logger.Log().Info("select list from db", zap.Int("count", len(list)))

	return list, nil
}

func GetFull(ctx context.Context, conn Connector, short string) (string, error) {
	rows, err := conn.QueryContext(ctx, "select s_full from t_data where s_short = $1", short)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		err := fmt.Errorf("full name for %s not found", short)
		logger.Log().Info("error", zap.Error(err))
		return "", fmt.Errorf("full name for %s not found", short)
	}

	var full *string
	if err := rows.Scan(&full); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}

	if err = rows.Err(); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}

	if full == nil {
		err := fmt.Errorf("full name for %s is empty", short)
		logger.Log().Error("error", zap.Error(err))
		return "", fmt.Errorf("full name for %s is empty", short)
	}

	return *full, nil
}

func GetShort(ctx context.Context, conn Connector, full string) (string, error) {
	if short, err := getShort(ctx, conn, full); err == nil {
		logger.Log().Info("GetShort from db ok", zap.String("full", full), zap.String("short", short))
		return short, utils.ErrConflict
	}

	short, err := utils.CreateShort(config.ShortLen)
	if err != nil {
		logger.Log().Error("CreateShort error", zap.Error(err))
		return "", err
	}

	if err := store(ctx, conn, full, short); err != nil {
		logger.Log().Error("store to db error", zap.Error(err))
		return "", err
	}
	return short, nil
}

func store(ctx context.Context, conn Connector, full, short string) error {
	_, err := conn.ExecContext(ctx, "insert into t_data(s_full, s_short) values($1, $2)", full, short)
	if err != nil {
		logger.Log().Error("insert error", zap.Error(err))
		return err
	}
	logger.Log().Info("inserted to db", zap.String("full", full), zap.String("short", short))
	return nil
}

func getShort(ctx context.Context, conn Connector, full string) (string, error) {
	rows, err := conn.QueryContext(ctx, "select s_short from t_data where s_full = $1", full)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		err = fmt.Errorf("short name for %s not found", full)
		logger.Log().Info("error", zap.Error(err))
		return "", fmt.Errorf("short name for %s not found", full)
	}

	var short *string
	if err := rows.Scan(&short); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}

	if err = rows.Err(); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return "", err
	}

	if short == nil {
		err = fmt.Errorf("short name for %s is empty", full)
		logger.Log().Error("error", zap.Error(err))
		return "", fmt.Errorf("short name for %s is empty", full)
	}

	return *short, nil
}
