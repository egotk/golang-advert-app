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

func TestUseCase_List(t *testing.T) {
	type listMockBehavior func(mr *Mockrepo, dto advertusecase.ListDTO)

	negativeLimit := int64(-1)

	testTable := []struct {
		name           string
		inputDTO       advertusecase.ListDTO
		mockBehavior   listMockBehavior
		expectedErrIs  error
		expectedCount  int64
		expectedLength int
	}{
		{
			name: "OK",
			inputDTO: advertusecase.ListDTO{
				UserID:   1,
				UserRole: "admin",
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.ListDTO) {
				mr.EXPECT().Count(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				mr.EXPECT().List(gomock.Any(), dto.Limit, dto.Offset, gomock.Any()).Return(
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
				Limit:  &negativeLimit,
			},
			mockBehavior:  func(mr *Mockrepo, dto advertusecase.ListDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "repo error",
			inputDTO: advertusecase.ListDTO{
				UserID:   1,
				UserRole: "admin",
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.ListDTO) {
				mr.EXPECT().Count(gomock.Any(), gomock.Any()).Return(int64(0), coreerrors.ErrInvalidArgument)
			},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "non-admin is scoped to active adverts",
			inputDTO: advertusecase.ListDTO{
				UserID:   1,
				UserRole: "user",
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.ListDTO) {
				assertActiveScope := func(_ context.Context, filter advertentity.Filter) {
					if assert.NotNil(t, filter.Status) {
						assert.Equal(t, advertentity.StatusActive, *filter.Status)
					}
				}
				mr.EXPECT().Count(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, filter advertentity.Filter) (int64, error) {
						assertActiveScope(ctx, filter)

						return 0, nil
					})
				mr.EXPECT().List(gomock.Any(), dto.Limit, dto.Offset, gomock.Any()).DoAndReturn(
					func(ctx context.Context, limit, offset *int64, filter advertentity.Filter) ([]advertentity.Advert, error) {
						assertActiveScope(ctx, filter)

						return nil, nil
					})
				mr.EXPECT().ListImagesByAdvertIDs(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			expectedCount:  0,
			expectedLength: 0,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			testCase.mockBehavior(mr, testCase.inputDTO)

			uc := advertusecase.New(NewMockfavRepo(c), mr, NewMockstorage(c))

			count, adverts, err := uc.List(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCount, count)
			assert.Len(t, adverts, testCase.expectedLength)
		})
	}
}
