package advertgrpc_test

import (
	"context"
	"io"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	advertgrpc "github.com/egotk/golang-advert-app/internal/features/advert/controller/grpc"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertpb "github.com/egotk/golang-advert-app/internal/gen/advert"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
)

// http.DetectContentType image/png
var pngHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

type fakeCreateImagesStream struct {
	grpc.ServerStream
	ctx      context.Context
	requests []*advertpb.CreateImagesRequest
	recvIdx  int
	response *advertpb.AdvertImagesResponse
}

func (s *fakeCreateImagesStream) Context() context.Context {
	return s.ctx
}

func (s *fakeCreateImagesStream) Recv() (*advertpb.CreateImagesRequest, error) {
	if s.recvIdx >= len(s.requests) {
		return nil, io.EOF
	}

	req := s.requests[s.recvIdx]
	s.recvIdx++

	return req, nil
}

func (s *fakeCreateImagesStream) SendAndClose(resp *advertpb.AdvertImagesResponse) error {
	s.response = resp

	return nil
}

func advertIDRequest(advertID int64) *advertpb.CreateImagesRequest {
	return &advertpb.CreateImagesRequest{
		ImgData: &advertpb.CreateImagesRequest_AdvertId{AdvertId: advertID},
	}
}

func imageRequest(name string, data []byte) *advertpb.CreateImagesRequest {
	return &advertpb.CreateImagesRequest{
		ImgData: &advertpb.CreateImagesRequest_Image{
			Image: &advertpb.UploadImage{Name: name, Data: data},
		},
	}
}

func imageRequests(name string, data []byte, count int) []*advertpb.CreateImagesRequest {
	requests := make([]*advertpb.CreateImagesRequest, count)
	for i := range requests {
		requests[i] = imageRequest(name, data)
	}

	return requests
}

func TestController_CreateImages(t *testing.T) {
	type createImagesMockBehavior func(muc *MockuseCase)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name          string
		requests      []*advertpb.CreateImagesRequest
		withClaims    bool
		mockBehavior  createImagesMockBehavior
		expectedErrIs error
		checkResponse func(t *testing.T, resp *advertpb.AdvertImagesResponse)
	}{
		{
			name: "OK",
			requests: []*advertpb.CreateImagesRequest{
				advertIDRequest(1),
				imageRequest("photo.png", pngHeader),
			},
			withClaims: true,
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().CreateImages(gomock.Any(), gomock.Any()).Return(
					[]advertentity.AdvertImage{{ID: 1, Name: "photo.png"}}, nil)
			},
			checkResponse: func(t *testing.T, resp *advertpb.AdvertImagesResponse) {
				assert.Equal(t, int64(1), resp.GetCount())
				assert.Len(t, resp.GetImages(), 1)
				assert.Equal(t, "photo.png", resp.GetImages()[0].GetName())
			},
		},
		{
			name: "non-positive advert id",
			requests: []*advertpb.CreateImagesRequest{
				advertIDRequest(0),
			},
			withClaims:    true,
			mockBehavior:  func(muc *MockuseCase) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "missing claims",
			requests: []*advertpb.CreateImagesRequest{
				advertIDRequest(1),
				imageRequest("photo.png", pngHeader),
			},
			withClaims:    false,
			mockBehavior:  func(muc *MockuseCase) {},
			expectedErrIs: coreerrors.ErrUnauthorized,
		},
		{
			name: "unsupported content type",
			requests: []*advertpb.CreateImagesRequest{
				advertIDRequest(1),
				imageRequest("note.txt", []byte("this is definitely not an image")),
			},
			withClaims:    true,
			mockBehavior:  func(muc *MockuseCase) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "too many images",
			requests: append(
				[]*advertpb.CreateImagesRequest{advertIDRequest(1)},
				imageRequests("photo.png", pngHeader, 11)...,
			),
			withClaims:    true,
			mockBehavior:  func(muc *MockuseCase) {},
			expectedErrIs: coreerrors.ErrInvalidArgument,
		},
		{
			name: "usecase error",
			requests: []*advertpb.CreateImagesRequest{
				advertIDRequest(1),
				imageRequest("photo.png", pngHeader),
			},
			withClaims: true,
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().CreateImages(gomock.Any(), gomock.Any()).Return(nil, coreerrors.ErrForbidden)
			},
			expectedErrIs: coreerrors.ErrForbidden,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc)

			controller := advertgrpc.New(muc)

			ctx := context.Background()
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			stream := &fakeCreateImagesStream{ctx: ctx, requests: testCase.requests}

			err := controller.CreateImages(stream)

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			testCase.checkResponse(t, stream.response)
		})
	}
}
