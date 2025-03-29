package storageserver

import (
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spaghettifunk/norman/internal/common/utils"
	"google.golang.org/grpc"
)

type StorageServer struct {
	server *grpc.Server
	config configuration.Configuration

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

func New(config configuration.Configuration) (*StorageServer, error) {
	return &StorageServer{
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
	log.Info().Msg("start gRPC server")
	if err := st.setupServer("<id-here>"); err != nil {
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
	log.Info().Msg("Shutting down server...")
	return nil
}
