package advertusecase_test

import (
	"os"
	"testing"

	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
)

func TestMain(m *testing.M) {
	if err := corevalidator.RegisterValidations(nil, advertentity.StructValidations()); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}