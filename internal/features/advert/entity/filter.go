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
	Title       *string `validate:"omitempty,min=1,max=100"`
	Description *string `validate:"omitempty,min=1,max=1500"`
	MinPrice    *int64  `validate:"omitempty,gte=0"`
	MaxPrice    *int64  `validate:"omitempty,gte=0"`
	CategoryID  *int64  `validate:"omitempty,gte=1"`
	Status      *Status `validate:"omitempty,oneof=initial active rejected blocked archived"`
	FromDate    *time.Time
	ToDate      *time.Time
	Sort        *Sort  `validate:"omitempty,oneof=price views date"`
	Order       *Order `validate:"omitempty,oneof=ASC DESC"`
}
