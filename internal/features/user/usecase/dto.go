package userusecase

type CreateDTO struct {
	Email       string `validate:"required,min=3,max=255,email"`
	FullName    string `validate:"required,min=3,max=100"`
	PhoneNumber string `validate:"required,min=4,max=20,phone_regex"`
	Password    string `validate:"required,min=8,max=64,bcrypt_password_byte_len"`
}

type LoginDTO struct {
	Email    string `validate:"required,min=3,max=255,email"`
	Password string `validate:"required,min=8,max=64,bcrypt_password_byte_len"`
}

type LoginResultDTO struct {
	UserID int64
	Tokens TokensDTO
}

type LogoutDTO struct {
	UserID       int64  `validate:"required,gt=0"`
	RefreshToken string `validate:"required,min=1,max=512"`
}

type RefreshTokensDTO struct {
	RefreshToken string `validate:"required,min=1,max=512"`
}

type TokensDTO struct {
	Access  string
	Refresh string
}
