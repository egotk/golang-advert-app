package categoryhttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	categoryhttp "github.com/egotk/golang-advert-app/internal/features/category/controller/rest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_Delete(t *testing.T) {
	type deleteMockBehavior func(muc *MockuseCase, id int64)

	testTable := []struct {
		name               string
		idPathValue        string
		inputID            int64
		mockBehavior       deleteMockBehavior
		expectedStatusCode int
	}{
		{
			name:        "OK",
			idPathValue: "1",
			inputID:     1,
			mockBehavior: func(muc *MockuseCase, id int64) {
				muc.EXPECT().Delete(gomock.Any(), id).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
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
				muc.EXPECT().Delete(gomock.Any(), id).Return(coreerrors.ErrNotFound)
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

			controller := categoryhttp.New(muc)

			r := httptest.NewRequest(http.MethodDelete, "/adverts/categories/"+testCase.idPathValue, nil)
			r.SetPathValue("id", testCase.idPathValue)

			ctx := corezaplogger.ToContext(r.Context(), &corezaplogger.Logger{Logger: zap.NewNop()})
			r = r.WithContext(ctx)

			rw := httptest.NewRecorder()

			controller.Delete(rw, r)

			assert.Equal(t, testCase.expectedStatusCode, rw.Code)
		})
	}
}
