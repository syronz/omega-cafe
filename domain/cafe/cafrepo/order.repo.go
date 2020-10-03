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

// OrderRepo for injecting engine
type OrderRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideOrderRepo is used in wire and initiate the Cols
func ProvideOrderRepo(engine *core.Engine) OrderRepo {
	return OrderRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(cafmodel.Order{}), cafmodel.OrderTable),
	}
}

// FindByID finds the order via its id
func (p *OrderRepo) FindByID(id types.RowID) (order cafmodel.Order, err error) {
	err = p.Engine.DB.Table(cafmodel.OrderTable).First(&order, id.ToUint64()).Error

	order.ID = id
	err = p.dbError(err, "E1572991", order, corterm.List)

	return
}

// List returns an array of orders
func (p *OrderRepo) List(params param.Param) (orders []cafmodel.Order, err error) {
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

	err = p.Engine.DB.Table(cafmodel.OrderTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&orders).Error

	err = p.dbError(err, "E1532861", cafmodel.Order{}, corterm.List)

	return
}

// Count of orders, mainly calls with List
func (p *OrderRepo) Count(params param.Param) (count uint64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1532288").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.DB.Table(cafmodel.OrderTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1539820", cafmodel.Order{}, corterm.List)
	return
}

// Save the order, in case it is not exist create it
func (p *OrderRepo) Save(order cafmodel.Order) (u cafmodel.Order, err error) {
	if err = p.Engine.DB.Table(cafmodel.OrderTable).Save(&order).Error; err != nil {
		err = p.dbError(err, "E1554817", order, corterm.Updated)
	}

	p.Engine.DB.Table(cafmodel.OrderTable).Where("id = ?", order.ID).Find(&u)
	return
}

// Create a order
func (p *OrderRepo) Create(order cafmodel.Order) (u cafmodel.Order, err error) {
	if err = p.Engine.DB.Table(cafmodel.OrderTable).Create(&order).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1553287", order, corterm.Created)
	}
	return
}

// Delete the order
func (p *OrderRepo) Delete(order cafmodel.Order) (err error) {
	if err = p.Engine.DB.Table(cafmodel.OrderTable).Unscoped().Delete(&order).Error; err != nil {
		err = p.dbError(err, "E1567392", order, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper datacafe error
func (p *OrderRepo) dbError(err error, code string, order cafmodel.Order, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, order.ID, cafterm.Orders)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(cafterm.Foods),
				dict.R(cafterm.Order), dict.R(action)).
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
