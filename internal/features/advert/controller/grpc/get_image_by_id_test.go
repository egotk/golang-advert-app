package advertgrpc_test

import (
	"context"
	"io"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertgrpc "github.com/egotk/golang-advert-app/internal/features/advert/controller/grpc"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
)

type fakeGetImageStream struct {
	grpc.ServerStream
	ctx  context.Context
	sent []*advertpb.AdvertImageResponse
}

func (s *fakeGetImageStream) Context() context.Context {
	return s.ctx
}

func (s *fakeGetImageStream) Send(resp *advertpb.AdvertImageResponse) error {
	s.sent = append(s.sent, resp)

	return nil
}

func TestController_GetImageByID(t *testing.T) {
	type getImageByIDMockBehavior func(muc *MockuseCase, dto advertusecase.GetImageDTO)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	const imageContent = "fake image bytes"

	testTable := []struct {
		name          string
		inputRequest  *advertpb.GetImageByIDRequest
		withClaims    bool
		mockBehavior  getImageByIDMockBehavior
		expectedErrIs error
		checkStream   func(t *testing.T, sent []*advertpb.AdvertImageResponse)
	}{
		{
			name:         "OK",
			inputRequest: &advertpb.GetImageByIDRequest{Id: 1},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.GetImageDTO) {
				muc.EXPECT().GetImageByID(gomock.Any(), dto).Return(
					io.NopCloser(strings.NewReader(imageContent)),
					advertentity.AdvertImage{ID: 1, Name: "photo.jpg"},
					nil,
				)
			},
			checkStream: func(t *testing.T, sent []*advertpb.AdvertImageResponse) {
				assert.GreaterOrEqual(t, len(sent), 2)
				assert.Equal(t, "photo.jpg", sent[0].GetImage().GetName())

				var data []byte
				for _, resp := range sent[1:] {
					data = append(data, resp.GetData()...)
				}
				assert.Equal(t, imageContent, string(data))
			},
		},
		{
			name:          "missing claims",
			inputRequest:  &advertpb.GetImageByIDRequest{Id: 1},
			withClaims:    false,
			mockBehavior:  func(muc *MockuseCase, dto advertusecase.GetImageDTO) {},
			expectedErrIs: coreerrors.ErrUnauthorized,
		},
		{
			name:         "usecase error",
			inputRequest: &advertpb.GetImageByIDRequest{Id: 1},
			withClaims:   true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.GetImageDTO) {
				muc.EXPECT().GetImageByID(gomock.Any(), dto).Return(nil, advertentity.AdvertImage{}, coreerrors.ErrNotFound)
			},
			expectedErrIs: coreerrors.ErrNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.GetImageDTO{
				UserID:   1,
				UserRole: "user",
				ImageID:  testCase.inputRequest.Id,
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertgrpc.New(muc)

			ctx := context.Background()
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			stream := &fakeGetImageStream{ctx: ctx}

			err := controller.GetImageByID(testCase.inputRequest, stream)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			testCase.checkStream(t, stream.sent)
		})
	}
}
