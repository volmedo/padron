package blob

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/alanshaw/libracha/capabilities/assert"
	"github.com/alanshaw/libracha/digestutil"
	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/ucan"
	"github.com/alanshaw/ucantone/ucan/delegation"
	"github.com/alanshaw/ucantone/ucan/delegation/policy"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	mh "github.com/multiformats/go-multihash"
	ucantodid "github.com/storacha/go-ucanto/did"
	"github.com/storacha/piri/pkg/store"
	"github.com/storacha/piri/pkg/store/acceptancestore"
	"github.com/storacha/piri/pkg/store/acceptancestore/acceptance"
	"github.com/storacha/piri/pkg/store/allocationstore"
	"github.com/storacha/piri/pkg/store/allocationstore/allocation"
	"github.com/storacha/piri/pkg/store/blobstore"
)

var log = logging.Logger("service/blob")

type Service struct {
	id          ucan.Signer
	publicURL   *url.URL
	blobs       blobstore.Blobstore
	allocations allocationstore.AllocationStore
	acceptances acceptancestore.AcceptanceStore
}

func NewService(
	id ucan.Signer,
	publicURL *url.URL,
	blobs blobstore.Blobstore,
	allocations allocationstore.AllocationStore,
	acceptances acceptancestore.AcceptanceStore,
) *Service {
	return &Service{
		id:          id,
		publicURL:   publicURL,
		blobs:       blobs,
		allocations: allocations,
		acceptances: nil,
	}
}

type Blob struct {
	Digest mh.Multihash
	Size   uint64
}

type Address struct {
	URL     *url.URL
	Headers http.Header
	Expires time.Time
}

func (s *Service) Allocate(ctx context.Context, space did.DID, blob Blob, cause ucan.Link) (uint64, *Address, error) {
	log := log.With("space", space.String(), "blob", digestutil.Format(blob.Digest))
	log.Infof("allocating blob of size %d", blob.Size)

	// check if we already have an allocation for the blob in this space
	allocs, err := s.allocations.List(ctx, blob.Digest)
	if err != nil {
		log.Errorw("getting allocations", "error", err)
		return 0, nil, fmt.Errorf("getting allocations: %w", err)
	}

	allocated := false
	for _, a := range allocs {
		if a.Space.String() == space.String() {
			allocated = true
			break
		}
	}

	received := false
	// check if we received the blob (only possible if we have an allocation)
	if len(allocs) > 0 {
		_, err = s.blobs.Get(ctx, blob.Digest)
		if err != nil && !errors.Is(err, store.ErrNotFound) {
			log.Errorw("getting blob", "error", err)
			return 0, nil, fmt.Errorf("getting blob: %w", err)
		}
		if err == nil {
			received = true
		}
	}

	// the size reported in the receipt is the number of bytes allocated
	// in the space - if a previous allocation already exists, this has
	// already been done, so the allocation size is 0
	size := blob.Size
	if allocated {
		log.Info("blob allocation already exists")
		size = 0
	}

	// nothing to do
	if allocated && received {
		log.Info("blob already received")
		return size, nil, nil
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	var address *Address
	// if not received yet, we need to generate a signed URL for the
	// upload, and include it in the receipt.
	if !received {
		address = &Address{
			URL:     s.getBlobURL(blob.Digest),
			Headers: http.Header{},
			Expires: expiresAt,
		}
	}

	// even if a previous allocation was made in this space, we create
	// another for the new invocation.
	sp, _ := ucantodid.Parse(space.String())
	err = s.allocations.Put(ctx, allocation.Allocation{
		Space:   sp,
		Blob:    allocation.Blob(blob),
		Expires: uint64(expiresAt.Unix()),
		Cause:   cidlink.Link{Cid: cid.Cid(cause)},
	})
	if err != nil {
		log.Errorw("putting allocation", "error", err)
		return 0, nil, fmt.Errorf("putting allocation: %w", err)
	}

	return size, address, nil
}

func (s *Service) Accept(ctx context.Context, space did.DID, blob Blob, cause ucan.Link) (*delegation.Delegation, error) {
	log := log.With("space", space.String(), "blob", digestutil.Format(blob.Digest))
	log.Infof("allocating blob of size %d", blob.Size)

	// check if we already got the blob
	_, err := s.blobs.Get(ctx, blob.Digest)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, fmt.Errorf("blob not found: %w", err)
		}

		log.Errorw("getting blob", "error", err)
		return nil, fmt.Errorf("getting blob: %w", err)
	}

	sp, _ := ucantodid.Parse(space.String())
	acc := acceptance.Acceptance{
		Space:      sp,
		Blob:       acceptance.Blob(blob),
		ExecutedAt: uint64(time.Now().Unix()),
		Cause:      cidlink.Link{Cid: cid.Cid(cause)},
	}

	err = s.acceptances.Put(ctx, acc)
	if err != nil {
		log.Errorw("putting acceptance for blob", "error", err)
		return nil, fmt.Errorf("putting acceptance for blob: %w", err)
	}

	// build location commitment
	locCommitment, err := assert.Location.Delegate(
		s.id,
		space,
		s.id,
		delegation.WithPolicyBuilder(
			policy.Equal(".space", space.String()),
			policy.Equal(".content", digestutil.Format(blob.Digest)),
			policy.Equal(".location[0].url", s.getBlobURL(blob.Digest).String()),
			policy.Equal(".range.offset", 0),
			policy.Equal(".range.length", blob.Size),
		),
		delegation.WithNoExpiration(),
	)
	if err != nil {
		log.Errorw("creating location commitment", "error", err)
		return nil, fmt.Errorf("creating location commitment: %w", err)
	}

	// TODO(vic): store and publish the location commitment

	return locCommitment, nil
}

func (s *Service) getBlobURL(digest mh.Multihash) *url.URL {
	return s.publicURL.JoinPath("blob", digestutil.Format(digest))
}
