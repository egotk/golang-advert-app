package advertrest_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertrest "github.com/egotk/golang-advert-app/internal/features/advert/controller/rest"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// image/png http.DetectContentType
var pngHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

func buildCreateImagesMultipartBody(t *testing.T, advertID string, imageBytes []byte) (string, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if advertID != "" {
		if err := writer.WriteField("advert_id", advertID); err != nil {
			t.Fatalf("write field: %v", err)
		}
	}

	if imageBytes != nil {
		part, err := writer.CreateFormFile("images", "photo.png")
		if err != nil {
			t.Fatalf("create form file: %v", err)
		}

		if _, err := part.Write(imageBytes); err != nil {
			t.Fatalf("write image bytes: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	return body.String(), writer.FormDataContentType()
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
		name               string
		advertID           string
		imageBytes         []byte
		withClaims         bool
		mockBehavior       createImagesMockBehavior
		expectedStatusCode int
		checkBody          func(t *testing.T, body []byte)
	}{
		{
			name:       "OK",
			advertID:   "1",
			imageBytes: pngHeader,
			withClaims: true,
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().CreateImages(gomock.Any(), gomock.Any()).Return(
					[]advertentity.AdvertImage{{ID: 1, Name: "photo.png"}}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			checkBody: func(t *testing.T, body []byte) {
				var resp struct {
					Count  int64 `json:"count"`
					Images []struct {
						ID   int64  `json:"id"`
						Name string `json:"name"`
					} `json:"images"`
				}
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("unmarshal response: %v", err)
				}

				assert.Equal(t, int64(1), resp.Count)
				assert.Len(t, resp.Images, 1)
				assert.Equal(t, "photo.png", resp.Images[0].Name)
			},
		},
		{
			name:               "invalid body",
			advertID:           "",
			imageBytes:         nil,
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "unsupported content type",
			advertID:           "1",
			imageBytes:         []byte("this is definitely not an image"),
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing claims",
			advertID:           "1",
			imageBytes:         pngHeader,
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:       "usecase error",
			advertID:   "1",
			imageBytes: pngHeader,
			withClaims: true,
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().CreateImages(gomock.Any(), gomock.Any()).Return(nil, coreerrors.ErrForbidden)
			},
			expectedStatusCode: http.StatusForbidden,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc)

			controller := advertrest.New(muc)

			body, contentType := buildCreateImagesMultipartBody(t, testCase.advertID, testCase.imageBytes)

			r := httptest.NewRequest(http.MethodPost, "/adverts/images", strings.NewReader(body))
			r.Header.Set("Content-Type", contentType)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.CreateImages(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.checkBody != nil {
				testCase.checkBody(t, rw.Body.Bytes())
			}
		})
	}
}
