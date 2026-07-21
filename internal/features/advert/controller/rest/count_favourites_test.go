package advertrest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertrest "github.com/egotk/golang-advert-app/internal/features/advert/controller/rest"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_CountFavourites(t *testing.T) {
	type countFavouritesMockBehavior func(muc *MockuseCase, dto advertusecase.CountDTO)

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
		mockBehavior         countFavouritesMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			queryString: "",
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.CountDTO) {
				muc.EXPECT().CountFavourites(gomock.Any(), dto).Return(int64(3), nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"count":3}`,
		},
		{
			name:               "invalid query param",
			queryString:        "?min_price=abc",
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.CountDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing claims",
			queryString:        "",
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.CountDTO) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:        "usecase error",
			queryString: "",
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.CountDTO) {
				muc.EXPECT().CountFavourites(gomock.Any(), dto).Return(int64(0), coreerrors.ErrInvalidArgument)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.CountDTO{
				UserID:   1,
				UserRole: "user",
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertrest.New(muc)

			r := httptest.NewRequest(http.MethodGet, "/adverts/favourites/count"+testCase.queryString, nil)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.CountFavourites(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}