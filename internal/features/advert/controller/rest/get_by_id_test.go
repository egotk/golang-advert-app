package advertrest_test

import (
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

func TestController_GetByID(t *testing.T) {
	type getByIDMockBehavior func(muc *MockuseCase, dto advertusecase.GetByIDDTO)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name                 string
		idPathValue          string
		withClaims           bool
		mockBehavior         getByIDMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			idPathValue: "1",
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.GetByIDDTO) {
				muc.EXPECT().GetByID(gomock.Any(), dto).Return(advertentity.Advert{
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
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"version":1,"user_id":1,"title":"Title","description":"Description","price":100,"category_id":1,"status":"active","advert_images":{"count":0,"images":[]},"created_at":"2024-01-01T12:00:00Z","updated_at":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:               "invalid id param",
			idPathValue:        "abc",
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.GetByIDDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing claims",
			idPathValue:        "1",
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.GetByIDDTO) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:        "usecase error",
			idPathValue: "1",
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.GetByIDDTO) {
				muc.EXPECT().GetByID(gomock.Any(), dto).Return(advertentity.Advert{}, coreerrors.ErrNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.GetByIDDTO{
				UserID:   1,
				UserRole: "user",
				AdvertID: 1,
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertrest.New(muc)

			r := httptest.NewRequest(http.MethodGet, "/adverts/"+testCase.idPathValue, nil)
			r.SetPathValue("id", testCase.idPathValue)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.GetByID(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}
