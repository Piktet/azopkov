package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var log *zap.Logger = zap.NewNop()

func Log() *zap.Logger {
	return log
}

func InitLogger(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	log = zl
	return nil
}

func WithLogging(h http.Handler) http.Handler {

	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		log.Info("statistic",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", responseData.status), // получаем перехваченный код статуса ответа
			zap.Int("size", responseData.size),     // получаем перехваченный размер ответа
			zap.String("headers", fmt.Sprintf("%+v", responseData.headers)),
		)

	}
	return http.HandlerFunc(logFn)
}
