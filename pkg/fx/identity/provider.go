package identity

import (
	"github.com/alanshaw/ucantone/principal"
	"go.uber.org/fx"

	"github.com/volmedo/padron/pkg/config/app"
)

var Module = fx.Module("identity",
	fx.Provide(ProvideIdentity),
)

// ProvideIdentity extracts the principal signer from the app config
func ProvideIdentity(cfg app.AppConfig) principal.Signer {
	return cfg.Identity.Signer
}
