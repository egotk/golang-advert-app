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

func TestController_GetMyAdverts(t *testing.T) {
	type getMyAdvertsMockBehavior func(muc *MockuseCase, dto advertusecase.ListDTO)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	userID := int64(1)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name             string
		inputRequest     *advertpb.GetMyAdvertsRequest
		withClaims       bool
		mockBehavior     getMyAdvertsMockBehavior
		expectedResponse *advertpb.AdvertsResponse
		expectedErrIs    error
	}{
		{
			name:         "OK",
			inputRequest: &advertpb.GetMyAdvertsRequest{},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.ListDTO) {
				muc.EXPECT().List(gomock.Any(), dto).Return(int64(1), []advertentity.Advert{
					{
						ID:        1,
						UserID:    1,
						Title:     "Title",
						Status:    advertentity.StatusActive,
						CreatedAt: fixTime,
						UpdatedAt: fixTime,
					},
				}, nil)
			},
			expectedResponse: &advertpb.AdvertsResponse{
				Count: 1,
				Adverts: []*advertpb.AdvertResponse{
					{
						Id:        1,
						UserId:    1,
						Title:     "Title",
						Status:    "active",
						CreatedAt: timestamppb.New(fixTime),
						UpdatedAt: timestamppb.New(fixTime),
					},
				},
			},
		},
		{
			name:          "missing claims",
			inputRequest:  &advertpb.GetMyAdvertsRequest{},
			withClaims:    false,
			mockBehavior:  func(muc *MockuseCase, dto advertusecase.ListDTO) {},
			expectedErrIs: coreerrors.ErrUnauthorized,
		},
		{
			name:         "usecase error",
			inputRequest: &advertpb.GetMyAdvertsRequest{},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.ListDTO) {
				muc.EXPECT().List(gomock.Any(), dto).Return(int64(0), nil, coreerrors.ErrInvalidArgument)
			},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.ListDTO{
				UserID: userID,
				Filter: advertentity.Filter{UserID: &userID},
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertgrpc.New(muc)

			ctx := context.Background()
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			response, err := controller.GetMyAdverts(ctx, testCase.inputRequest)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}