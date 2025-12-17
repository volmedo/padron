package blob

import (
	"net/url"
	"time"

	"github.com/alanshaw/ucantone/execution/bindexec"
	logging "github.com/ipfs/go-log/v2"

	blobcap "github.com/volmedo/padron/pkg/capabilities/blob"
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
					URL:     blobcap.CborURL(*url),
					Headers: map[string]string{},
					Expires: blobcap.CborTime(time.Now().Add(24 * time.Hour)),
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
