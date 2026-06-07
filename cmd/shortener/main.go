package main

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/Piktet/azopkov.git/internal/compress"
	"github.com/Piktet/azopkov.git/internal/config"
	"github.com/Piktet/azopkov.git/internal/handler"
	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/Piktet/azopkov.git/internal/repository/connloader"
	"github.com/Piktet/azopkov.git/internal/repository/fileloader"
	"github.com/go-chi/chi/v5"
)

// addr — адрес
// Формат: "хост:порт" — localhost:8080.
//const addr = "localhost:8080"

func main() {
	ctx, fnCancel := context.WithCancelCause(context.Background())
	defer fnCancel(errors.New("exit"))
	run(ctx)
}

func run(ctx context.Context) {
	// Создаём новый экземпляр HTTP-сервера
	//srv := handler.New(addr)
	cfg := config.New()
	if err := logger.InitLogger("info"); err != nil {
		panic(err)
	}

	logger.Log().Info("config",
		zap.String("serverAddress", cfg.GetServerAddress()),
		zap.String("baseAddress", cfg.GetBaseAddress()),
		zap.String("logLevel", cfg.GetLogLevel()),
		zap.String("fileName", cfg.GetFileName()),
		zap.String("connAddress", cfg.GetConnAddress()),
	)

	srv := handler.New(cfg.GetBaseAddress())
	var loader model.StorageLoader
	var conn model.ConnLoader
	if cfg.GetConnAddress() != "" {
		loader = connloader.New(cfg.GetConnAddress())
		conn, _ = loader.(model.ConnLoader)
		if err := srv.Load(context.Background(), loader); err != nil {
			logger.Log().Info("conn not loaded", zap.Error(err))
			loader = nil
		} else {
			logger.Log().Info("conn storage usage")
		}

	}

	if loader == nil && cfg.GetFileName() != "" {
		loader = fileloader.New(cfg.GetFileName())

		if err := srv.Load(context.Background(), loader); err != nil {
			logger.Log().Error("file not loaded", zap.Error(err))
			loader = nil
		} else {
			logger.Log().Info("file storage usage")
		}
	}

	if loader == nil {
		logger.Log().Info("memory storage usage")
	}

	connServer := handler.NewConn(conn)

	router := chi.NewRouter()

	router.Use(logger.WithLogging)
	router.Use(compress.WithCompress)

	router.Post("/", srv.HandlerPostFull)
	router.Post("/api/shorten", srv.HandlerPostFullJSON)
	router.Post("/api/shorten/batch", srv.HandlerPostBatch)
	router.Get("/{id}", srv.HandlerGetFull)
	router.Get("/ping", connServer.HandlerGetPing)

	go func() {
		if err := http.ListenAndServe(cfg.GetServerAddress(), router); err != nil {
			panic(err)
		}
	}()
	// Запускаем HTTP-сервер
	//panic при ошибке
	logger.Log().Info("listen port", zap.String("serverAddress", cfg.GetServerAddress()))

	<-ctx.Done()
	logger.Log().Info("exit")
}
