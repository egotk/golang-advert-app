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

func TestUseCase_Count(t *testing.T) {
	type countMockBehavior func(mr *Mockrepo, dto advertusecase.CountDTO)

	invalidStatus := advertentity.Status("unknown")

	testTable := []struct {
		name          string
		inputDTO      advertusecase.CountDTO
		mockBehavior  countMockBehavior
		expectedErrIs error
		expectedCount int64
	}{
		{
			name: "OK",
			inputDTO: advertusecase.CountDTO{
				UserID:   1,
				UserRole: "admin",
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.CountDTO) {
				mr.EXPECT().Count(gomock.Any(), gomock.Any()).Return(int64(5), nil)
			},
			expectedCount: 5,
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
				UserRole: "admin",
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.CountDTO) {
				mr.EXPECT().Count(gomock.Any(), gomock.Any()).Return(int64(0), coreerrors.ErrInvalidArgument)
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

			count, err := uc.Count(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCount, count)
		})
	}
}
