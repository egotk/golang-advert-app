package advertentity

import "time"

type Status string

const (
	StatusInitial  Status = "initial"
	StatusActive   Status = "active"
	StatusRejected Status = "rejected"
	StatusBlocked  Status = "blocked"
	StatusArchived Status = "archived"
)

type Advert struct {
	ID          int
	Version     int
	UserID      int
	Title       string
	Description string
	Price       int
	CategoryID  int
	Status      Status
	ViewsCount  int
	Images      []AdvertImage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewInitial(
	userID int,
	title string,
	description string,
	price int,
	categoryID int,
) Advert {
	const (
		initialID         = -1
		initialVersion    = 1
		initialViewsCount = 0
		initialStatus     = StatusInitial
	)

	now := time.Now()

	return Advert{
		ID:          initialID,
		Version:     initialVersion,
		UserID:      userID,
		Title:       title,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
		Status:      initialStatus,
		ViewsCount:  initialViewsCount,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (a Advert) IsPublic() bool {
	if a.Status == StatusActive || a.Status == StatusArchived {
		return true
	}

	return false
}
