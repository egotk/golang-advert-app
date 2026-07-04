package usergrpc

import (
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userpb "github.com/egotk/golang-advert-app/internal/gen/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func userToResponse(u userentity.User) *userpb.UserResponse {
	var lockedUntil *timestamppb.Timestamp
	if u.LockedUntil != nil {
		lockedUntil = timestamppb.New(*u.LockedUntil)
	}

	return &userpb.UserResponse{
		Id:          u.ID,
		Version:     u.Version,
		Email:       u.Email,
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber,
		Role:        u.Role,
		LockedUntil: lockedUntil,
		CreatedAt:   timestamppb.New(u.CreatedAt),
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
		ImagePath:   u.ImagePath,
	}
}

func usersToResponse(users []userentity.User) *userpb.UsersResponse {
	userResponses := make([]*userpb.UserResponse, len(users))

	for i, u := range users {
		userResponses[i] = userToResponse(u)
	}

	return &userpb.UsersResponse{
		Users: userResponses,
	}
}
