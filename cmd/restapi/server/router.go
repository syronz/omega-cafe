package server

import (
	"net/http"
	"omega/domain/base"
	"omega/domain/base/basmid"
	"omega/domain/cafe"
	"omega/internal/core"

	"github.com/gin-gonic/gin"
)

// Route trigger router and api methods
func Route(rg gin.RouterGroup, engine *core.Engine) {
	// Base Domain
	basAuthAPI := initAuthAPI(engine)
	basUserAPI := initUserAPI(engine)
	basRoleAPI := initRoleAPI(engine)
	basSettingAPI := initSettingAPI(engine)
	basActivityAPI := initActivityAPI(engine)

	// Html Domain
	htmErrDescAPI := initErrDescAPI(engine)

	// Cafe Domain
	cafFoodAPI := initFoodAPI(engine)
	cafOrderAPI := initOrderAPI(engine)
	cafOrderFoodAPI := initOrderFoodAPI(engine)

	rg.GET("/error-list", htmErrDescAPI.List)
	rg.StaticFS("/public", http.Dir("public"))

	rg.POST("/login", basAuthAPI.Login)

	rg.Use(basmid.AuthGuard(engine))
	access := basmid.NewAccessMid(engine)

	rg.POST("/logout", basAuthAPI.Logout)

	// Base Domain
	rg.GET("/temporary/token", basAuthAPI.TemporaryToken)

	rg.GET("/settings", access.Check(base.SettingRead), basSettingAPI.List)
	rg.GET("/settings/:settingID", access.Check(base.SettingRead), basSettingAPI.FindByID)
	rg.PUT("/settings/:settingID", access.Check(base.SettingWrite), basSettingAPI.Update)
	rg.GET("/excel/settings", access.Check(base.SettingExcel), basSettingAPI.Excel)

	rg.GET("/roles", access.Check(base.RoleRead), basRoleAPI.List)
	rg.GET("/roles/:roleID", access.Check(base.RoleRead), basRoleAPI.FindByID)
	rg.POST("/roles", access.Check(base.RoleWrite), basRoleAPI.Create)
	rg.PUT("/roles/:roleID", access.Check(base.RoleWrite), basRoleAPI.Update)
	rg.DELETE("/roles/:roleID", access.Check(base.RoleWrite), basRoleAPI.Delete)
	rg.GET("/excel/roles", access.Check(base.RoleExcel), basRoleAPI.Excel)

	rg.GET("/username/:username", access.Check(base.UserRead), basUserAPI.FindByUsername)
	rg.GET("/users", access.Check(base.UserRead), basUserAPI.List)
	rg.GET("/users/:userID", access.Check(base.UserRead), basUserAPI.FindByID)
	rg.POST("/users", access.Check(base.UserWrite), basUserAPI.Create)
	rg.PUT("/users/:userID", access.Check(base.UserWrite), basUserAPI.Update)
	rg.DELETE("/users/:userID", access.Check(base.UserWrite), basUserAPI.Delete)
	rg.GET("/excel/users", access.Check(base.UserExcel), basUserAPI.Excel)

	rg.GET("/activities", access.Check(base.ActivityAll), basActivityAPI.List)

	// Cafe Domain
	rg.GET("/foods", access.Check(cafe.FoodRead), cafFoodAPI.List)
	rg.GET("/foods/:foodID", access.Check(cafe.FoodRead), cafFoodAPI.FindByID)
	rg.POST("/foods", access.Check(cafe.FoodWrite), cafFoodAPI.Create)
	rg.PUT("/foods/:foodID", access.Check(cafe.FoodWrite), cafFoodAPI.Update)
	rg.DELETE("/foods/:foodID", access.Check(cafe.FoodWrite), cafFoodAPI.Delete)
	rg.GET("/excel/foods", access.Check(cafe.FoodExcel), cafFoodAPI.Excel)

	rg.GET("/orders", access.Check(cafe.OrderRead), cafOrderAPI.List)
	rg.GET("/orders/:orderID", access.Check(cafe.OrderRead), cafOrderAPI.FindByID)
	rg.POST("/orders", access.Check(cafe.OrderWrite), cafOrderAPI.Create)
	rg.PUT("/orders/:orderID", access.Check(cafe.OrderWrite), cafOrderAPI.Update)
	rg.DELETE("/orders/:orderID", access.Check(cafe.OrderWrite), cafOrderAPI.Delete)
	rg.GET("/excel/orders", access.Check(cafe.OrderExcel), cafOrderAPI.Excel)

	rg.GET("/order-food", access.Check(cafe.OrderRead), cafOrderFoodAPI.List)
	rg.GET("/order-food/:orderFoodID", access.Check(cafe.OrderRead), cafOrderFoodAPI.FindByID)
	// rg.POST("/order-food", access.Check(cafe.OrderWrite), cafOrderFoodAPI.Create)
	rg.PUT("/order-food/:orderFoodID", access.Check(cafe.OrderWrite), cafOrderFoodAPI.Update)
	rg.DELETE("/order-food/:orderFoodID", access.Check(cafe.OrderWrite), cafOrderFoodAPI.Delete)
	rg.GET("/excel/order-food", access.Check(cafe.OrderExcel), cafOrderFoodAPI.Excel)

}
