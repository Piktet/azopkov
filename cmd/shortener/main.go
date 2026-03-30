package main

import (
	"net/http"

	"github.com/Piktet/azopkov.git/internal/config"
	"github.com/Piktet/azopkov.git/internal/handler"
	"github.com/go-chi/chi/v5"
)

// addr — адрес
// Формат: "хост:порт" — localhost:8080.
//const addr = "localhost:8080"

func main() {
	// Создаём новый экземпляр HTTP-сервера
	//srv := handler.New(addr)
	config := config.New()
	srv := handler.New(config.GetBaseAddress())

	// Инициализируем роутер
	router := chi.NewRouter()

	// Регистрируем обработчики:
	// - POST / → создание короткого URL
	// - GET /{id} → редирект по ID
	router.Post(`/`, srv.HandlerPostFull)
	router.Get(`/{id}`, srv.HandlerGetFull)

	// Запускаем HTTP-сервер
	//panic при ошибке
	if err := http.ListenAndServe(config.GetServerAddress(), router); err != nil {
		panic(err)
	}
}
