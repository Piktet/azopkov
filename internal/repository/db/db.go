package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Piktet/azopkov.git/internal/config"
	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
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
		CREATE UNIQUE INDEX IF NOT EXISTS idx_data_short ON t_data(s_short);
		ALTER TABLE IF EXISTS t_data ADD COLUMN IF NOT EXISTS u_user uuid;
		ALTER TABLE IF EXISTS t_data ADD COLUMN IF NOT EXISTS b_deleted boolean;`)

	if err != nil {
		logger.Log().Error("error", zap.Error(err))
	}
	return err
}

func LoadList(ctx context.Context, conn Connector) (map[string]string, error) {

	rows, err := conn.QueryContext(ctx, "select s_full, s_short from t_data where not coalesce(b_deleted, false)")
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
			list[*full] = *short
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
	rows, err := conn.QueryContext(ctx, "select s_full, coalesce(b_deleted, false) from t_data where s_short = $1", short)
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
	var deleted bool
	if err := rows.Scan(&full, &deleted); err != nil {
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

	if deleted {
		err := fmt.Errorf("full name for %s is deleted", short)
		logger.Log().Error("error", zap.Error(err))
		return "", model.ErrorDeleted

	}

	return *full, nil
}

func GetShort(ctx context.Context, conn Connector, full string, user string) (string, error) {
	if short, err := getShort(ctx, conn, full); err == nil {
		logger.Log().Info("GetShort from db ok", zap.String("full", full), zap.String("short", short))
		return short, utils.ErrConflict
	}

	short, err := utils.CreateShort(config.ShortLen)
	if err != nil {
		logger.Log().Error("CreateShort error", zap.Error(err))
		return "", err
	}

	if err := store(ctx, conn, full, short, user); err != nil {
		logger.Log().Error("store to db error", zap.Error(err))
		return "", err
	}
	return short, nil
}

func GetUser(ctx context.Context, conn Connector, user string) ([]model.StoreItem, error) {
	rows, err := conn.QueryContext(ctx, "select s_full, s_short from t_data where u_user = $1 and not coalesce(b_deleted, false)", user)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	resp := make([]model.StoreItem, 0)
	var full, short string
	for rows.Next() {
		if err := rows.Scan(&full, &short); err != nil {
			logger.Log().Error("error", zap.Error(err))
			return nil, err
		}
		resp = append(resp, model.StoreItem{Full: full, Short: short})
	}

	if err = rows.Err(); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func Delete(ctx context.Context, conn Connector, short []string, user string) error {
	_, err := conn.ExecContext(ctx, "update t_data set b_deleted = true where s_short = any($1) and u_user = $2", short, user)
	if err != nil {
		logger.Log().Error("delete error", zap.Error(err))
		return err
	}
	logger.Log().Error("delete Ok", zap.String("short", fmt.Sprintf("%v", short)))

	return nil
}

func store(ctx context.Context, conn Connector, full, short, user string) error {
	_, err := conn.ExecContext(ctx, "insert into t_data(s_full, s_short, u_user) values($1, $2, $3)", full, short, user)
	if err != nil {
		logger.Log().Error("insert error", zap.Error(err))
		return err
	}
	logger.Log().Info("inserted to db", zap.String("full", full), zap.String("short", short))
	return nil
}

func getShort(ctx context.Context, conn Connector, full string) (string, error) {
	rows, err := conn.QueryContext(ctx, "select s_short, coalesce(b_deleted, false) from t_data where s_full = $1", full)
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
	var deleted bool
	if err := rows.Scan(&short, &deleted); err != nil {
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

	if deleted {
		err = fmt.Errorf("short name for %s is deleted", full)
		logger.Log().Error("error", zap.Error(err))
		return "", model.ErrorDeleted
	}

	return *short, nil
}
