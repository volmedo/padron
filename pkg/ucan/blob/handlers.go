package blob

import (
	"net/url"
	"time"

	"github.com/alanshaw/libracha/capabilities"
	blobcap "github.com/alanshaw/libracha/capabilities/blob"
	"github.com/alanshaw/ucantone/execution/bindexec"
	logging "github.com/ipfs/go-log/v2"

	blobsvc "github.com/volmedo/padron/pkg/service/blob"
	"github.com/volmedo/padron/pkg/ucan"
)

var log = logging.Logger("ucan/blob")

const serverURL = "http://localhost:3000"

func NewBlobAllocateHandler(svc *blobsvc.Service) *ucan.Handler {
	return &ucan.Handler{
		Capability: blobcap.Allocate,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*blobcap.AllocateArguments]) (*bindexec.Response[*blobcap.AllocateOK], error) {
				log.Debugf("%+v", req.Task().BindArguments())
				args := req.Task().BindArguments()
				// site, err := svc.Allocate(req.Context(), digest)
				// if err != nil {
				// 	return nil, err
				// }
				url, err := url.Parse(serverURL + "/blobs/")
				if err != nil {
					return nil, err
				}

				address := &blobcap.BlobAddress{
					URL:     capabilities.CborURL(*url),
					Headers: map[string]string{},
					Expires: capabilities.CborTime(time.Now().Add(24 * time.Hour)),
				}

				ok := &blobcap.AllocateOK{
					Size:    args.Blob.Size,
					Address: address,
				}

				return bindexec.NewResponse(bindexec.WithSuccess(ok))
			},
		),
	}
}
