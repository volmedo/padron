package blob

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alanshaw/libracha/digestutil"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/storacha/piri/pkg/store/allocationstore"
	"github.com/storacha/piri/pkg/store/blobstore"
)

func NewBlobGetHandler(blobs blobstore.Blobstore) echo.HandlerFunc {
	if fsblobs, ok := blobs.(blobstore.FileSystemer); ok {
		serveHTTP := http.FileServer(fsblobs.FileSystem()).ServeHTTP
		return func(ctx echo.Context) error {
			r := ctx.Request()
			w := ctx.Response()
			r.URL.Path = r.URL.Path[len("/blob"):]
			serveHTTP(w, r)
			return nil
		}
	}

	log.Error("blobstore does not support filesystem access")
	return func(ctx echo.Context) error {
		return echo.ErrMethodNotAllowed
	}
}

func NewBlobPutHandler(allocs allocationstore.AllocationStore, blobs blobstore.Blobstore) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		digest := ctx.Param("blob")
		mh, err := digestutil.Parse(digest)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("decoding multibase encoded digest: %w", err))
		}

		results, err := allocs.List(ctx.Request().Context(), mh)
		if err != nil {
			return fmt.Errorf("listing allocations: %w", err)
		}

		if len(results) == 0 {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("missing allocation for write to: %s", digestutil.Format(mh)))
		}

		expired := true
		for _, a := range results {
			exp := a.Expires
			if exp > uint64(time.Now().Unix()) {
				expired = false
				break
			}
		}

		if expired {
			return echo.NewHTTPError(http.StatusForbidden, "expired allocation")
		}

		log.Infof("Found %d allocations for write to: %s", len(results), digestutil.Format(mh))

		contentLength, err := strconv.ParseInt(ctx.Request().Header.Get("Content-Length"), 10, 64)
		if err != nil {
			return fmt.Errorf("parsing Content-Length header: %w", err)
		}

		err = blobs.Put(ctx.Request().Context(), mh, uint64(contentLength), ctx.Request().Body)
		if err != nil {
			log.Errorf("writing to %s: %w", digestutil.Format(mh), err)
			if errors.Is(err, blobstore.ErrDataInconsistent) {
				return echo.NewHTTPError(http.StatusConflict, "data consistency check failed")
			}

			return fmt.Errorf("write failed: %w", err)
		}

		ctx.Response().WriteHeader(http.StatusOK)
		return nil
	}
}
