package cafrepo

import (
	"omega/domain/cafe/cafmodel"
	"omega/domain/cafe/message/cafterm"
	"omega/internal/core"
	"omega/internal/core/corerr"
	"omega/internal/core/corterm"
	"omega/internal/core/validator"
	"omega/internal/param"
	"omega/internal/types"
	"omega/pkg/dict"
	"omega/pkg/helper"
	"omega/pkg/limberr"
	"reflect"
)

// OrderFoodRepo for injecting engine
type OrderFoodRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideOrderFoodRepo is used in wire and initiate the Cols
func ProvideOrderFoodRepo(engine *core.Engine) OrderFoodRepo {
	return OrderFoodRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(cafmodel.OrderFood{}), cafmodel.OrderFoodTable),
	}
}

// FindByID finds the order_food via its id
func (p *OrderFoodRepo) FindByID(id types.RowID) (order_food cafmodel.OrderFood, err error) {
	err = p.Engine.DB.Table(cafmodel.OrderFoodTable).First(&order_food, id.ToUint64()).Error

	order_food.ID = id
	err = p.dbError(err, "E1572991", order_food, corterm.List)

	return
}

// List returns an array of order_foods
func (p *OrderFoodRepo) List(params param.Param) (order_foods []cafmodel.OrderFood, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E1584438").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1532278").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.DB.Table(cafmodel.OrderFoodTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&order_foods).Error

	err = p.dbError(err, "E1532861", cafmodel.OrderFood{}, corterm.List)

	return
}

// Count of order_foods, mainly calls with List
func (p *OrderFoodRepo) Count(params param.Param) (count uint64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1532288").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.DB.Table(cafmodel.OrderFoodTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1539820", cafmodel.OrderFood{}, corterm.List)
	return
}

// Save the order_food, in case it is not exist create it
func (p *OrderFoodRepo) Save(order_food cafmodel.OrderFood) (u cafmodel.OrderFood, err error) {
	if err = p.Engine.DB.Table(cafmodel.OrderFoodTable).Save(&order_food).Error; err != nil {
		err = p.dbError(err, "E1554817", order_food, corterm.Updated)
	}

	p.Engine.DB.Table(cafmodel.OrderFoodTable).Where("id = ?", order_food.ID).Find(&u)
	return
}

// Create a order_food
func (p *OrderFoodRepo) Create(order_food cafmodel.OrderFood) (u cafmodel.OrderFood, err error) {
	if err = p.Engine.DB.Table(cafmodel.OrderFoodTable).Create(&order_food).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1553287", order_food, corterm.Created)
	}
	return
}

// Delete the order_food
func (p *OrderFoodRepo) Delete(order_food cafmodel.OrderFood) (err error) {
	if err = p.Engine.DB.Table(cafmodel.OrderFoodTable).Unscoped().Delete(&order_food).Error; err != nil {
		err = p.dbError(err, "E1567392", order_food, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper datacafe error
func (p *OrderFoodRepo) dbError(err error, code string, order_food cafmodel.OrderFood, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, order_food.ID, cafterm.Foods)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(cafterm.Foods),
				dict.R(cafterm.Food), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
