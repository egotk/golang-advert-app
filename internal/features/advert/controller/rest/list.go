package advertrest

import (
	"fmt"
	"net/http"

	corehttprequest "github.com/egotk/golang-advert-app/internal/core/http/request"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	advertentity "github.com/egotk/golang-advert-app/internal/features/advert/entity"
	advertusecase "github.com/egotk/golang-advert-app/internal/features/advert/usecase"
)

func (c *Controller) list(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := corezaplogger.FromContext(ctx)
	responseHandler := corehttpresponse.New(log, rw)

	limit, offset, err := corehttprequest.GetLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit/offset query param")

		return
	}

	filter, err := getListFilterQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get list filter query params")

		return
	}

	userID, userRole, err := corejwt.UserInfoFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'UserInfo' from JWT")

		return
	}

	dto := advertusecase.ListDTO{
		UserID:   userID,
		UserRole: userRole,
		Limit:    limit,
		Offset:   offset,
		Filter:   filter,
	}

	count, adverts, err := c.useCase.List(ctx, dto)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to list adverts")

		return
	}

	response := advertsResponseFromEntities(count, adverts)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getListFilterQueryParams(r *http.Request) (advertentity.Filter, error) {
	const (
		titleKey       = "title"
		descriptionKey = "description"
		minPriceKey    = "minPrice"
		maxPriceKey    = "maxPrice"
		categoryIDKey  = "categoryID"
		statusKey      = "status"
		fromDateKey    = "fromDate"
		toDateKey      = "toDate"
		sortKey        = "sort"
		orderKey       = "order"
	)

	title := corehttprequest.GetStringQueryParam(titleKey, r)
	description := corehttprequest.GetStringQueryParam(descriptionKey, r)

	minPrice, err := corehttprequest.GetIntQueryParam(minPriceKey, r)
	if err != nil {
		return advertentity.Filter{}, fmt.Errorf("get 'minPrice' query param: %w", err)
	}

	maxPrice, err := corehttprequest.GetIntQueryParam(maxPriceKey, r)
	if err != nil {
		return advertentity.Filter{}, fmt.Errorf("get 'maxPrice' query param: %w", err)
	}

	categoryID, err := corehttprequest.GetIntQueryParam(categoryIDKey, r)
	if err != nil {
		return advertentity.Filter{}, fmt.Errorf("get 'categoryID' query param: %w", err)
	}

	status := corehttprequest.GetStringQueryParam(statusKey, r)

	fromDate, err := corehttprequest.GetTimeQueryParam(fromDateKey, r)
	if err != nil {
		return advertentity.Filter{}, fmt.Errorf("get 'fromDate' query param: %w", err)
	}

	toDate, err := corehttprequest.GetTimeQueryParam(toDateKey, r)
	if err != nil {
		return advertentity.Filter{}, fmt.Errorf("get 'toDate' query param: %w", err)
	}

	sort := corehttprequest.GetStringQueryParam(sortKey, r)

	order := corehttprequest.GetStringQueryParam(orderKey, r)

	filter := advertentity.Filter{
		Title:       title,
		Description: description,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		CategoryID:  categoryID,
		Status:      (*advertentity.Status)(status),
		FromDate:    fromDate,
		ToDate:      toDate,
		Sort:        (*advertentity.Sort)(sort),
		Order:       (*advertentity.Order)(order),
	}

	return filter, nil
}
