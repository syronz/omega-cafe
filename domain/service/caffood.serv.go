package service

import (
	"fmt"
	"omega/domain/cafe/cafmodel"
	"omega/domain/cafe/cafrepo"
	"omega/internal/core"
	"omega/internal/core/coract"
	"omega/internal/core/corerr"
	"omega/internal/param"
	"omega/internal/types"
	"omega/pkg/glog"
)

// BasFoodServ for injecting auth cafrepo
type BasFoodServ struct {
	Repo   cafrepo.FoodRepo
	Engine *core.Engine
}

// ProvideBasFoodService for food is used in wire
func ProvideBasFoodService(p cafrepo.FoodRepo) BasFoodServ {
	return BasFoodServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting food by it's id
func (p *BasFoodServ) FindByID(id types.RowID) (food cafmodel.Food, err error) {
	if food, err = p.Repo.FindByID(id); err != nil {
		err = corerr.Tick(err, "E1543183", "can't fetch the food", id)
		return
	}

	return
}

// List of foods, it support pagination and search and return back count
func (p *BasFoodServ) List(params param.Param) (foods []cafmodel.Food,
	count uint64, err error) {

	if foods, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in foods list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in foods count")
	}

	return
}

// Create a food
func (p *BasFoodServ) Create(food cafmodel.Food) (createdFood cafmodel.Food, err error) {

	if err = food.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1598554", "validation failed in creating the food", food)
		return
	}

	if createdFood, err = p.Repo.Create(food); err != nil {
		err = corerr.Tick(err, "E1542894", "food not created", food)
		return
	}

	return
}

// Save a food, if it is exist update it, if not create it
func (p *BasFoodServ) Save(food cafmodel.Food) (savedFood cafmodel.Food, err error) {
	if err = food.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E1537119", corerr.ValidationFailed, food)
		return
	}

	if savedFood, err = p.Repo.Save(food); err != nil {
		err = corerr.Tick(err, "E1578742", "food not saved")
		return
	}

	BasAccessResetFullCache()
	return
}

// Delete food, it is soft delete
func (p *BasFoodServ) Delete(foodID types.RowID) (food cafmodel.Food, err error) {
	if food, err = p.FindByID(foodID); err != nil {
		err = corerr.Tick(err, "E1552861", "food not found for deleting")
		return
	}

	if err = p.Repo.Delete(food); err != nil {
		err = corerr.Tick(err, "E1517987", "food not deleted")
		return
	}

	BasAccessResetFullCache()
	return
}

// Excel is used for export excel file
func (p *BasFoodServ) Excel(params param.Param) (foods []cafmodel.Food, err error) {
	params.Limit = p.Engine.Envs.ToUint64(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", cafmodel.FoodTable)

	if foods, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E1567385", "cant generate the excel list for foods")
		return
	}

	return
}
