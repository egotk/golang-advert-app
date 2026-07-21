package categoryhttp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	"github.com/egotk/golang-advert-app/internal/core/nullable"
	categoryhttp "github.com/egotk/golang-advert-app/internal/features/category/controller/rest"
	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_Patch(t *testing.T) {
	type patchMockBehavior func(muc *MockuseCase, dto categoryusecase.PatchDTO)

	newParentID := int64(2)
	newName := "Updated"

	testTable := []struct {
		name                 string
		idPathValue          string
		inputBody            string
		inputDTO             categoryusecase.PatchDTO
		mockBehavior         patchMockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			idPathValue: "1",
			inputBody:   `{"parent_id":2,"name":"Updated"}`,
			inputDTO: categoryusecase.PatchDTO{
				ID:       1,
				ParentID: nullable.Nullable[int64]{Set: true, Value: &newParentID},
				Name:     &newName,
			},
			mockBehavior: func(muc *MockuseCase, dto categoryusecase.PatchDTO) {
				muc.EXPECT().Patch(gomock.Any(), dto).Return(
					categoryentity.Category{
						ID:       1,
						ParentID: dto.ParentID.Value,
						Name:     *dto.Name,
					}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"parent_id":2,"name":"Updated"}`,
		},
		{
			name:               "invalid body",
			idPathValue:        "1",
			inputBody:          `invalid json`,
			mockBehavior:       func(muc *MockuseCase, dto categoryusecase.PatchDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid id param",
			idPathValue:        "abc",
			inputBody:          `{"name":"Updated"}`,
			mockBehavior:       func(muc *MockuseCase, dto categoryusecase.PatchDTO) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "usecase error",
			idPathValue: "1",
			inputBody:   `{"name":"Updated"}`,
			inputDTO: categoryusecase.PatchDTO{
				ID:   1,
				Name: &newName,
			},
			mockBehavior: func(muc *MockuseCase, dto categoryusecase.PatchDTO) {
				muc.EXPECT().Patch(gomock.Any(), dto).Return(categoryentity.Category{}, coreerrors.ErrNotFound)
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

			controller := categoryhttp.New(muc)

			r := httptest.NewRequest(http.MethodPatch, "/adverts/categories/"+testCase.idPathValue, strings.NewReader(testCase.inputBody))
			r.SetPathValue("id", testCase.idPathValue)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
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
