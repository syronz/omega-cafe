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

// FoodAPI for injecting food service
type FoodAPI struct {
	Service service.BasFoodServ
	Engine  *core.Engine
}

// ProvideFoodAPI for food is used in wire
func ProvideFoodAPI(c service.BasFoodServ) FoodAPI {
	return FoodAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a food by it's id
func (p *FoodAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var err error
	var food cafmodel.Food

	if food.ID, err = resp.GetRowID(c.Param("foodID"), "E1553982", cafterm.Food); err != nil {
		return
	}

	if food, err = p.Service.FindByID(food.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.ViewFood)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, cafterm.Food).
		JSON(food)
}

// List of foods
func (p *FoodAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, cafmodel.FoodTable, cafe.Domain)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.ListFood)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, cafterm.Foods).
		JSON(data)
}

// Create food
func (p *FoodAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var food, createdFood cafmodel.Food
	var err error

	if err = resp.Bind(&food, "E1588259", cafe.Domain, cafterm.Food); err != nil {
		return
	}

	if createdFood, err = p.Service.Create(food); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(cafe.CreateFood, food)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, cafterm.Food).
		JSON(createdFood)
}

// Update food
func (p *FoodAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var err error

	var food, foodBefore, foodUpdated cafmodel.Food

	if food.ID, err = resp.GetRowID(c.Param("foodID"), "E1582097", cafterm.Food); err != nil {
		return
	}

	if err = resp.Bind(&food, "E1576117", cafe.Domain, cafterm.Food); err != nil {
		return
	}

	if foodBefore, err = p.Service.FindByID(food.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	if foodUpdated, err = p.Service.Save(food); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.UpdateFood, foodBefore, food)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, cafterm.Food).
		JSON(foodUpdated)
}

// Delete food
func (p *FoodAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, cafe.Domain)
	var err error
	var food cafmodel.Food

	if food.ID, err = resp.GetRowID(c.Param("foodID"), "E1574329", cafterm.Food); err != nil {
		return
	}

	if food, err = p.Service.Delete(food.ID); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.DeleteFood, food)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, cafterm.Food).
		JSON()
}

// Excel generate excel files cafed on search
func (p *FoodAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, cafterm.Foods, cafe.Domain)

	foods, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("food")
	ex.AddSheet("Foods").
		AddSheet("Summary").
		Active("Foods").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "B", 15.3).
		SetColWidth("C", "C", 80).
		SetColWidth("D", "E", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Foods").
		WriteHeader("ID", "Name", "Resources", "Description", "Updated At").
		SetSheetFields("ID", "Name", "Resources", "Description", "UpdatedAt").
		WriteData(foods).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(cafe.ExcelFood)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
