package advertlocal

import (
	"fmt"
	"io"
	"os"
)

func (s *Storage) GetByPath(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}

	return file, nil
}
