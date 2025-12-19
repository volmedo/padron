package store

import (
	"github.com/volmedo/padron/pkg/fx/store/filesystem"
	"github.com/volmedo/padron/pkg/fx/store/memory"
)

// FileSystemStoreModule provides filesystem-backed stores
var FileSystemStoreModule = filesystem.Module

// MemoryStoreModule provides memory-backed stores
var MemoryStoreModule = memory.Module
