package json

type JsonRepository struct {
	RootDir string
}

func NewJsonRepository(rootDir string) *JsonRepository {
	return &JsonRepository{RootDir: rootDir}
}
