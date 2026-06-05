package corejwt

import "time"

type RefreshToken struct {
	Token     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}
