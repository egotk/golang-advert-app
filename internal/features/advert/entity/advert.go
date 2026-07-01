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
	ID          int64
	Version     int64
	UserID      int64
	Title       string
	Description string
	Price       int64
	CategoryID  int64
	Status      Status
	ViewsCount  int64
	Images      []AdvertImage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewInitial(
	userID int64,
	title string,
	description string,
	price int64,
	categoryID int64,
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
