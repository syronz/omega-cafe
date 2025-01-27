package cafe

import "omega/internal/types"

const (
	Domain string = "cafe"

	FoodRead    types.Resource = "food:read"
	FoodWrite   types.Resource = "food:write"
	FoodExcel   types.Resource = "food:excel"
	FoodConsume types.Resource = "food-consume:report"

	OrderRead          types.Resource = "order:read"
	OrderWrite         types.Resource = "order:write"
	OrderExcel         types.Resource = "order:excel"
	OrderMonthlyReport types.Resource = "order-monthly:report"
)
