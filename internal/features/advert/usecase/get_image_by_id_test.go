package advertusecase_test

import (
	"context"
	"io"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_GetImageByID(t *testing.T) {
	type getImageByIDMockBehavior func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.GetImageDTO)

	testTable := []struct {
		name          string
		inputDTO      advertusecase.GetImageDTO
		mockBehavior  getImageByIDMockBehavior
		expectedErrIs error
		expectedName  string
	}{
		{
			name: "OK",
			inputDTO: advertusecase.GetImageDTO{
				UserID:   1,
				UserRole: "user",
				ImageID:  1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.GetImageDTO) {
				mr.EXPECT().GetImageByID(gomock.Any(), dto.ImageID).Return(int64(1), advertentity.AdvertImage{
					ID:   1,
					Name: "photo.jpg",
					Path: "path/to/photo.jpg",
				}, nil)
				mr.EXPECT().GetByID(gomock.Any(), int64(1)).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
					Status: advertentity.StatusActive,
				}, nil)
				ms.EXPECT().GetByPath("path/to/photo.jpg").Return(io.NopCloser(strings.NewReader("fake image bytes")), nil)
			},
			expectedName: "photo.jpg",
		},
		{
			name: "forbidden",
			inputDTO: advertusecase.GetImageDTO{
				UserID:   2,
				UserRole: "user",
				ImageID:  1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.GetImageDTO) {
				mr.EXPECT().GetImageByID(gomock.Any(), dto.ImageID).Return(int64(1), advertentity.AdvertImage{
					ID:   1,
					Name: "photo.jpg",
					Path: "path/to/photo.jpg",
				}, nil)
				mr.EXPECT().GetByID(gomock.Any(), int64(1)).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
					Status: advertentity.StatusInitial,
				}, nil)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
		{
			name: "repo error",
			inputDTO: advertusecase.GetImageDTO{
				UserID:   1,
				UserRole: "user",
				ImageID:  1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.GetImageDTO) {
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

			rc, image, err := uc.GetImageByID(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedName, image.Name)
			assert.NotNil(t, rc)
			rc.Close()
		})
	}
}
