package blob

import (
	"go.uber.org/fx"

	blobsvc "github.com/volmedo/padron/pkg/service/blob"
	"github.com/volmedo/padron/pkg/ucan/blob"
)

var Module = fx.Module("blob",
	fx.Provide(
		NewBlobService,
		fx.Annotate(
			blob.NewBlobAllocateHandler,
			fx.ResultTags(`group:"ucan_handlers"`),
		),
	),
)

func NewBlobService() (*blobsvc.Service, error) {
	return blobsvc.NewService(), nil
}
