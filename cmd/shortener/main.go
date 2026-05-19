package main

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	// "context"

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
	// Создаём новый экземпляр HTTP-сервера
	//srv := handler.New(addr)
	cfg := config.New()
	if err := logger.InitLogger("info"); err != nil {
		panic(err)
	}

	srv := handler.New(cfg.GetBaseAddress())
	var loader model.StorageLoader
	if cfg.GetConnAddress() != "" {
		loader = connloader.New(cfg.GetConnAddress())
		if err := srv.Load(context.Background(), loader); err != nil {
			logger.Log().Error("conn not loaded", zap.Error(err))
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

	conn, _ := loader.(model.ConnLoader)
	connServer := handler.NewConn(conn)

	// Инициализируем роутер
	router := chi.NewRouter()

	// Регистрируем обработчики:
	// - POST / → создание короткого URL
	// - GET /{id} → редирект по ID
	router.Post(`/`, logger.WithLogging(compress.WithCompress(srv.HandlerPostFull)))
	router.Get(`/{id}`, logger.WithLogging(compress.WithCompress(srv.HandlerGetFull)))
	router.Post(`/api/shorten`, logger.WithLogging(compress.WithCompress(srv.HandlerPostFullJSON)))
	router.Get(`/ping`, logger.WithLogging(compress.WithCompress(connServer.HandlerGetPing)))

	// Запускаем HTTP-сервер
	//panic при ошибке
	if err := http.ListenAndServe(cfg.GetServerAddress(), router); err != nil {
		panic(err)
	}
}
