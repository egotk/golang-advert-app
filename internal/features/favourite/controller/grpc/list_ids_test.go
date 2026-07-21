package favgrpc_test

import (
	"context"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	favgrpc "github.com/egotk/golang-advert-app/internal/features/favourite/controller/grpc"
	favpb "github.com/egotk/golang-advert-app/internal/gen/favourite"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestController_ListIDs(t *testing.T) {
	type listIDsMockBehavior func(muc *MockuseCase, userID int64)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name             string
		withClaims       bool
		mockBehavior     listIDsMockBehavior
		expectedResponse *favpb.ListIDsResponse
		expectedErr      error
	}{
		{
			name:       "OK",
			withClaims: true,
			mockBehavior: func(muc *MockuseCase, userID int64) {
				muc.EXPECT().ListIDs(gomock.Any(), userID).Return([]int64{1, 2, 3}, nil)
			},
			expectedResponse: &favpb.ListIDsResponse{Ids: []int64{1, 2, 3}},
		},
		{
			name:         "missing claims",
			withClaims:   false,
			mockBehavior: func(muc *MockuseCase, userID int64) {},
			expectedErr:  coreerrors.ErrUnauthorized,
		},
		{
			name:       "usecase error",
			withClaims: true,
			mockBehavior: func(muc *MockuseCase, userID int64) {
				muc.EXPECT().ListIDs(gomock.Any(), userID).Return(nil, coreerrors.ErrNotFound)
			},
			expectedErr: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, 1)

			controller := favgrpc.New(muc)

			ctx := context.Background()
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			response, err := controller.ListIDs(ctx, &emptypb.Empty{})

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)

				return
			}

			assert.NoError(t, err)
			assert.True(t, proto.Equal(testCase.expectedResponse, response))
		})
	}
}
