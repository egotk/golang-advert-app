package advertrest_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestController_GetImageByID(t *testing.T) {
	type getImageByIDMockBehavior func(muc *MockuseCase, dto advertusecase.GetImageDTO)

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
		mockBehavior       getImageByIDMockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:        "OK",
			idPathValue: "1",
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.GetImageDTO) {
				muc.EXPECT().GetImageByID(gomock.Any(), dto).Return(
					io.NopCloser(strings.NewReader("fake image bytes")),
					advertentity.AdvertImage{ID: 1, Name: "photo.jpg", Path: "path/to/photo.jpg"},
					nil,
				)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       "fake image bytes",
		},
		{
			name:               "invalid id param",
			idPathValue:        "abc",
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.GetImageDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing claims",
			idPathValue:        "1",
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.GetImageDTO) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:        "usecase error",
			idPathValue: "1",
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.GetImageDTO) {
				muc.EXPECT().GetImageByID(gomock.Any(), dto).Return(nil, advertentity.AdvertImage{}, coreerrors.ErrNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.GetImageDTO{
				UserID:   1,
				UserRole: "user",
				ImageID:  1,
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertrest.New(muc)

			r := httptest.NewRequest(http.MethodGet, "/adverts/images/"+testCase.idPathValue, nil)
			r.SetPathValue("id", testCase.idPathValue)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.GetImageByID(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedBody != "" {
				assert.Equal(t, testCase.expectedBody, rw.Body.String())
			}
		})
	}
}