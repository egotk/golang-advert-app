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

func TestUseCase_Create(t *testing.T) {
	type createMockBehavior func(mr *Mockrepo, ms *Mockstorage)

	testTable := []struct {
		name          string
		inputDTO      advertusecase.CreateDTO
		mockBehavior  createMockBehavior
		expectedErrIs error
		checkResult   func(t *testing.T, advert advertentity.Advert)
	}{
		{
			name: "OK",
			inputDTO: advertusecase.CreateDTO{
				UserID:      1,
				Title:       "Title",
				Description: "Description",
				Price:       100,
				CategoryID:  1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage) {
				mr.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, advert *advertentity.Advert) error {
						advert.ID = 1

						return nil
					})
			},
			checkResult: func(t *testing.T, advert advertentity.Advert) {
				assert.Equal(t, int64(1), advert.ID)
				assert.Equal(t, int64(1), advert.UserID)
				assert.Equal(t, "Title", advert.Title)
				assert.Equal(t, advertentity.StatusInitial, advert.Status)
			},
		},
		{
			name: "invalid dto",
			inputDTO: advertusecase.CreateDTO{
				UserID:      1,
				Title:       "",
				Description: "Description",
				Price:       100,
				CategoryID:  1,
			},
			mockBehavior:  func(mr *Mockrepo, ms *Mockstorage) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "repo error",
			inputDTO: advertusecase.CreateDTO{
				UserID:      1,
				Title:       "Title",
				Description: "Description",
				Price:       100,
				CategoryID:  1,
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage) {
				mr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(coreerrors.ErrConflict)
			},
			expectedErrIs: coreerrors.ErrConflict,
		},
		{
			name: "storage save error",
			inputDTO: advertusecase.CreateDTO{
				UserID:      1,
				Title:       "Title",
				Description: "Description",
				Price:       100,
				CategoryID:  1,
				Images: []imageentity.Image{
					{Name: "photo.jpg", Extension: ".jpg", File: strings.NewReader("fake image bytes")},
				},
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage) {
				ms.EXPECT().Save(".jpg", gomock.Any()).Return("", coreerrors.ErrConflict)
			},
			expectedErrIs: coreerrors.ErrConflict,
		},
		{
			name: "saved images rolled back on repo error",
			inputDTO: advertusecase.CreateDTO{
				UserID:      1,
				Title:       "Title",
				Description: "Description",
				Price:       100,
				CategoryID:  1,
				Images: []imageentity.Image{
					{Name: "photo.jpg", Extension: ".jpg", File: strings.NewReader("fake image bytes")},
				},
			},
			mockBehavior: func(mr *Mockrepo, ms *Mockstorage) {
				ms.EXPECT().Save(".jpg", gomock.Any()).Return("path/to/photo.jpg", nil)
				mr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(coreerrors.ErrConflict)
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
			testCase.mockBehavior(mr, ms)

			uc := advertusecase.New(NewMockfavRepo(c), mr, ms)

			ctx := corezaplogger.ToContext(context.Background(), &corezaplogger.Logger{Logger: zap.NewNop()})

			advert, err := uc.Create(ctx, testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			testCase.checkResult(t, advert)
		})
	}
}
