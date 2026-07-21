package categoryhttp_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	categoryhttp "github.com/egotk/golang-advert-app/internal/features/category/controller/rest"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_Create(t *testing.T) {
	type createMockBehavior func(muc *MockuseCase, dto categoryusecase.CreateDTO)

	parentID := int64(1)

	testTable := []struct {
		name                 string
		inputBody            string
		inputDTO             categoryusecase.CreateDTO
		mockBehavior         createMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"parent_id":1,"name":"Electronics"}`,
			inputDTO: categoryusecase.CreateDTO{
				ParentID: &parentID,
				Name:     "Electronics",
			},
			mockBehavior: func(muc *MockuseCase, dto categoryusecase.CreateDTO) {
				muc.EXPECT().Create(gomock.Any(), dto).Return(
					categoryentity.Category{
						ID:       1,
						ParentID: dto.ParentID,
						Name:     dto.Name,
					}, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":1,"parent_id":1,"name":"Electronics"}`,
		},
		{
			name:               "invalid body",
			inputBody:          `invalid json`,
			mockBehavior:       func(muc *MockuseCase, dto categoryusecase.CreateDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "usecase error",
			inputBody: `{"name":"Electronics"}`,
			inputDTO: categoryusecase.CreateDTO{
				Name: "Electronics",
			},
			mockBehavior: func(muc *MockuseCase, dto categoryusecase.CreateDTO) {
				muc.EXPECT().Create(gomock.Any(), dto).Return(categoryentity.Category{}, coreerrors.ErrConflict)
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

			controller := categoryhttp.New(muc)

			r := httptest.NewRequest(http.MethodPost, "/adverts/categories", bytes.NewBufferString(testCase.inputBody))
			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.Create(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)

			if testCase.expectedResponseBody != "" {
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSuffix(rw.Body.String(), "\n"))
			}
		})
	}
}
