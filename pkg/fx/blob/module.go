package blob

import (
	"github.com/alanshaw/ucantone/principal"
	"github.com/labstack/echo/v4"
	"github.com/storacha/piri/pkg/store/acceptancestore"
	"github.com/storacha/piri/pkg/store/allocationstore"
	"github.com/storacha/piri/pkg/store/blobstore"
	"go.uber.org/fx"

	"github.com/volmedo/padron/pkg/config/app"
	echofx "github.com/volmedo/padron/pkg/fx/echo"
	blobsvr "github.com/volmedo/padron/pkg/server/blob"
	blobsvc "github.com/volmedo/padron/pkg/service/blob"
	blobucan "github.com/volmedo/padron/pkg/ucan/blob"
)

var Module = fx.Module("blob",
	fx.Provide(
		NewBlobService,
		fx.Annotate(
			blobucan.NewBlobAllocateHandler,
			fx.ResultTags(`group:"ucan_handlers"`),
		),
		fx.Annotate(
			blobucan.NewBlobAcceptHandler,
			fx.ResultTags(`group:"ucan_handlers"`),
		),
		fx.Annotate(
			NewBlobServer,
			fx.As(new(echofx.RouteRegistrar)),
			fx.ResultTags(`group:"route_registrar"`),
		),
	),
)

func NewBlobService(
	cfg app.ServerConfig,
	id principal.Signer,
	blobs blobstore.Blobstore,
	allocs allocationstore.AllocationStore,
	acceptances acceptancestore.AcceptanceStore,
) *blobsvc.Service {
	return blobsvc.NewService(id, cfg.PublicURL, blobs, allocs, acceptances)
}

var _ echofx.RouteRegistrar = (*Server)(nil)

type Server struct {
	allocs allocationstore.AllocationStore
	blobs  blobstore.Blobstore
}

func NewBlobServer(allocs allocationstore.AllocationStore, blobs blobstore.Blobstore) *Server {
	return &Server{
		allocs: allocs,
		blobs:  blobs,
	}
}

func (srv *Server) RegisterRoutes(e *echo.Echo) {
	e.GET("/blob/:blob", blobsvr.NewBlobGetHandler(srv.blobs))
	e.PUT("/blob/:blob", blobsvr.NewBlobPutHandler(srv.allocs, srv.blobs))
}
