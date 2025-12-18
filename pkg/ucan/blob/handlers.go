package blob

import (
	"fmt"

	"github.com/alanshaw/libracha/capabilities"
	blobcap "github.com/alanshaw/libracha/capabilities/blob"
	"github.com/alanshaw/ucantone/execution/bindexec"
	logging "github.com/ipfs/go-log/v2"

	blobsvc "github.com/volmedo/padron/pkg/service/blob"
	"github.com/volmedo/padron/pkg/ucan"
)

var log = logging.Logger("ucan/blob")

const maxUploadSize = 127 * (1 << 25)

func NewBlobAllocateHandler(svc *blobsvc.Service) *ucan.Handler {
	return &ucan.Handler{
		Capability: blobcap.Allocate,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*blobcap.AllocateArguments]) (*bindexec.Response[*blobcap.AllocateOK], error) {
				args := req.Task().BindArguments()
				log.Debugf("%+v", args)

				// enforce max upload size requirements
				if args.Blob.Size > maxUploadSize {
					return nil, fmt.Errorf("blob size %d exceeds maximum upload size of %d bytes", args.Blob.Size, maxUploadSize)
				}

				size, address, err := svc.Allocate(
					req.Context(),
					req.Invocation().Subject().DID(),
					blobsvc.Blob{
						Digest: args.Blob.Digest,
						Size:   args.Blob.Size,
					},
					req.Invocation().Link(),
				)
				if err != nil {
					return nil, fmt.Errorf("allocation failed: %w", err)
				}

				hdrs := make(map[string]string, len(address.Headers))
				for k, v := range address.Headers {
					hdrs[k] = v[0]
				}

				var addr *blobcap.BlobAddress
				if address != nil {
					addr = &blobcap.BlobAddress{
						URL:     capabilities.CborURL(*address.URL),
						Headers: hdrs,
						Expires: capabilities.CborTime(address.Expires),
					}
				}

				ok := &blobcap.AllocateOK{
					Size:    size,
					Address: addr,
				}

				return bindexec.NewResponse(bindexec.WithSuccess(ok))
			},
		),
	}
}

func NewBlobAcceptHandler(svc *blobsvc.Service) *ucan.Handler {
	return &ucan.Handler{
		Capability: blobcap.Accept,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*blobcap.AcceptArguments]) (*bindexec.Response[*blobcap.AcceptOK], error) {
				args := req.Task().BindArguments()
				log.Debugf("%+v", args)

				locCommitment, err := svc.Accept(
					req.Context(),
					req.Invocation().Subject().DID(),
					blobsvc.Blob{
						Digest: args.Blob.Digest,
						Size:   args.Blob.Size,
					},
					req.Invocation().Link(),
				)
				if err != nil {
					return nil, fmt.Errorf("accept failed: %w", err)
				}

				ok := &blobcap.AcceptOK{
					Site: locCommitment.Link(),
				}
				
				return bindexec.NewResponse(bindexec.WithSuccess(ok))
			},
		),
	}
}
