package cafe

import "omega/internal/types"

const (
	CreateFood types.Event = "food-create"
	UpdateFood types.Event = "food-update"
	DeleteFood types.Event = "food-delete"
	ListFood   types.Event = "food-list"
	ViewFood   types.Event = "food-view"
	ExcelFood  types.Event = "food-excel"

	CreateOrder   types.Event = "order-create"
	UpdateOrder   types.Event = "order-update"
	DeleteOrder   types.Event = "order-delete"
	ListOrder     types.Event = "order-list"
	ViewOrder     types.Event = "order-view"
	ExcelOrder    types.Event = "order-excel"
	MonthlyReport types.Event = "order-monthly-report"
)
