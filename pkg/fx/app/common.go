package app

import (
	"github.com/volmedo/padron/pkg/config/app"
	"github.com/volmedo/padron/pkg/fx/echo"
	"go.uber.org/fx"
)

func CommonModules(cfg app.AppConfig) fx.Option {
	var modules = []fx.Option{
		// Supply top level config, and it's sub-configs
		// this allows dependencies to be taken on, for example, app.IdentityConfig or app.ServerConfig
		// instead of needing to depend on the top level app.AppConfig
		fx.Supply(cfg),
		fx.Supply(cfg.Identity),
		fx.Supply(cfg.Server),

		echo.Module, // Provides Echo server with route registration
	}

	return fx.Module("common", modules...)
}
