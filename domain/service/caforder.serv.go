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

// BasOrderServ for injecting auth cafrepo
type BasOrderServ struct {
	Repo   cafrepo.OrderRepo
	Engine *core.Engine
}

// ProvideBasOrderService for order is used in wire
func ProvideBasOrderService(p cafrepo.OrderRepo) BasOrderServ {
	return BasOrderServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting order by it's id
func (p *BasOrderServ) FindByID(id types.RowID) (order cafmodel.Order, err error) {
	if order, err = p.Repo.FindByID(id); err != nil {
		err = corerr.Tick(err, "E3243183", "can't fetch the order", id)
		return
	}

	return
}

// List of orders, it support pagination and search and return back count
func (p *BasOrderServ) List(params param.Param) (orders []cafmodel.Order,
	count uint64, err error) {

	if orders, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in orders list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in orders count")
	}

	return
}

// Create a order
func (p *BasOrderServ) Create(order cafmodel.Order) (createdOrder cafmodel.Order, err error) {

	if err = order.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E3298554", "validation failed in creating the order", order)
		return
	}

	clonedEngine := p.Engine.Clone()
	clonedEngine.DB = clonedEngine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			glog.LogError(fmt.Errorf("panic happened in transaction mode for %v",
				"orders, order_foods table"), "rollback recover")
			clonedEngine.DB.Rollback()
		}
	}()

	glog.Debug(order)

	orderRepo := cafrepo.ProvideOrderRepo(clonedEngine)
	if createdOrder, err = orderRepo.Create(order); err != nil {
		err = corerr.Tick(err, "E3242894", "order not created", order)
		clonedEngine.DB.Rollback()
		return
	}

	foodService := ProvideBasFoodService(cafrepo.ProvideFoodRepo(clonedEngine))

	orderFoodRepo := cafrepo.ProvideOrderFoodRepo(clonedEngine)
	for _, v := range order.Foods {
		v.OrderID = createdOrder.ID
		var food cafmodel.Food
		if food, err = foodService.FindByID(v.FoodID); err != nil {
			err = corerr.Tick(err, "E3266553", "food not exist for creation", v)
			clonedEngine.DB.Rollback()
			return
		}
		v.Price = food.Price
		if _, err = orderFoodRepo.Create(v); err != nil {
			err = corerr.Tick(err, "E3278321", "order-food not created", v)
			clonedEngine.DB.Rollback()
			return
		}
	}

	clonedEngine.DB.Commit()

	return
}

// Save a order, if it is exist update it, if not create it
func (p *BasOrderServ) Save(order cafmodel.Order) (savedOrder cafmodel.Order, err error) {
	if err = order.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E3237119", corerr.ValidationFailed, order)
		return
	}

	if savedOrder, err = p.Repo.Save(order); err != nil {
		err = corerr.Tick(err, "E3278742", "order not saved")
		return
	}

	BasAccessResetFullCache()
	return
}

// Delete order, it is soft delete
func (p *BasOrderServ) Delete(orderID types.RowID) (order cafmodel.Order, err error) {
	if order, err = p.FindByID(orderID); err != nil {
		err = corerr.Tick(err, "E3252861", "order not found for deleting")
		return
	}

	if err = p.Repo.Delete(order); err != nil {
		err = corerr.Tick(err, "E3217987", "order not deleted")
		return
	}

	BasAccessResetFullCache()
	return
}

// Excel is used for export excel file
func (p *BasOrderServ) Excel(params param.Param) (orders []cafmodel.Order, err error) {
	params.Limit = p.Engine.Envs.ToUint64(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", cafmodel.OrderTable)

	if orders, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E3267385", "cant generate the excel list for orders")
		return
	}

	return
}
