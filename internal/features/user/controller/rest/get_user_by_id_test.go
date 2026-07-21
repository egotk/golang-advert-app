package userhttp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	userhttp "github.com/egotk/golang-advert-app/internal/features/user/controller/rest"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_GetUserByID(t *testing.T) {
	type getUserByIDMockBehavior func(muc *MockuseCase, id int64)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	testTable := []struct {
		name                 string
		idPathValue          string
		inputID              int64
		mockBehavior         getUserByIDMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			idPathValue: "1",
			inputID:     1,
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().GetUserByID(gomock.Any(), id).Return(
					userentity.User{
						ID:          1,
						Version:     1,
						Email:       "test@example.com",
						FullName:    "Test User",
						PhoneNumber: "1234567890",
						Role:        "user",
						CreatedAt:   fixTime,
						UpdatedAt:   fixTime,
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"version":1,"email":"test@example.com","full_name":"Test User","phone_number":"1234567890","role":"user","locked_until":null,"created_at":"2024-01-01T12:00:00Z","updated_at":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:               "invalid id param",
			idPathValue:        "abc",
			mockBehavior:       func(muc *MockuseCase, id int64) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "usecase error",
			idPathValue: "1",
			inputID:     1,
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().GetUserByID(gomock.Any(), id).Return(userentity.User{}, coreerrors.ErrNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, testCase.inputID)

			controller := userhttp.New(muc)

			r := httptest.NewRequest(http.MethodGet, "/users/"+testCase.idPathValue, nil)
			r.SetPathValue("id", testCase.idPathValue)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.GetUserByID(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}
