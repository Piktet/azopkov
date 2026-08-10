package service

import (
	"context"

	"github.com/Piktet/azopkov.git/internal/model"
	"github.com/Piktet/azopkov.git/internal/repository/connloader"
	"github.com/Piktet/azopkov.git/pkg/utils"
	"github.com/google/uuid"
)

func ExampleServer_GetFull() {
	val, _ := utils.CreateShort(6)
	full := "http://yandex.ru/" + val
	server := New()
	user := uuid.New().String()

	loader := connloader.New("postgres://user:user@localhost:5433/practicum?sslmode=disable")
	server.Load(context.TODO(), loader)

	short, err := server.GetShort(context.TODO(), full, user)
	if err != nil {
		return
	}
	_, _ = server.GetFull(context.TODO(), short)
}

func ExampleServer_GetShort() {
	val, _ := utils.CreateShort(6)
	full := "http://yandex.ru/" + val
	server := New()
	user := uuid.New().String()

	loader := connloader.New("postgres://user:user@localhost:5433/practicum?sslmode=disable")
	server.Load(context.TODO(), loader)

	_, err := server.GetShort(context.TODO(), full, user)
	if err != nil {
		return
	}

	_, _ = server.GetShort(context.TODO(), full, user)

}

func ExampleServer_GetShortList() {
	fullList := []model.FullItem{{Corr: "aaa", Full: "full"}}
	user := uuid.New().String()

	p := New()
	_, _ = p.GetShortList(context.Background(), fullList, user)

}
