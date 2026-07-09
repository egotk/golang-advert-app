package userhttp

import (
	"time"

	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

type dtoResponse struct {
	ID          int64      `json:"id"`
	Version     int64      `json:"version"`
	Email       string     `json:"email"`
	FullName    string     `json:"full_name"`
	PhoneNumber string     `json:"phone_number"`
	Role        string     `json:"role"`
	LockedUntil *time.Time `json:"locked_until"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func dtoResponseFromEntity(u userentity.User) dtoResponse {
	return dtoResponse{
		ID:          u.ID,
		Version:     u.Version,
		Email:       u.Email,
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber,
		Role:        u.Role,
		LockedUntil: u.LockedUntil,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func dtoResponseFromEntities(users []userentity.User) []dtoResponse {
	res := make([]dtoResponse, len(users))

	for i, u := range users {
		res[i] = dtoResponseFromEntity(u)
	}

	return res
}
