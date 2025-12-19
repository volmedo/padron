package ucan

import (
	"github.com/alanshaw/ucantone/server"
	logging "github.com/ipfs/go-log/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/volmedo/padron/pkg/config/app"
	"github.com/volmedo/padron/pkg/fx/blob"
	echofx "github.com/volmedo/padron/pkg/fx/echo"
	"github.com/volmedo/padron/pkg/ucan"
)

var log = logging.Logger("fx/ucan")

var _ echofx.RouteRegistrar = (*Server)(nil)

type Server struct {
	ucanServer *server.HTTPServer
}

var Module = fx.Module("ucan/server",
	fx.Provide(
		NewServer,
		fx.Annotate(
			NewServer,
			fx.As(new(echofx.RouteRegistrar)),
			fx.ResultTags(`group:"route_registrar"`),
		),
	),
	blob.Module,
)

type Params struct {
	fx.In
	Identity app.IdentityConfig
	Handlers []*ucan.Handler     `group:"ucan_handlers"`
	Options  []server.HTTPOption `group:"ucan_options"`
}

func NewServer(p Params) (*Server, error) {
	ucanSvr := server.NewHTTP(p.Identity.Signer, p.Options...)
	log.Infof("Registering %d UCAN handlers", len(p.Handlers))
	for _, h := range p.Handlers {
		log.Infof("Registering %q UCAN handler", h.Capability.Command())
		ucanSvr.Handle(h.Capability, h.Handler)
	}
	return &Server{ucanSvr}, nil
}

func (s *Server) RegisterRoutes(e *echo.Echo) {
	e.POST("/", echo.WrapHandler(s.ucanServer))
}
