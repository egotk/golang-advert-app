package favrest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	favrest "github.com/egotk/golang-advert-app/internal/features/favourite/controller/rest"
	favusecase "github.com/egotk/golang-advert-app/internal/features/favourite/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_Remove(t *testing.T) {
	type removeMockBehavior func(muc *MockuseCase, dto favusecase.RemoveDTO)

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
		inputDTO           favusecase.RemoveDTO
		mockBehavior       removeMockBehavior
		expectedStatusCode int
	}{
		{
			name:        "OK",
			idPathValue: "1",
			withClaims:  true,
			inputDTO:    favusecase.RemoveDTO{AdvertID: 1, UserID: 1},
			mockBehavior: func(muc *MockuseCase, dto favusecase.RemoveDTO) {
				muc.EXPECT().Remove(gomock.Any(), dto).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "invalid id path param",
			idPathValue:        "abc",
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase, dto favusecase.RemoveDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing claims",
			idPathValue:        "1",
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase, dto favusecase.RemoveDTO) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:        "usecase error",
			idPathValue: "1",
			withClaims:  true,
			inputDTO:    favusecase.RemoveDTO{AdvertID: 1, UserID: 1},
			mockBehavior: func(muc *MockuseCase, dto favusecase.RemoveDTO) {
				muc.EXPECT().Remove(gomock.Any(), dto).Return(coreerrors.ErrNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputDTO)

			controller := favrest.New(muc)

			r := httptest.NewRequest(http.MethodDelete, "/adverts/favourites/"+testCase.idPathValue, nil)
			r.SetPathValue("id", testCase.idPathValue)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})

			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}

			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.Remove(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)
		})
	}
}
