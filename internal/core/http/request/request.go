package corehttprequest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
)

func Decode(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

func GetIntPathParam(key string, r *http.Request) (int64, error) {
	param := r.PathValue(key)
	if param == "" {
		return 0, fmt.Errorf(
			"no key='%s' in path params: %w",
			key,
			coreerrors.ErrInvalidArgument,
		)
	}

	value, err := strconv.ParseInt(param, 10, 64)
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

func GetIntQueryParam(key string, r *http.Request) (*int64, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	value, err := strconv.ParseInt(param, 10, 64)
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

func GetLimitOffsetQueryParams(r *http.Request) (*int64, *int64, error) {
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
