package advertentity

import "time"

type AdvertImage struct {
	ID        int
	Name      string
	Position  int
	Path      string
	CreatedAt time.Time
}

func NewAdvertImageUninitialized(
	name string,
	path string,
) AdvertImage {
	const initialID = -1

	return AdvertImage{
		ID:        initialID,
		Name:      name,
		Path:      path,
		CreatedAt: time.Now(),
	}
}
