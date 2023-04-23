package storageserver

import (
	"fmt"

	configuration "github.com/spaghettifunk/norman/internal/common"
)

type StorageServer struct {
	Name   string
	config configuration.Configuration
	// gRPC server to receive requests from Commander
}

func New(config configuration.Configuration) *StorageServer {
	_ = fmt.Sprintf("%s:%d", config.Storage.Address, config.Storage.Port)
	return &StorageServer{
		Name:   "storage",
		config: config,
	}
}
