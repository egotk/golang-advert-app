package corehttprequest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corevalidator "github.com/egotk/golang-advert-app/internal/core/validator"
)

type validatable interface {
	Validate() error
}

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	var err error

	value, ok := dest.(validatable)
	if ok {
		err = value.Validate()
	} else {
		validator := corevalidator.Instance()

		err = validator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"validate request: %v: %w",
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

func GetIntPathParam(key string, r *http.Request) (int, error) {
	param := r.PathValue(key)
	if param == "" {
		return 0, fmt.Errorf(
			"no key='%s' in path params: %w",
			key,
			coreerrors.ErrInvalidArgument,
		)
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf(
			"path param='%s' by key='%s' is not a valid int: %v: %w",
			param,
			key,
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	return value, nil
}

func GetIntQueryParam(key string, r *http.Request) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not a valid int: %v: %w",
			param,
			key,
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	return &value, nil
}

func GetStringQueryParam(key string, r *http.Request) *string {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}

	return &param
}

func GetTimeQueryParam(key string, r *http.Request) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	time, err := time.Parse("2006-01-02 15:04:05", param)
	if err != nil {
		return nil, fmt.Errorf("get time query param: %w", err)
	}

	return &time, nil
}

func GetLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := GetIntQueryParam(limitQueryParamKey, r)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := GetIntQueryParam(offsetQueryParamKey, r)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
