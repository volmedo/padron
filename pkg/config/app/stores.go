package app

// StoreConfig contains all storage paths and directories
type StoreConfig struct {
	// Root directories
	DataDir string
	TempDir string

	// Service-specific storage subdirectories
	Blobs       BlobStoreConfig
	Allocations AllocationStoreConfig
	Acceptance  AcceptanceStoreConfig
}

// BlobStoreConfig contains blob-specific storage paths
type BlobStoreConfig struct {
	Dir    string
	TmpDir string
}

// AllocationStoreConfig contains allocation-specific storage paths
type AllocationStoreConfig struct {
	Dir string
}

// AcceptanceStoreConfig contains acceptance-specific storage paths
type AcceptanceStoreConfig struct {
	Dir string
}
