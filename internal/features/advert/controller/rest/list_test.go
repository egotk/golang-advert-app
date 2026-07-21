package advertrest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestController_List(t *testing.T) {
	type listMockBehavior func(muc *MockuseCase, dto advertusecase.ListDTO)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name                 string
		queryString          string
		withClaims           bool
		mockBehavior         listMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			queryString: "",
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.ListDTO) {
				muc.EXPECT().List(gomock.Any(), dto).Return(int64(1), []advertentity.Advert{
					{
						ID:          1,
						Version:     1,
						UserID:      1,
						Title:       "Title",
						Description: "Description",
						Price:       100,
						CategoryID:  1,
						Status:      advertentity.StatusActive,
						CreatedAt:   fixTime,
						UpdatedAt:   fixTime,
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"count":1,"adverts":[{"id":1,"version":1,"user_id":1,"title":"Title","description":"Description","price":100,"category_id":1,"status":"active","advert_images":{"count":0,"images":[]},"created_at":"2024-01-01T12:00:00Z","updated_at":"2024-01-01T12:00:00Z"}]}`,
		},
		{
			name:               "invalid query param",
			queryString:        "?limit=abc",
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.ListDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing claims",
			queryString:        "",
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.ListDTO) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:        "usecase error",
			queryString: "",
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.ListDTO) {
				muc.EXPECT().List(gomock.Any(), dto).Return(int64(0), nil, errors.New("boom"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.ListDTO{
				UserID:   1,
				UserRole: "user",
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertrest.New(muc)

			r := httptest.NewRequest(http.MethodGet, "/adverts"+testCase.queryString, nil)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.List(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}
