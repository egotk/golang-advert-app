package userentity

import "time"

type RefreshToken struct {
	Hash      string
	UserID    int64
	IssuedAt  time.Time
	ExpiresAt time.Time
}
