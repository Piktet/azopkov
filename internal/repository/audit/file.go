package audit

import (
	"context"
	"os"
)

const observerFileName = "file"

type FileObserver struct {
	name string
}

func NewFileObserver(name string) *FileObserver {
	// TODO: проверить существование файла, попоробовать открыть его для записи
	return &FileObserver{name: name}
}

func (p *FileObserver) GetName() string {
	return observerFileName
}

func (p *FileObserver) Send(ctx context.Context, data []byte) error {

	f, err := os.OpenFile(p.name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	data = append(data, '\n')

	if _, err = f.Write(data); err != nil {
		return err
	}
	return nil
}
