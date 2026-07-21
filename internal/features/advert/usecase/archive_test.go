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

func TestUseCase_Archive(t *testing.T) {
	type archiveMockBehavior func(mr *Mockrepo, dto advertusecase.ArchiveDTO)

	testTable := []struct {
		name          string
		inputDTO      advertusecase.ArchiveDTO
		mockBehavior  archiveMockBehavior
		expectedErrIs error
		expectedID    int64
	}{
		{
			name: "OK",
			inputDTO: advertusecase.ArchiveDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 1,
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.ArchiveDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
				}, nil)
				mr.EXPECT().ChangeStatus(gomock.Any(), dto.AdvertID, advertentity.StatusActive, advertentity.StatusArchived).Return(
					advertentity.Advert{ID: 1, Status: advertentity.StatusArchived}, nil)
			},
			expectedID: 1,
		},
		{
			name: "invalid id",
			inputDTO: advertusecase.ArchiveDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 0,
			},
			mockBehavior:  func(mr *Mockrepo, dto advertusecase.ArchiveDTO) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "forbidden",
			inputDTO: advertusecase.ArchiveDTO{
				UserID:   2,
				UserRole: "user",
				AdvertID: 1,
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.ArchiveDTO) {
				mr.EXPECT().GetByID(gomock.Any(), dto.AdvertID).Return(advertentity.Advert{
					ID:     1,
					UserID: 1,
				}, nil)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
		{
			name: "repo error",
			inputDTO: advertusecase.ArchiveDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 1,
			},
			mockBehavior: func(mr *Mockrepo, dto advertusecase.ArchiveDTO) {
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
			testCase.mockBehavior(mr, testCase.inputDTO)

			uc := advertusecase.New(NewMockfavRepo(c), mr, NewMockstorage(c))

			advert, err := uc.Archive(context.Background(), testCase.inputDTO)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedID, advert.ID)
		})
	}
}
