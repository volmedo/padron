package filesystem

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/storacha/piri/pkg/store/acceptancestore"
	"github.com/storacha/piri/pkg/store/allocationstore"
	"github.com/storacha/piri/pkg/store/blobstore"
	"go.uber.org/fx"

	"github.com/volmedo/padron/pkg/config/app"
)

var Module = fx.Module("filesystem-store",
	fx.Provide(
		ProvideConfigs,
		NewAllocationStore,
		NewAcceptanceStore,
		NewBlobStore,
	),
)

type Configs struct {
	fx.Out
	Blob       app.BlobStoreConfig
	Allocation app.AllocationStoreConfig
	Acceptance app.AcceptanceStoreConfig
}

// ProvideConfigs provides the fields of a storage config
func ProvideConfigs(cfg app.StoreConfig) Configs {
	return Configs{
		Allocation: cfg.Allocations,
		Blob:       cfg.Blobs,
		Acceptance: cfg.Acceptance,
	}
}

func NewBlobStore(cfg app.BlobStoreConfig) (blobstore.Blobstore, error) {
	if cfg.Dir == "" {
		return nil, fmt.Errorf("no data dir provided for blob store")
	}
	var tmpDir = cfg.TmpDir
	if tmpDir == "" {
		tmpDir = filepath.Join(os.TempDir(), "storage")
	}

	bs, err := blobstore.NewFsBlobstore(cfg.Dir, tmpDir)
	if err != nil {
		return nil, fmt.Errorf("creating blob store: %w", err)
	}
	return bs, nil
}

func NewAllocationStore(cfg app.AllocationStoreConfig, lc fx.Lifecycle) (allocationstore.AllocationStore, error) {
	if cfg.Dir == "" {
		return nil, fmt.Errorf("no data dir provided for allocation store")
	}

	ds, err := newDs(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("creating allocation store: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ds.Close()
		},
	})

	return allocationstore.NewDsAllocationStore(ds)
}

func NewAcceptanceStore(cfg app.AcceptanceStoreConfig, lc fx.Lifecycle) (acceptancestore.AcceptanceStore, error) {
	if cfg.Dir == "" {
		return nil, fmt.Errorf("no data dir provided for acceptance store")
	}

	ds, err := newDs(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("creating allocation store: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ds.Close()
		},
	})

	return acceptancestore.NewDsAcceptanceStore(ds)
}

func newDs(path string) (*leveldb.Datastore, error) {
	dirPath, err := mkdirp(path)
	if err != nil {
		return nil, fmt.Errorf("creating leveldb for store at path %s: %w", path, err)
	}
	return leveldb.NewDatastore(dirPath, nil)
}

func mkdirp(dirpath ...string) (string, error) {
	dir := filepath.Join(dirpath...)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("creating directory: %s: %w", dir, err)
	}
	return dir, nil
}
