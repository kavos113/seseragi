package json

import "os"

type JsonRepository struct {
	RootDir string
}

func NewJsonRepository(rootDir string) *JsonRepository {
	os.MkdirAll(rootDir, os.ModePerm)

	return &JsonRepository{RootDir: rootDir}
}
