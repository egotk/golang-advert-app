package categoryusecase_test

import (
	"context"
	"errors"
	"testing"

	categoryentity "github.com/egotk/golang-advert-app/internal/features/category/entity"
	categoryusecase "github.com/egotk/golang-advert-app/internal/features/category/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_List(t *testing.T) {
	type listMockBehavior func(mr *Mockrepo)

	mockErr := errors.New("boom")

	testTable := []struct {
		name               string
		mockBehavior       listMockBehavior
		expectedErrIs      error
		expectedCategories []categoryentity.Category
	}{
		{
			name: "OK",
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().List(gomock.Any()).Return(
					[]categoryentity.Category{{ID: 1, Name: "Electronics"}}, nil)
			},
			expectedCategories: []categoryentity.Category{{ID: 1, Name: "Electronics"}},
		},
		{
			name: "repo error",
			mockBehavior: func(mr *Mockrepo) {
				mr.EXPECT().List(gomock.Any()).Return(nil, mockErr)
			},
			expectedErrIs: mockErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mr := NewMockrepo(c)
			testCase.mockBehavior(mr)

			uc := categoryusecase.New(mr)

			categories, err := uc.List(context.Background())

			if testCase.expectedErrIs != nil {
				assert.ErrorIs(t, err, testCase.expectedErrIs)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCategories, categories)
		})
	}
}
