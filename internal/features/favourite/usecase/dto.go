package favusecase

type RemoveDTO struct {
	AdvertID int64 `validate:"gt=0"`
	UserID   int64
}
