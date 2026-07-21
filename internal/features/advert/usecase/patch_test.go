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

func TestUseCase_Patch(t *testing.T) {
	type patchMockBehavior func(mr *Mockrepo, dto advertusecase.PatchDTO)

	newTitle := "New Title"
	whitespaceTitle := "   "

	testTable := []struct {
		name          string
		inputDTO      advertusecase.PatchDTO
		mockBehavior  patchMockBehavior
		expectedErrIs error
		checkResult   func(t *testing.T, advert advertentity.Advert)
	}{
		{
			name: "OK",
			inputDTO: advertusecase.PatchDTO{
				UserID:   1,
				UserRole: "user",
				ID:       1,
				Version:  1,
				Title:    &newTitle,
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
					Status: advertentity.StatusActive,
				}, nil)
				mr.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResult: func(t *testing.T, advert advertentity.Advert) {
				assert.Equal(t, "New Title", advert.Title)
				assert.Equal(t, advertentity.StatusInitial, advert.Status)
			},
		},
		{
			name: "empty patch",
			inputDTO: advertusecase.PatchDTO{
				UserID:   1,
				UserRole: "user",
				ID:       1,
				Version:  1,
			},
			mockBehavior:  func(mr *Mockrepo, dto advertusecase.PatchDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "forbidden",
			inputDTO: advertusecase.PatchDTO{
				UserID:   2,
				UserRole: "user",
				ID:       1,
				Version:  1,
				Title:    &newTitle,
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
					Status: advertentity.StatusActive,
				}, nil)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
		{
			name: "repo GetByID error",
			inputDTO: advertusecase.PatchDTO{
				UserID:   1,
				UserRole: "user",
				ID:       1,
				Version:  1,
				Title:    &newTitle,
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(advertentity.Advert{}, coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
		{
			name: "whitespace title",
			inputDTO: advertusecase.PatchDTO{
				UserID:   1,
				UserRole: "user",
				ID:       1,
				Version:  1,
				Title:    &whitespaceTitle,
			},
			mockBehavior:  func(mr *Mockrepo, dto advertusecase.PatchDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "user cant patch blocked advert",
			inputDTO: advertusecase.PatchDTO{
				UserID:   1,
				UserRole: "user",
				ID:       1,
				Version:  1,
				Title:    &newTitle,
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
					Status: advertentity.StatusBlocked,
				}, nil)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
		{
			name: "admin patches blocked advert without status reset",
			inputDTO: advertusecase.PatchDTO{
				UserID:   99,
				UserRole: "admin",
				ID:       1,
				Version:  1,
				Title:    &newTitle,
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
					Status: advertentity.StatusBlocked,
				}, nil)
				mr.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResult: func(t *testing.T, advert advertentity.Advert) {
				assert.Equal(t, "New Title", advert.Title)
				assert.Equal(t, advertentity.StatusBlocked, advert.Status)
			},
		},
		{
			name: "repo Patch error",
			inputDTO: advertusecase.PatchDTO{
				UserID:   1,
				UserRole: "user",
				ID:       1,
				Version:  1,
				Title:    &newTitle,
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.PatchDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.ID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
					Status: advertentity.StatusActive,
				}, nil)
				mr.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(coreerrors.ErrConflict)
			},
			expectedErrIs: coreerrors.ErrConflict,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			testCase.mockBehavior(mr, testCase.inputDTO)

			uc := advertusecase.New(NewMockfavRepo(c), mr, NewMockstorage(c))

			advert, err := uc.Patch(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			testCase.checkResult(t, advert)
		})
	}
}
