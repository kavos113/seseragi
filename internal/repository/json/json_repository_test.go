package json

import (
	"os"
	"testing"
)

func setupTestRepository(t *testing.T) *JsonRepository {
	t.Helper()

	tmpDir := t.TempDir()
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return NewJsonRepository(tmpDir)
}
