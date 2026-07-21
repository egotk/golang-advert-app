package advertusecase_test

import (
	"context"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	imageentity "github.com/egotk/golang-advert-app/internal/features/image/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestUseCase_CreateImages(t *testing.T) {
	type createImagesMockBehavior func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.CreateImagesDTO)

	testTable := []struct {
		name          string
		inputDTO      advertusecase.CreateImagesDTO
		mockBehavior  createImagesMockBehavior
		expectedErrIs error
		checkResult   func(t *testing.T, images []advertentity.AdvertImage)
	}{
		{
			name: "OK",
			inputDTO: advertusecase.CreateImagesDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 1,
				Images: []imageentity.Image{
					{Name: "photo.jpg", Extension: ".jpg", File: strings.NewReader("fake image bytes")},
				},
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.CreateImagesDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
				}, nil)
				ms.EXPECT().Save(".jpg", gomock.Any()).Return("path/to/photo.jpg", nil)
				mr.EXPECT().CreateImages(gomock.Any(), dto.AdvertID, gomock.Any()).Return(nil)
			},
			checkResult: func(t *testing.T, images []advertentity.AdvertImage) {
				assert.Len(t, images, 1)
				assert.Equal(t, "photo.jpg", images[0].Name)
				assert.Equal(t, "path/to/photo.jpg", images[0].Path)
			},
		},
		{
			name: "invalid dto",
			inputDTO: advertusecase.CreateImagesDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 0,
			},
			mockBehavior:  func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.CreateImagesDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "forbidden",
			inputDTO: advertusecase.CreateImagesDTO{
				UserID:   2,
				UserRole: "user",
				AdvertID: 1,
				Images: []imageentity.Image{
					{Name: "photo.jpg", Extension: ".jpg", File: strings.NewReader("fake image bytes")},
				},
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.CreateImagesDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
				}, nil)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
		{
			name: "repo error",
			inputDTO: advertusecase.CreateImagesDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 1,
				Images: []imageentity.Image{
					{Name: "photo.jpg", Extension: ".jpg", File: strings.NewReader("fake image bytes")},
				},
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.CreateImagesDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{}, coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
		{
			name: "empty images",
			inputDTO: advertusecase.CreateImagesDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 1,
			},
			mockBehavior:  func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.CreateImagesDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "too many images",
			inputDTO: advertusecase.CreateImagesDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 1,
				Images: []imageentity.Image{
					{Name: "photo.jpg", Extension: ".jpg", File: strings.NewReader("fake image bytes")},
				},
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.CreateImagesDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
					Images: make([]advertentity.AdvertImage, 10),
				}, nil)
			},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "saved images rolled back on repo error",
			inputDTO: advertusecase.CreateImagesDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 1,
				Images: []imageentity.Image{
					{Name: "photo.jpg", Extension: ".jpg", File: strings.NewReader("fake image bytes")},
				},
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage, dto advertusecase.CreateImagesDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
				}, nil)
				ms.EXPECT().Save(".jpg", gomock.Any()).Return("path/to/photo.jpg", nil)
				mr.EXPECT().CreateImages(gomock.Any(), dto.AdvertID, gomock.Any()).Return(coreerrors.ErrConflict)
				ms.EXPECT().DeleteByPath("path/to/photo.jpg").Return(nil)
			},
			expectedErrIs: coreerrors.ErrConflict,
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

			images, err := uc.CreateImages(ctx, testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			testCase.checkResult(t, images)
		})
	}
}
