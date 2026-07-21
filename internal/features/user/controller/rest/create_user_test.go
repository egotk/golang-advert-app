package userhttp_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	userhttp "github.com/egotk/golang-advert-app/internal/features/user/controller/rest"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	userusecase "github.com/egotk/golang-advert-app/internal/features/user/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_CreateUser(t *testing.T) {
	type createUserMockBehavior func(muc *MockuseCase, dto userusecase.CreateDTO)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	testTable := []struct {
		name                 string
		inputBody            string
		inputDTO             userusecase.CreateDTO
		mockBehavior         createUserMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"test@example.com","full_name":"Test User","phone_number":"1234567890","password":"password123"}`,
			inputDTO: userusecase.CreateDTO{
				Email:       "test@example.com",
				FullName:    "Test User",
				PhoneNumber: "1234567890",
				Password:    "password123",
			},
			mockBehavior: func(muc *MockuseCase, dto userusecase.CreateDTO) {
				muc.EXPECT().CreateUser(gomock.Any(), dto).Return(
					userentity.User{
						ID:          1,
						Version:     1,
						Email:       dto.Email,
						FullName:    dto.FullName,
						PhoneNumber: dto.PhoneNumber,
						Role:        "user",
						CreatedAt:   fixTime,
						UpdatedAt:   fixTime,
					}, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":1,"version":1,"email":"test@example.com","full_name":"Test User","phone_number":"1234567890","role":"user","locked_until":null,"created_at":"2024-01-01T12:00:00Z","updated_at":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:               "invalid body",
			inputBody:          `invalid json`,
			mockBehavior:       func(muc *MockuseCase, dto userusecase.CreateDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "usecase error",
			inputBody: `{"email":"test@example.com","full_name":"Test User","phone_number":"1234567890","password":"password123"}`,
			inputDTO: userusecase.CreateDTO{
				Email:       "test@example.com",
				FullName:    "Test User",
				PhoneNumber: "1234567890",
				Password:    "password123",
			},
			mockBehavior: func(muc *MockuseCase, dto userusecase.CreateDTO) {
				muc.EXPECT().CreateUser(gomock.Any(), dto).Return(userentity.User{}, coreerrors.ErrConflict)
			},
			expectedStatusCode: http.StatusConflict,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputDTO)

			controller := userhttp.New(muc)

			r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(testCase.inputBody))
			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.CreateUser(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}
