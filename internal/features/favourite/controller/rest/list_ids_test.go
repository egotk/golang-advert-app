package favrest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	favrest "github.com/egotk/golang-advert-app/internal/features/favourite/controller/rest"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
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
		name                 string
		withClaims           bool
		mockBehavior         listIDsMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "OK",
			withClaims: true,
			mockBehavior: func(muc *MockuseCase, userID int64) {
				muc.EXPECT().ListIDs(gomock.Any(), userID).Return([]int64{1, 2, 3}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"ids":[1,2,3]}`,
		},
		{
			name:               "missing claims",
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase, userID int64) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:       "usecase error",
			withClaims: true,
			mockBehavior: func(muc *MockuseCase, userID int64) {
				muc.EXPECT().ListIDs(gomock.Any(), userID).Return(nil, coreerrors.ErrNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, 1)

			controller := favrest.New(muc)

			r := httptest.NewRequest(http.MethodGet, "/adverts/favourites/ids", nil)
			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})

			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.ListIDs(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}
