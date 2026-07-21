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

func TestUseCase_ListFavourites(t *testing.T) {
	type listFavouritesMockBehavior func(mr *Mockrepo, dto advertusecase.ListDTO)

	negativeOffset := int64(-1)

	testTable := []struct {
		name           string
		inputDTO       advertusecase.ListDTO
		mockBehavior   listFavouritesMockBehavior
		expectedErrIs  error
		expectedCount  int64
		expectedLength int
	}{
		{
			name: "OK",
			inputDTO: advertusecase.ListDTO{
				UserID:   1,
				UserRole: "user",
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.ListDTO) {
				mr.EXPECT().CountFavourites(gomock.Any(), dto.UserID, gomock.Any()).Return(int64(1), nil)
				mr.EXPECT().ListFavourites(gomock.Any(), dto.UserID, dto.Limit, dto.Offset, gomock.Any()).Return(
					[]advertentity.Advert{{ID: 1}}, nil)
				mr.EXPECT().ListImagesByAdvertIDs(gomock.Any(), []int64{1}).Return(nil, nil)
			},
			expectedCount:  1,
			expectedLength: 1,
		},
		{
			name: "invalid dto",
			inputDTO: advertusecase.ListDTO{
				UserID: 1,
				Offset: &negativeOffset,
			},
			mockBehavior:  func(mr *Mockrepo, dto advertusecase.ListDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "repo error",
			inputDTO: advertusecase.ListDTO{
				UserID:   1,
				UserRole: "user",
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.ListDTO) {
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

			count, favs, err := uc.ListFavourites(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCount, count)
			assert.Len(t, favs, testCase.expectedLength)
		})
	}
}
