package app

import (
	"go.uber.org/fx"

	"github.com/volmedo/padron/pkg/config/app"
	"github.com/volmedo/padron/pkg/fx/echo"
	"github.com/volmedo/padron/pkg/fx/identity"
	"github.com/volmedo/padron/pkg/fx/store"
)

func CommonModules(cfg app.AppConfig) fx.Option {
	var modules = []fx.Option{
		// Supply top level config, and it's sub-configs
		// this allows dependencies to be taken on, for example, app.IdentityConfig or app.ServerConfig
		// instead of needing to depend on the top level app.AppConfig
		fx.Supply(cfg),
		fx.Supply(cfg.Identity),
		fx.Supply(cfg.Server),
		fx.Supply(cfg.Stores),

		identity.Module, // Provides principal.Signer
		echo.Module,     // Provides Echo server with route registration
	}

	if cfg.Stores.DataDir == "" {
		modules = append(modules, store.MemoryStoreModule)
	} else {
		modules = append(modules, store.FileSystemStoreModule)
	}

	return fx.Module("common", modules...)
}
