package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Piktet/azopkov.git/internal/auth"
	"github.com/Piktet/azopkov.git/internal/compress"
	"github.com/Piktet/azopkov.git/internal/config"
	"github.com/Piktet/azopkov.git/internal/handler"
	"github.com/Piktet/azopkov.git/internal/logger"
	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/Piktet/azopkov.git/internal/repository/audit"
	"github.com/Piktet/azopkov.git/internal/repository/connloader"
	"github.com/Piktet/azopkov.git/internal/repository/fileloader"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// addr — адрес
// Формат: "хост:порт" — localhost:8080.
// const addr = "localhost:8080"
const stopTimeout = 5 * time.Second

func main() {
	if err := run(context.WithCancelCause(context.Background())); err != nil {
		log.Fatalf("exist with error: %v", err)
	}
}

func run(ctx context.Context, fnCancel context.CancelCauseFunc) error {
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

	auditEvent := audit.NewAuditEvent()

	if cfg.GetAuditFile() != "" {
		auditEvent.Register(audit.NewFileObserver(cfg.GetAuditFile()))
	}

	if cfg.GetAuditAddress() != "" {
		auditEvent.Register(audit.NewAddressObserver(cfg.GetAuditFile()))
	}

	srv.SetAudit(auditEvent)

	router := chi.NewRouter()

	router.Use(logger.WithLogging)
	router.Use(compress.WithCompress)
	router.Use(auth.WithAuth)

	router.Mount("/debug", middleware.Profiler())
	router.Post("/", srv.HandlerPostFull)
	router.Post("/api/shorten", srv.HandlerPostFullJSON)
	router.Post("/api/shorten/batch", srv.HandlerPostBatch)
	router.Get("/{id}", srv.HandlerGetFull)
	router.Get("/ping", connServer.HandlerGetPing)
	router.Get("/api/user/urls", srv.HandlerGetUser)
	router.Delete("/api/user/urls", srv.HandlerDelete)

	if err := run(ctx, &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}); err != nil {
		fnCancel(err)
		return err
	}
	fnCancel(nil)

	return nil
}

func run(ctx context.Context, srv *http.Server) error {

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		select {
		case s := <-sigint:
			logger.Log().Info("stop with signal", zap.String("signal", s.String()))
		case <-ctx.Done():
			logger.Log().Info("stop with context", zap.Error(context.Cause(ctx)))
		}

		stopCtx, cancel := context.WithTimeoutCause(context.Background(), stopTimeout, fmt.Errorf("server Shutdown with timeout %v", stopTimeout))
		defer cancel()
		if err := srv.Shutdown(stopCtx); err != nil {
			logger.Log().Info("HTTP server shutdown", zap.Error(err))
		}
	}()
	// Запускаем HTTP-сервер
	//panic при ошибке
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log().Info("HTTP server ListenAndServe", zap.Error(err))
		return err
	}

	logger.Log().Info("exit")
	return nil
}
