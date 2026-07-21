package userhttp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	userhttp "github.com/egotk/golang-advert-app/internal/features/user/controller/rest"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_Login(t *testing.T) {
	type loginMockBehavior func(muc *MockuseCase, dto userusecase.LoginDTO)

	testTable := []struct {
		name                 string
		inputBody            string
		inputDTO             userusecase.LoginDTO
		mockBehavior         loginMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"test@example.com","password":"password123"}`,
			inputDTO: userusecase.LoginDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockBehavior: func(muc *MockuseCase, dto userusecase.LoginDTO) {
				muc.EXPECT().Login(gomock.Any(), dto).Return(
					userusecase.LoginResultDTO{
						UserID: 1,
						Tokens: userusecase.TokensDTO{
							Access:  "access-token",
							Refresh: "refresh-token",
						},
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"user_id":1,"access_token":"access-token","refresh_token":"refresh-token"}`,
		},
		{
			name:               "invalid body",
			inputBody:          `invalid json`,
			mockBehavior:       func(muc *MockuseCase, dto userusecase.LoginDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "usecase error",
			inputBody: `{"email":"test@example.com","password":"wrong-password"}`,
			inputDTO: userusecase.LoginDTO{
				Email:    "test@example.com",
				Password: "wrong-password",
			},
			mockBehavior: func(muc *MockuseCase, dto userusecase.LoginDTO) {
				muc.EXPECT().Login(gomock.Any(), dto).Return(userusecase.LoginResultDTO{}, coreerrors.ErrUnauthorized)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputDTO)

			controller := userhttp.New(muc)

			r := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(testCase.inputBody))
			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.Login(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}