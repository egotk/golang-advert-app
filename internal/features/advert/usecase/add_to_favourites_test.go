package advertusecase_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_AddToFavourites(t *testing.T) {
	type addToFavouritesMockBehavior func(mr *Mockrepo, mf *MockfavRepo, dto advertusecase.AddToFavouritesDTO)

	testTable := []struct {
		name          string
		inputDTO      advertusecase.AddToFavouritesDTO
		mockBehavior  addToFavouritesMockBehavior
		expectedErrIs error
	}{
		{
			name: "OK",
			inputDTO: advertusecase.AddToFavouritesDTO{
				AdvertID: 1,
				UserID:   2,
				UserRole: "user",
			},
			mockBehavior: func(mr *Mockrepo, mf *MockfavRepo, dto advertusecase.AddToFavouritesDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
					Status: advertentity.StatusActive,
				}, nil)
				mf.EXPECT().Add(gomock.Any(), dto.AdvertID, dto.UserID).Return(nil)
			},
		},
		{
			name: "invalid dto",
			inputDTO: advertusecase.AddToFavouritesDTO{
				AdvertID: 0,
				UserID:   2,
				UserRole: "user",
			},
			mockBehavior:  func(mr *Mockrepo, mf *MockfavRepo, dto advertusecase.AddToFavouritesDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "forbidden",
			inputDTO: advertusecase.AddToFavouritesDTO{
				AdvertID: 1,
				UserID:   2,
				UserRole: "user",
			},
			mockBehavior: func(mr *Mockrepo, mf *MockfavRepo, dto advertusecase.AddToFavouritesDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
					Status: advertentity.StatusInitial,
				}, nil)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
		{
			name: "repo error",
			inputDTO: advertusecase.AddToFavouritesDTO{
				AdvertID: 1,
				UserID:   2,
				UserRole: "user",
			},
			mockBehavior: func(mr *Mockrepo, mf *MockfavRepo, dto advertusecase.AddToFavouritesDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{}, coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			mf := NewMockfavRepo(c)
			testCase.mockBehavior(mr, mf, testCase.inputDTO)

			uc := advertusecase.New(mf, mr, NewMockstorage(c))

			err := uc.AddToFavourites(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
		})
	}
}
