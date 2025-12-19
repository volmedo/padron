package root

import (
	"github.com/alanshaw/ucantone/principal"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	echofx "github.com/volmedo/padron/pkg/fx/echo"
	"github.com/volmedo/padron/pkg/server"
)

var Module = fx.Module("root-handler",
	fx.Provide(
		fx.Annotate(
			NewRootHandler,
			fx.As(new(echofx.RouteRegistrar)),
			fx.ResultTags(`group:"route_registrar"`),
		),
	),
)

var _ echofx.RouteRegistrar = (*Handler)(nil)

type Handler struct {
	id principal.Signer
}

func NewRootHandler(id principal.Signer) *Handler {
	return &Handler{id: id}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", echo.WrapHandler(server.NewRootHandler(h.id)))
}
