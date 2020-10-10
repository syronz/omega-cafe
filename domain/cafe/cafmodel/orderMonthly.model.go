package cafmodel

type OrderMonthly struct {
	Month    string `json:"month"`
	Total    int    `json:"total"`
	Discount int    `json:"discount"`
	SubTotal int    `json:"sub_total"`
}
