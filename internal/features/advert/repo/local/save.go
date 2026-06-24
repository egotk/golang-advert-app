package advertlocal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func (s *Storage) Save(
	extension string,
	reader io.Reader,
) (string, error) {
	if err := os.MkdirAll(s.root, 0755); err != nil {
		return "", fmt.Errorf("mkdir advert_images: %w", err)
	}

	fileName := uuid.NewString() + extension
	path := filepath.Join(s.root, fileName)

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("copy: %w", err)
	}

	return path, nil
}
