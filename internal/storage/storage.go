package storageserver

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spaghettifunk/norman/internal/common/utils"
	"google.golang.org/grpc"
)

type StorageServer struct {
	Name     string
	ID       uuid.UUID
	Hostname string
	server   *grpc.Server
	config   configuration.Configuration

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

func New(config configuration.Configuration) (*StorageServer, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// get the hostname from the machine
	hn, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &StorageServer{
		Name:      "storage",
		ID:        id,
		Hostname:  hn,
		config:    config,
		shutdowns: make(chan struct{}),
	}, nil
}

func (st *StorageServer) setupServer(id string) error {
	var err error
	st.server, err = NewGRPCServer(id)
	if err != nil {
		return err
	}

	rpcAddr, err := utils.RPCAddr(st.config.Storage.BindAddr, st.config.Storage.RPCPort)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return err
	}

	go func() {
		if err := st.server.Serve(ln); err != nil {
			_ = st.ShutdownGRPC()
		}
	}()
	return err
}

func (st *StorageServer) StartServer(address string) error {
	// register to consul
	log.Info().Msg("register and declare Storage to Consul")

	log.Info().Msg("start gRPC server")
	if err := st.setupServer(st.ID.String()); err != nil {
		return err
	}
	return nil
}

func (st *StorageServer) ShutdownGRPC() error {
	st.shutdownLock.Lock()
	defer st.shutdownLock.Unlock()

	if st.shutdown {
		return nil
	}

	st.shutdown = true
	close(st.shutdowns)

	shutdown := []func() error{
		func() error {
			st.server.GracefulStop()
			return nil
		},
	}
	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (st *StorageServer) ShutdownServer() error {
	// deregister to consul
	log.Info().Msg("Shutting down server...")
	return nil
}

func (st *StorageServer) GetHost() string {
	return st.Hostname
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
