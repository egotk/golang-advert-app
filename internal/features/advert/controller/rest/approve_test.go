package advertrest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertrest "github.com/egotk/golang-advert-app/internal/features/advert/controller/rest"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_Approve(t *testing.T) {
	type approveMockBehavior func(muc *MockuseCase, id int64)

	fixTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	testTable := []struct {
		name                 string
		idPathValue          string
		mockBehavior         approveMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			idPathValue: "1",
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().Approve(gomock.Any(), id).Return(advertentity.Advert{
					ID:        1,
					Version:   2,
					UserID:    1,
					Title:     "Title",
					Status:    advertentity.StatusActive,
					CreatedAt: fixTime,
					UpdatedAt: fixTime,
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"version":2,"user_id":1,"title":"Title","description":"","price":0,"category_id":0,"status":"active","advert_images":{"count":0,"images":[]},"created_at":"2024-01-01T12:00:00Z","updated_at":"2024-01-01T12:00:00Z"}`,
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
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().Approve(gomock.Any(), id).Return(advertentity.Advert{}, coreerrors.ErrNotFound)
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

			controller := advertrest.New(muc)

			r := httptest.NewRequest(http.MethodPost, "/adverts/"+testCase.idPathValue+"/approve", nil)
			r.SetPathValue("id", testCase.idPathValue)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.Approve(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}