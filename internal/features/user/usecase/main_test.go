package userusecase_test

import (
	"os"
	"testing"

	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	userentity "github.com/egotk/golang-advert-app/internal/features/user/entity"
)

func TestMain(m *testing.M) {
	if err := corevalidator.RegisterValidations(userentity.Validations(), nil); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}