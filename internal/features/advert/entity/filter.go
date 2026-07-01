package advertentity

import "time"

type Sort string

const (
	SortByPrice Sort = "price"
	SortByViews Sort = "views"
	SortByDate  Sort = "date"
)

type Order string

const (
	OrderAsc  Order = "ASC"
	OrderDesc Order = "DESC"
)

type Filter struct {
	UserID      *int64
	Title       *string
	Description *string
	MinPrice    *int64
	MaxPrice    *int64
	CategoryID  *int64
	Status      *Status
	FromDate    *time.Time
	ToDate      *time.Time
	Sort        *Sort
	Order       *Order
}
