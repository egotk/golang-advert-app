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

func TestUseCase_CountFavourites(t *testing.T) {
	type countFavouritesMockBehavior func(mr *Mockrepo, dto advertusecase.CountDTO)

	invalidStatus := advertentity.Status("unknown")

	testTable := []struct {
		name          string
		inputDTO      advertusecase.CountDTO
		mockBehavior  countFavouritesMockBehavior
		expectedErrIs error
		expectedCount int64
	}{
		{
			name: "OK",
			inputDTO: advertusecase.CountDTO{
				UserID:   1,
				UserRole: "user",
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.CountDTO) {
				mr.EXPECT().CountFavourites(gomock.Any(), dto.UserID, gomock.Any()).Return(int64(3), nil)
			},
			expectedCount: 3,
		},
		{
			name: "invalid dto",
			inputDTO: advertusecase.CountDTO{
				UserID: 1,
				Filter: advertentity.Filter{Status: &invalidStatus},
			},
			mockBehavior:  func(mr *Mockrepo, dto advertusecase.CountDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "repo error",
			inputDTO: advertusecase.CountDTO{
				UserID:   1,
				UserRole: "user",
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.CountDTO) {
				mr.EXPECT().CountFavourites(gomock.Any(), dto.UserID, gomock.Any()).Return(int64(0), coreerrors.ErrInvalidArgument)
			},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			testCase.mockBehavior(mr, testCase.inputDTO)

			uc := advertusecase.New(NewMockfavRepo(c), mr, NewMockstorage(c))

			count, err := uc.CountFavourites(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCount, count)
		})
	}
}
