package main

import (
	"net/http"

	"github.com/Piktet/azopkov.git/internal/handler"
)

// addr — адрес
// Формат: "хост:порт" — localhost:8080.
const addr = "localhost:8080"

func main() {
	// Создаём новый экземпляр HTTP-сервера
	srv := handler.New(addr)

	// Инициализируем роутер
	mux := http.NewServeMux()

	// Регистрируем обработчики:
	// - POST / → создание короткого URL
	// - GET /{id} → редирект по ID
	mux.HandleFunc("/", srv.HandlerPostFull)
	mux.HandleFunc("/{id}", srv.HandlerGetFull)

	// Запускаем HTTP-сервер
	//panic при ошибке
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
