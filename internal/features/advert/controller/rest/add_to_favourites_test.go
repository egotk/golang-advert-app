package advertrest_test

import (
	"net/http"
	"net/http/httptest"
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

func TestController_AddToFavourites(t *testing.T) {
	type addToFavouritesMockBehavior func(muc *MockuseCase, dto advertusecase.AddToFavouritesDTO)

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name               string
		idPathValue        string
		withClaims         bool
		mockBehavior       addToFavouritesMockBehavior
		expectedStatusCode int
	}{
		{
			name:        "OK",
			idPathValue: "1",
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.AddToFavouritesDTO) {
				muc.EXPECT().AddToFavourites(gomock.Any(), dto).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "invalid id param",
			idPathValue:        "abc",
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.AddToFavouritesDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing claims",
			idPathValue:        "1",
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.AddToFavouritesDTO) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:        "usecase error",
			idPathValue: "1",
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.AddToFavouritesDTO) {
				muc.EXPECT().AddToFavourites(gomock.Any(), dto).Return(coreerrors.ErrForbidden)
			},
			expectedStatusCode: http.StatusForbidden,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.AddToFavouritesDTO{
				AdvertID: 1,
				UserID:   1,
				UserRole: "user",
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertrest.New(muc)

			r := httptest.NewRequest(http.MethodPost, "/adverts/favourites/"+testCase.idPathValue, nil)
			r.SetPathValue("id", testCase.idPathValue)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.AddToFavourites(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)
		})
	}
}