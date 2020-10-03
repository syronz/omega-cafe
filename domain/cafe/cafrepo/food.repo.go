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

// FoodRepo for injecting engine
type FoodRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideFoodRepo is used in wire and initiate the Cols
func ProvideFoodRepo(engine *core.Engine) FoodRepo {
	return FoodRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(cafmodel.Food{}), cafmodel.FoodTable),
	}
}

// FindByID finds the food via its id
func (p *FoodRepo) FindByID(id types.RowID) (food cafmodel.Food, err error) {
	err = p.Engine.DB.Table(cafmodel.FoodTable).First(&food, id.ToUint64()).Error

	food.ID = id
	err = p.dbError(err, "E1572991", food, corterm.List)

	return
}

// List returns an array of foods
func (p *FoodRepo) List(params param.Param) (foods []cafmodel.Food, err error) {
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

	err = p.Engine.DB.Table(cafmodel.FoodTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&foods).Error

	err = p.dbError(err, "E1532861", cafmodel.Food{}, corterm.List)

	return
}

// Count of foods, mainly calls with List
func (p *FoodRepo) Count(params param.Param) (count uint64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E1532288").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.DB.Table(cafmodel.FoodTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E1539820", cafmodel.Food{}, corterm.List)
	return
}

// Save the food, in case it is not exist create it
func (p *FoodRepo) Save(food cafmodel.Food) (u cafmodel.Food, err error) {
	if err = p.Engine.DB.Table(cafmodel.FoodTable).Save(&food).Error; err != nil {
		err = p.dbError(err, "E1554817", food, corterm.Updated)
	}

	p.Engine.DB.Table(cafmodel.FoodTable).Where("id = ?", food.ID).Find(&u)
	return
}

// Create a food
func (p *FoodRepo) Create(food cafmodel.Food) (u cafmodel.Food, err error) {
	if err = p.Engine.DB.Table(cafmodel.FoodTable).Create(&food).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E1553287", food, corterm.Created)
	}
	return
}

// Delete the food
func (p *FoodRepo) Delete(food cafmodel.Food) (err error) {
	if err = p.Engine.DB.Table(cafmodel.FoodTable).Unscoped().Delete(&food).Error; err != nil {
		err = p.dbError(err, "E1567392", food, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper datacafe error
func (p *FoodRepo) dbError(err error, code string, food cafmodel.Food, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, food.ID, cafterm.Foods)

	// case corerr.ForeignErr:
	// 	err = limberr.Take(err, code).
	// 		Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(cafterm.Users),
	// 			dict.R(cafterm.Food), dict.R(action)).
	// 		Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(cafterm.Food), food.Name).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, food.Name)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
