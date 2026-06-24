package imageentity

import (
	"fmt"
	"io"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
)

type Image struct {
	Name      string
	Extension string
	File      io.Reader
}

var extensionsByContentType = map[string]string{
	"image/webp": ".webp",
	"image/png":  ".png",
	"image/jpeg": ".jpeg",
}

func ParseExtension(contentType string) (string, error) {
	ext, valid := extensionsByContentType[contentType]
	if !valid {
		return "", fmt.Errorf(
			"%s is not supported image type: %w",
			contentType,
			coreerrors.ErrInvalidArgument,
		)
	}

	return ext, nil
}
