package storageserver

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spaghettifunk/norman/pkg/consul"
)

type StorageServer struct {
	Name   string
	ID     uuid.UUID
	consul *consul.Consul
	config configuration.Configuration
	// gRPC server to receive requests from Commander
}

func New(config configuration.Configuration) (*StorageServer, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// initialize consul client
	cs := consul.New()
	if err := cs.Init(); err != nil {
		return nil, err
	}

	return &StorageServer{
		Name:   "storage",
		ID:     id,
		config: config,
		consul: cs,
	}, nil
}

func (st *StorageServer) StartServer(address string) error {
	// register to consul
	log.Info().Msg("register and declare Storage to Consul")
	if err := st.consul.Start(st); err != nil {
		return err
	}
	if err := st.consul.Declare(st); err != nil {
		return err
	}
	return nil
}

func (st *StorageServer) ShutdownServer() error {
	// deregister to consul
	log.Info().Msg("deregister Storage to Consul")
	if err := st.consul.Stop(st); err != nil {
		return err
	}

	log.Info().Msg("Shutting down server...")
	return nil
}

func (st *StorageServer) GetHost() string {
	hn, err := os.Hostname()
	if err != nil {
		panic(err.Error())
	}
	return hn
}

func (st *StorageServer) GetPort() string {
	return fmt.Sprint(st.config.Storage.Port)
}

func (st *StorageServer) GetName() string {
	return st.Name
}

func (st *StorageServer) GetID() string {
	return st.ID.String()
}

func (st *StorageServer) GetMetadata() map[string]string {
	return nil
}
