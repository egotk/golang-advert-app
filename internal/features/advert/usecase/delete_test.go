package advertusecase_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestUseCase_Delete(t *testing.T) {
	type deleteMockBehavior func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.DeleteDTO)

	testTable := []struct {
		name          string
		inputDTO      advertusecase.DeleteDTO
		mockBehavior  deleteMockBehavior
		expectedErrIs error
	}{
		{
			name: "OK",
			inputDTO: advertusecase.DeleteDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.DeleteDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
				}, nil)
				mr.EXPECT().ListImagesByAdvertID(gomock.Any(), int64(1)).Return(nil, nil)
				mr.EXPECT().DeleteByID(gomock.Any(), dto.AdvertID).Return(nil)
			},
		},
		{
			name: "forbidden",
			inputDTO: advertusecase.DeleteDTO{
				UserID:   2,
				UserRole: "user",
				AdvertID: 1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.DeleteDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
				}, nil)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
		{
			name: "repo error",
			inputDTO: advertusecase.DeleteDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.DeleteDTO) {
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
			ms := NewMockstorage(c)
			testCase.mockBehavior(mr, ms, testCase.inputDTO)

			uc := advertusecase.New(NewMockfavRepo(c), mr, ms)

			ctx := corezaplogger.ToContext(context.Background(), &corezaplogger.Logger{Logger: zap.NewNop()})

			err := uc.Delete(ctx, testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
		})
	}
}