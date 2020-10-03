package startoff

import (
	"fmt"
	"omega/domain/base/basmodel"
	"omega/domain/cafe/cafmodel"
	"omega/internal/core"
)

// Migrate the database for creating tables
func Migrate(engine *core.Engine) {

	// Base Domain
	engine.DB.Table(basmodel.SettingTable).AutoMigrate(&basmodel.Setting{})
	engine.DB.Table(basmodel.RoleTable).AutoMigrate(&basmodel.Role{})
	engine.DB.Table(basmodel.UserTable).AutoMigrate(&basmodel.User{}).
		AddForeignKey("role_id", fmt.Sprintf("%v(id)", basmodel.RoleTable), "RESTRICT", "RESTRICT")
	engine.ActivityDB.Table(basmodel.ActivityTable).AutoMigrate(&basmodel.Activity{})

	// Cafe Domain
	engine.DB.Table(cafmodel.FoodTable).AutoMigrate(&cafmodel.Food{})
	engine.DB.Table(cafmodel.OrderTable).AutoMigrate(&cafmodel.Order{})
	engine.DB.Table(cafmodel.OrderFoodTable).AutoMigrate(&cafmodel.OrderFood{}).
		AddForeignKey("order_id", "caf_orders(id)", "CASCADE", "CASCADE").
		AddForeignKey("food_id", "caf_foods(id)", "RESTRICT", "RESTRICT")
}
