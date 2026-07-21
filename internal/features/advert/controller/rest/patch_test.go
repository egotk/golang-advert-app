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

func TestController_Patch(t *testing.T) {
	type patchMockBehavior func(muc *MockuseCase, dto advertusecase.PatchDTO)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	newTitle := "New Title"

	claims := corejwt.Claims{
		Role: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	testTable := []struct {
		name                 string
		idPathValue          string
		inputBody            string
		withClaims           bool
		mockBehavior         patchMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			idPathValue: "1",
			inputBody:   `{"version":1,"title":"New Title"}`,
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.PatchDTO) {
				muc.EXPECT().Patch(gomock.Any(), dto).Return(advertentity.Advert{
					ID:      1,
					Version: 1,
					UserID:  1,
					Title:   "New Title",
					Status:  advertentity.StatusInitial,

					CreatedAt: fixTime,
					UpdatedAt: fixTime,
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"version":1,"user_id":1,"title":"New Title","description":"","price":0,"category_id":0,"status":"initial","advert_images":{"count":0,"images":[]},"created_at":"2024-01-01T12:00:00Z","updated_at":"2024-01-01T12:00:00Z"}`,
		},
		{
			name:               "invalid body",
			idPathValue:        "1",
			inputBody:          `invalid json`,
			withClaims:         true,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.PatchDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing claims",
			idPathValue:        "1",
			inputBody:          `{"version":1,"title":"New Title"}`,
			withClaims:         false,
			mockBehavior:       func(muc *MockuseCase, dto advertusecase.PatchDTO) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:        "usecase error",
			idPathValue: "1",
			inputBody:   `{"version":1,"title":"New Title"}`,
			withClaims:  true,
			mockBehavior: func(muc *MockuseCase, dto advertusecase.PatchDTO) {
				muc.EXPECT().Patch(gomock.Any(), dto).Return(advertentity.Advert{}, coreerrors.ErrForbidden)
			},
			expectedStatusCode: http.StatusForbidden,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			dto := advertusecase.PatchDTO{
				UserID:   1,
				UserRole: "user",
				ID:       1,
				Version:  1,
				Title:    &newTitle,
			}

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc, dto)

			controller := advertrest.New(muc)

			r := httptest.NewRequest(http.MethodPatch, "/adverts/"+testCase.idPathValue, strings.NewReader(testCase.inputBody))
			r.SetPathValue("id", testCase.idPathValue)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			if testCase.withClaims {
				ctx = corejwt.ClaimsToContext(ctx, claims)
			}
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.Patch(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}
