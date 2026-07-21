package advertrest_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertrest "github.com/egotk/golang-advert-app/internal/features/advert/controller/rest"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func buildMultipartBody(t *testing.T, fields map[string]string) (string, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			t.Fatalf("write field: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	return body.String(), writer.FormDataContentType()
}

func TestController_Create(t *testing.T) {
	type createMockBehavior func(muc *MockuseCase, dto advertusecase.CreateDTO)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name                 string
		formFields           map[string]string
		withClaims           bool
		mockBehavior         createMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			formFields: map[string]string{
				"title":       "Title",
				"description": "Description",
				"price":       "100",
				"category_id": "1",
			},
			withClaims: true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.CreateDTO) {
				muc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(advertentity.Advert{
					ID:          1,
					Version:     1,
					UserID:      1,
					Title:       "Title",
					Description: "Description",
					Price:       100,
					CategoryID:  1,
					Status:      advertentity.StatusInitial,
					CreatedAt:   fixTime,
					UpdatedAt:   fixTime,
				}, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":1,"version":1,"user_id":1,"title":"Title","description":"Description","price":100,"category_id":1,"status":"initial","advert_images":{"count":0,"images":[]},"created_at":"2024-01-01T12:00:00Z","updated_at":"2024-01-01T12:00:00Z"}`,
		},
		{
			name: "invalid body",
			formFields: map[string]string{
				"title": "Title",
			},
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.CreateDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "missing claims",
			formFields: map[string]string{
				"title":       "Title",
				"description": "Description",
				"price":       "100",
				"category_id": "1",
			},
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.CreateDTO) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "usecase error",
			formFields: map[string]string{
				"title":       "Title",
				"description": "Description",
				"price":       "100",
				"category_id": "1",
			},
			withClaims: true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.CreateDTO) {
				muc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(advertentity.Advert{}, coreerrors.ErrInvalidArgument)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, advertusecase.CreateDTO{})

			controller := advertrest.New(muc)

			body, contentType := buildMultipartBody(t, testCase.formFields)

			r := httptest.NewRequest(http.MethodPost, "/adverts", strings.NewReader(body))
			r.Header.Set("Content-Type", contentType)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.Create(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}
