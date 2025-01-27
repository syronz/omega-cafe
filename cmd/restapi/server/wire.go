// +build wireinject

package server

import (
	"omega/domain/base/basapi"
	"omega/domain/base/basrepo"
	"omega/domain/cafe/cafapi"
	"omega/domain/cafe/cafrepo"
	"omega/domain/html/htmapi"
	"omega/domain/service"

	"omega/internal/core"

	"github.com/google/wire"
)

// Base Domain
func initSettingAPI(e *core.Engine) basapi.SettingAPI {
	wire.Build(basrepo.ProvideSettingRepo, service.ProvideBasSettingService,
		basapi.ProvideSettingAPI)
	return basapi.SettingAPI{}
}

func initRoleAPI(e *core.Engine) basapi.RoleAPI {
	wire.Build(basrepo.ProvideRoleRepo, service.ProvideBasRoleService,
		basapi.ProvideRoleAPI)
	return basapi.RoleAPI{}
}

func initUserAPI(engine *core.Engine) basapi.UserAPI {
	wire.Build(basrepo.ProvideUserRepo, service.ProvideBasUserService, basapi.ProvideUserAPI)
	return basapi.UserAPI{}
}

func initAuthAPI(e *core.Engine) basapi.AuthAPI {
	wire.Build(service.ProvideBasAuthService, basapi.ProvideAuthAPI)
	return basapi.AuthAPI{}
}

func initActivityAPI(engine *core.Engine) basapi.ActivityAPI {
	wire.Build(basrepo.ProvideActivityRepo, service.ProvideBasActivityService, basapi.ProvideActivityAPI)
	return basapi.ActivityAPI{}
}

// Html Domain
func initErrDescAPI(e *core.Engine) htmapi.ErrDescAPI {
	wire.Build(htmapi.GenErrDescAPI)
	return htmapi.ErrDescAPI{}
}

// Cafe domain
func initFoodAPI(e *core.Engine) cafapi.FoodAPI {
	wire.Build(cafrepo.ProvideFoodRepo, service.ProvideBasFoodService,
		cafapi.ProvideFoodAPI)
	return cafapi.FoodAPI{}
}

func initOrderAPI(e *core.Engine) cafapi.OrderAPI {
	wire.Build(cafrepo.ProvideOrderRepo, service.ProvideBasOrderService,
		cafapi.ProvideOrderAPI)
	return cafapi.OrderAPI{}
}

func initOrderFoodAPI(e *core.Engine) cafapi.OrderFoodAPI {
	wire.Build(cafrepo.ProvideOrderFoodRepo, service.ProvideBasOrderFoodService,
		cafapi.ProvideOrderFoodAPI)
	return cafapi.OrderFoodAPI{}
}
