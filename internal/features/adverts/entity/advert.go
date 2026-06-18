package advertentity

import "time"

type Status string

const (
	StatusInitial  Status = "initial"
	StatusActive   Status = "active"
	StatusRejected Status = "rejected"
	StatusBlocked  Status = "blocked"
	StatusArchived Status = "archived"

	initialID         = 0
	initialVersion    = 1
	initialViewsCount = 0
	initialStatus     = StatusInitial
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
