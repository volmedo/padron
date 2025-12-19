package config

import (
	"os"
	"path/filepath"

	"github.com/volmedo/padron/pkg/config/app"
)

type StoreConfig struct {
	DataDir string `mapstructure:"data_dir" validate:"required" flag:"data-dir" toml:"data_dir"`
	TempDir string `mapstructure:"temp_dir" validate:"required" flag:"temp-dir" toml:"temp_dir"`
}

func (r StoreConfig) Validate() error {
	return validateConfig(r)
}

func (r StoreConfig) ToAppConfig() (app.StoreConfig, error) {
	if r.DataDir == "" {
		// Return empty config for memory stores
		return app.StoreConfig{}, nil
	}

	if err := os.MkdirAll(r.DataDir, 0755); err != nil {
		return app.StoreConfig{}, err
	}
	if err := os.MkdirAll(r.TempDir, 0755); err != nil {
		return app.StoreConfig{}, err
	}

	out := app.StoreConfig{
		DataDir: r.DataDir,
		TempDir: r.TempDir,
		Blobs: app.BlobStoreConfig{
			Dir:    filepath.Join(r.DataDir, "blobs"),
			TmpDir: filepath.Join(r.TempDir, "storage"),
		},
		Allocations: app.AllocationStoreConfig{
			Dir: filepath.Join(r.DataDir, "allocation"),
		},
		Acceptance: app.AcceptanceStoreConfig{
			Dir: filepath.Join(r.DataDir, "acceptance"),
		},
	}

	return out, nil
}
