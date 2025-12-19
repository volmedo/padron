package memory

import (
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/sync"
	"go.uber.org/fx"

	"github.com/storacha/piri/pkg/store/acceptancestore"
	"github.com/storacha/piri/pkg/store/allocationstore"
	"github.com/storacha/piri/pkg/store/blobstore"
)

var Module = fx.Module("memory-store",
	fx.Provide(
		NewAllocationStore,
		NewAcceptanceStore,
		NewBlobStore,
	),
)

func NewAllocationStore() (allocationstore.AllocationStore, error) {
	ds := sync.MutexWrap(datastore.NewMapDatastore())
	return allocationstore.NewDsAllocationStore(ds)
}

func NewAcceptanceStore() (acceptancestore.AcceptanceStore, error) {
	ds := sync.MutexWrap(datastore.NewMapDatastore())
	return acceptancestore.NewDsAcceptanceStore(ds)
}

func NewBlobStore() blobstore.Blobstore {
	ds := sync.MutexWrap(datastore.NewMapDatastore())
	return blobstore.NewDsBlobstore(ds)
}
