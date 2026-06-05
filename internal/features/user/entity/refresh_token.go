package userentity

import "time"

type RefreshToken struct {
	Hash      string
	UserID    int
	IssuedAt  time.Time
	ExpiresAt time.Time
}
