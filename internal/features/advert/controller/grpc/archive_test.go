package advertgrpc_test

import (
	"context"
	"testing"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertgrpc "github.com/egotk/golang-advert-app/internal/features/advert/controller/grpc"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestController_Archive(t *testing.T) {
	type archiveMockBehavior func(muc *MockuseCase, dto advertusecase.ArchiveDTO)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name             string
		inputRequest     *advertpb.ArchiveRequest
		withClaims       bool
		mockBehavior     archiveMockBehavior
		expectedResponse *advertpb.AdvertResponse
		expectedErrIs    error
	}{
		{
			name:         "OK",
			inputRequest: &advertpb.ArchiveRequest{Id: 1},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.ArchiveDTO) {
				muc.EXPECT().Archive(gomock.Any(), dto).Return(advertentity.Advert{
					ID:        1,
					Version:   2,
					UserID:    1,
					Title:     "Title",
					Status:    advertentity.StatusArchived,
					CreatedAt: fixTime,
					UpdatedAt: fixTime,
				}, nil)
			},
			expectedResponse: &advertpb.AdvertResponse{
				Id:        1,
				Version:   2,
				UserId:    1,
				Title:     "Title",
				Status:    "archived",
				CreatedAt: timestamppb.New(fixTime),
				UpdatedAt: timestamppb.New(fixTime),
			},
		},
		{
			name:          "missing claims",
			inputRequest:  &advertpb.ArchiveRequest{Id: 1},
			withClaims:    false,
			mockBehavior:  func(muc *MockuseCase, dto advertusecase.ArchiveDTO) {},
			expectedErrIs: coreerrors.ErrUnauthorized,
		},
		{
			name:         "usecase error",
			inputRequest: &advertpb.ArchiveRequest{Id: 1},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.ArchiveDTO) {
				muc.EXPECT().Archive(gomock.Any(), dto).Return(advertentity.Advert{}, coreerrors.ErrForbidden)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.ArchiveDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: testCase.inputRequest.Id,
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertgrpc.New(muc)

			ctx := context.Background()
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			response, err := controller.Archive(ctx, testCase.inputRequest)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}