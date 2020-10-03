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

// OrderFoodAPI for injecting order_food service
type OrderFoodAPI struct {
	Service service.BasOrderFoodServ
	Engine  *core.Engine
}

// ProvideOrderFoodAPI for order_food is used in wire
func ProvideOrderFoodAPI(c service.BasOrderFoodServ) OrderFoodAPI {
	return OrderFoodAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a order_food by it's id
func (p *OrderFoodAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var err error
	var order_food cafmodel.OrderFood

	if order_food.ID, err = resp.GetRowID(c.Param("order_foodID"), "E1553978", cafterm.Food); err != nil {
		return
	}

	if order_food, err = p.Service.FindByID(order_food.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.ViewOrder)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, cafterm.Food).
		JSON(order_food)
}

// List of order_foods
func (p *OrderFoodAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, cafmodel.OrderFoodTable, cafe.Domain)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, cafterm.Foods).
		JSON(data)
}

// Create order_food
func (p *OrderFoodAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var order_food, createdOrderFood cafmodel.OrderFood
	var err error

	if err = resp.Bind(&order_food, "E1588256", cafe.Domain, cafterm.Food); err != nil {
		return
	}

	if createdOrderFood, err = p.Service.Create(order_food); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, cafterm.Food).
		JSON(createdOrderFood)
}

// Update order_food
func (p *OrderFoodAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var err error

	var order_food, order_foodBefore, order_foodUpdated cafmodel.OrderFood

	if order_food.ID, err = resp.GetRowID(c.Param("order_foodID"), "E1582094", cafterm.Food); err != nil {
		return
	}

	if err = resp.Bind(&order_food, "E1576114", cafe.Domain, cafterm.Food); err != nil {
		return
	}

	if order_foodBefore, err = p.Service.FindByID(order_food.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	if order_foodUpdated, err = p.Service.Save(order_food); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.UpdateOrder, order_foodBefore, order_food)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, cafterm.Food).
		JSON(order_foodUpdated)
}

// Delete order_food
func (p *OrderFoodAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var err error
	var order_food cafmodel.OrderFood

	if order_food.ID, err = resp.GetRowID(c.Param("order_foodID"), "E1574326", cafterm.Food); err != nil {
		return
	}

	if order_food, err = p.Service.Delete(order_food.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.DeleteOrder, order_food)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, cafterm.Food).
		JSON()
}

// Excel generate excel files cafed on search
func (p *OrderFoodAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, cafterm.Foods, cafe.Domain)

	order_foods, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("order_food")
	ex.AddSheet("OrderFoods").
		AddSheet("Summary").
		Active("OrderFoods").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "B", 15.3).
		SetColWidth("C", "C", 80).
		SetColWidth("D", "E", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("OrderFoods").
		WriteHeader("ID", "Name", "Resources", "Description", "Updated At").
		SetSheetFields("ID", "Name", "Resources", "Description", "UpdatedAt").
		WriteData(order_foods).
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
