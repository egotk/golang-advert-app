package advertlocal

import (
	"fmt"
	"os"
)

func (s *Storage) DeleteByPath(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("delete image: %w", err)
	}

	return nil
}
