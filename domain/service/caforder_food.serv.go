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

// BasOrderFoodServ for injecting auth cafrepo
type BasOrderFoodServ struct {
	Repo   cafrepo.OrderFoodRepo
	Engine *core.Engine
}

// ProvideBasOrderFoodService for order_food is used in wire
func ProvideBasOrderFoodService(p cafrepo.OrderFoodRepo) BasOrderFoodServ {
	return BasOrderFoodServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting order_food by it's id
func (p *BasOrderFoodServ) FindByID(id types.RowID) (order_food cafmodel.OrderFood, err error) {
	if order_food, err = p.Repo.FindByID(id); err != nil {
		err = corerr.Tick(err, "E2243183", "can't fetch the order_food", id)
		return
	}

	return
}

// List of order_foods, it support pagination and search and return back count
func (p *BasOrderFoodServ) List(params param.Param) (order_foods []cafmodel.OrderFood,
	count uint64, err error) {

	if order_foods, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in order_foods list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in order_foods count")
	}

	return
}

// Create a order_food
func (p *BasOrderFoodServ) Create(order_food cafmodel.OrderFood) (createdOrderFood cafmodel.OrderFood, err error) {

	if err = order_food.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E2298554", "validation failed in creating the order_food", order_food)
		return
	}

	if createdOrderFood, err = p.Repo.Create(order_food); err != nil {
		err = corerr.Tick(err, "E2242894", "order_food not created", order_food)
		return
	}

	return
}

// Save a order_food, if it is exist update it, if not create it
func (p *BasOrderFoodServ) Save(order_food cafmodel.OrderFood) (savedOrderFood cafmodel.OrderFood, err error) {
	if err = order_food.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E2237119", corerr.ValidationFailed, order_food)
		return
	}

	if savedOrderFood, err = p.Repo.Save(order_food); err != nil {
		err = corerr.Tick(err, "E2278742", "order_food not saved")
		return
	}

	BasAccessResetFullCache()
	return
}

// Delete order_food, it is soft delete
func (p *BasOrderFoodServ) Delete(order_foodID types.RowID) (order_food cafmodel.OrderFood, err error) {
	if order_food, err = p.FindByID(order_foodID); err != nil {
		err = corerr.Tick(err, "E2252861", "order_food not found for deleting")
		return
	}

	if err = p.Repo.Delete(order_food); err != nil {
		err = corerr.Tick(err, "E2217987", "order_food not deleted")
		return
	}

	BasAccessResetFullCache()
	return
}

// Excel is used for export excel file
func (p *BasOrderFoodServ) Excel(params param.Param) (order_foods []cafmodel.OrderFood, err error) {
	params.Limit = p.Engine.Envs.ToUint64(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", cafmodel.OrderFoodTable)

	if order_foods, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E2267385", "cant generate the excel list for order_foods")
		return
	}

	return
}
