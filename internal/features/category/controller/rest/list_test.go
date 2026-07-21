package categoryhttp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	categoryhttp "github.com/egotk/golang-advert-app/internal/features/category/controller/rest"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_List(t *testing.T) {
	type listMockBehavior func(muc *MockuseCase)

	parentID := int64(1)

	testTable := []struct {
		name                 string
		mockBehavior         listMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().List(gomock.Any()).Return(
					[]categoryentity.Category{
						{
							ID:       1,
							ParentID: nil,
							Name:     "Electronics",
						},
						{
							ID:       2,
							ParentID: &parentID,
							Name:     "Phones",
						},
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"categories":[{"id":1,"parent_id":null,"name":"Electronics"},{"id":2,"parent_id":1,"name":"Phones"}]}`,
		},
		{
			name: "usecase error",
			mockBehavior: func(muc *MockuseCase) {
				muc.EXPECT().List(gomock.Any()).Return(nil, errors.New("boom"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			muc := NewMockuseCase(c)
			testCase.mockBehavior(muc)

			controller := categoryhttp.New(muc)

			r := httptest.NewRequest(http.MethodGet, "/adverts/categories", nil)
			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.List(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}
