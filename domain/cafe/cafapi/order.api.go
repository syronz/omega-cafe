package cafapi

import (
	"net/http"
	"omega/domain/cafe"
	"omega/domain/cafe/cafmodel"
	"omega/domain/cafe/message/cafterm"
	"omega/domain/service"
	"omega/internal/core"
	"omega/internal/core/corterm"
	"omega/internal/response"
	"omega/pkg/excel"

	"github.com/gin-gonic/gin"
)

// OrderAPI for injecting order service
type OrderAPI struct {
	Service service.BasOrderServ
	Engine  *core.Engine
}

// ProvideOrderAPI for order is used in wire
func ProvideOrderAPI(c service.BasOrderServ) OrderAPI {
	return OrderAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a order by it's id
func (p *OrderAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var err error
	var order cafmodel.Order

	if order.ID, err = resp.GetRowID(c.Param("orderID"), "E1553982", cafterm.Order); err != nil {
		return
	}

	if order, err = p.Service.FindByID(order.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.ViewOrder)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, cafterm.Order).
		JSON(order)
}

// List of orders
func (p *OrderAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, cafmodel.OrderTable, cafe.Domain)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.ListOrder)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, cafterm.Orders).
		JSON(data)
}

// Create order
func (p *OrderAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var order, createdOrder cafmodel.Order
	var err error

	if err = resp.Bind(&order, "E1588259", cafe.Domain, cafterm.Order); err != nil {
		return
	}

	if createdOrder, err = p.Service.Create(order); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(cafe.CreateOrder, order)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, cafterm.Order).
		JSON(createdOrder)
}

// Update order
func (p *OrderAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var err error

	var order, orderBefore, orderUpdated cafmodel.Order

	if order.ID, err = resp.GetRowID(c.Param("orderID"), "E1582097", cafterm.Order); err != nil {
		return
	}

	if err = resp.Bind(&order, "E1576117", cafe.Domain, cafterm.Order); err != nil {
		return
	}

	if orderBefore, err = p.Service.FindByID(order.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	if orderUpdated, err = p.Service.Save(order); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.UpdateOrder, orderBefore, order)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, cafterm.Order).
		JSON(orderUpdated)
}

// Delete order
func (p *OrderAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var err error
	var order cafmodel.Order

	if order.ID, err = resp.GetRowID(c.Param("orderID"), "E1574329", cafterm.Order); err != nil {
		return
	}

	if order, err = p.Service.Delete(order.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.DeleteOrder, order)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, cafterm.Order).
		JSON()
}

// Excel generate excel files cafed on search
func (p *OrderAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, cafterm.Orders, cafe.Domain)

	orders, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("order")
	ex.AddSheet("Orders").
		AddSheet("Summary").
		Active("Orders").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "B", 15.3).
		SetColWidth("C", "C", 80).
		SetColWidth("D", "E", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Orders").
		WriteHeader("ID", "Name", "Resources", "Description", "Updated At").
		SetSheetFields("ID", "Name", "Resources", "Description", "UpdatedAt").
		WriteData(orders).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.ExcelOrder)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}

// Print generate html to be printed
func (p *OrderAPI) Print(c *gin.Context) {
	// resp := response.New(p.Engine, c, cafe.Domain)

	// resp.Status(http.StatusOK).
	// 	MessageT(corterm.VInfo, cafterm.Order).
	// 	JSON(gin.H{"ok": "no"})

	data := gin.H{
		"title": "hello diako",
		"companyInfo": gin.H{
			"name": "footar",
		},
		"order": gin.H{
			"customer":   "Diako Amir",
			"phone":      "",
			"address":    "",
			"id":         323,
			"created_at": "2020-10-15 13:17:28",
			"created_by": 16,
		},
		"foods": []cafmodel.OrderFood{
			{
				Food:  "chicken",
				Qty:   3,
				Price: 1000,
				Total: 3000,
			},
			{
				Food:  "sandwitch",
				Qty:   1,
				Price: 1500,
				Total: 1500,
			},
		},
		"dict": gin.H{
			"Agent":      "Agent",
			"Food":       "Food",
			"Qty":        "Qty",
			"Price":      "Price",
			"Total":      "Total",
			"SubTotal":   "Sub Total",
			"Discount":   "Discount",
			"GrandTotal": "Grand Total",
		},
	}

	c.HTML(http.StatusOK, "order.tmpl", data)
}
