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

func TestController_Patch(t *testing.T) {
	type patchMockBehavior func(muc *MockuseCase, dto advertusecase.PatchDTO)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	newTitle := "New Title"

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name             string
		inputRequest     *advertpb.PatchRequest
		withClaims       bool
		mockBehavior     patchMockBehavior
		expectedResponse *advertpb.AdvertResponse
		expectedErrIs    error
	}{
		{
			name: "OK",
			inputRequest: &advertpb.PatchRequest{
				Id:      1,
				Version: 1,
				Title:   &newTitle,
			},
			withClaims: true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.PatchDTO) {
				muc.EXPECT().Patch(gomock.Any(), dto).Return(advertentity.Advert{
					ID:        1,
					Version:   1,
					UserID:    1,
					Title:     "New Title",
					Status:    advertentity.StatusInitial,
					CreatedAt: fixTime,
					UpdatedAt: fixTime,
				}, nil)
			},
			expectedResponse: &advertpb.AdvertResponse{
				Id:        1,
				Version:   1,
				UserId:    1,
				Title:     "New Title",
				Status:    "initial",
				CreatedAt: timestamppb.New(fixTime),
				UpdatedAt: timestamppb.New(fixTime),
			},
		},
		{
			name: "missing claims",
			inputRequest: &advertpb.PatchRequest{
				Id:      1,
				Version: 1,
				Title:   &newTitle,
			},
			withClaims:    false,
			mockBehavior:  func(muc *MockuseCase, dto advertusecase.PatchDTO) {},
			expectedErrIs: coreerrors.ErrUnauthorized,
		},
		{
			name: "usecase error",
			inputRequest: &advertpb.PatchRequest{
				Id:      1,
				Version: 1,
				Title:   &newTitle,
			},
			withClaims: true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.PatchDTO) {
				muc.EXPECT().Patch(gomock.Any(), dto).Return(advertentity.Advert{}, coreerrors.ErrForbidden)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.PatchDTO{
				UserID:   1,
				UserRole: "user",
				ID:       testCase.inputRequest.Id,
				Version:  testCase.inputRequest.Version,
				Title:    testCase.inputRequest.Title,
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertgrpc.New(muc)

			ctx := context.Background()
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			response, err := controller.Patch(ctx, testCase.inputRequest)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}