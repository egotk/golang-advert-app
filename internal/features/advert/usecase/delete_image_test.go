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

func TestUseCase_DeleteImage(t *testing.T) {
	type deleteImageMockBehavior func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.DeleteImageDTO)

	testTable := []struct {
		name          string
		inputDTO      advertusecase.DeleteImageDTO
		mockBehavior  deleteImageMockBehavior
		expectAnyErr  bool
		expectedErrIs error
	}{
		{
			name: "OK",
			inputDTO: advertusecase.DeleteImageDTO{
				UserID:   1,
				UserRole: "user",
				ImageID:  1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.DeleteImageDTO) {
				mr.EXPECT().GetImageByID(gomock.Any(), dto.ImageID).Return(int64(1), advertentity.AdvertImage{
					ID:   1,
					Path: "path/to/photo.jpg",
				}, nil)
				mr.EXPECT().GetByID(gomock.Any(), int64(1)).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
				}, nil)
				mr.EXPECT().DeleteImageByID(gomock.Any(), dto.ImageID).Return(nil)
				ms.EXPECT().DeleteByPath("path/to/photo.jpg").Return(nil)
			},
		},
		{
			name: "invalid image id",
			inputDTO: advertusecase.DeleteImageDTO{
				UserID:   1,
				UserRole: "user",
				ImageID:  0,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.DeleteImageDTO) {},
			expectAnyErr: true,
		},
		{
			name: "forbidden",
			inputDTO: advertusecase.DeleteImageDTO{
				UserID:   2,
				UserRole: "user",
				ImageID:  1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.DeleteImageDTO) {
				mr.EXPECT().GetImageByID(gomock.Any(), dto.ImageID).Return(int64(1), advertentity.AdvertImage{
					ID:   1,
					Path: "path/to/photo.jpg",
				}, nil)
				mr.EXPECT().GetByID(gomock.Any(), int64(1)).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
				}, nil)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
		{
			name: "repo error",
			inputDTO: advertusecase.DeleteImageDTO{
				UserID:   1,
				UserRole: "user",
				ImageID:  1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.DeleteImageDTO) {
				mr.EXPECT().GetImageByID(gomock.Any(), dto.ImageID).Return(int64(0), advertentity.AdvertImage{}, coreerrors.ErrNotFound)
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

			err := uc.DeleteImage(ctx, testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			if testCase.expectAnyErr {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
		})
	}
}
