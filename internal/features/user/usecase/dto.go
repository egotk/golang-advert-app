package userusecase

type CreateDTO struct {
	Email       string
	FullName    string
	PhoneNumber string
	Password    string
}

type LoginDTO struct {
	Email    string
	Password string
}

type LoginResultDTO struct {
	UserID int
	Tokens TokensDTO
}

type LogoutDTO struct {
	UserID       int
	RefreshToken string
}

type RefreshTokensDTO struct {
	RefreshToken string
}

type TokensDTO struct {
	Access  string
	Refresh string
}
