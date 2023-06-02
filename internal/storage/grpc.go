package storageserver

import (
	"context"
	"fmt"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/spaghettifunk/norman/internal/common/manager"
	"github.com/spaghettifunk/norman/pkg/consul"
	storageproto "github.com/spaghettifunk/norman/proto/v1/storage"
)

type grpcServer struct {
	storageID           string
	consul              *consul.Consul
	ingestionJobManager *manager.IngestionJobManager
}

func newgrpcServer(id string, cs *consul.Consul) (*grpcServer, error) {
	// Create the new job m
	m, err := manager.NewIngestionJobManager(cs)
	if err != nil {
		return nil, err
	}
	// initialize job manager
	m.Initialize()

	srv := &grpcServer{
		storageID:           id,
		consul:              cs,
		ingestionJobManager: m,
	}
	return srv, nil
}

// InterceptorLogger adapts zerolog logger to interceptor logger.
func InterceptorLogger(l zerolog.Logger) grpc_logging.Logger {
	return grpc_logging.LoggerFunc(func(ctx context.Context, lvl grpc_logging.Level, msg string, fields ...any) {
		l = l.With().Fields(fields).Logger()

		switch lvl {
		case grpc_logging.LevelDebug:
			l.Debug().Msg(msg)
		case grpc_logging.LevelInfo:
			l.Info().Msg(msg)
		case grpc_logging.LevelWarn:
			l.Warn().Msg(msg)
		case grpc_logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func NewGRPCServer(id string, cs *consul.Consul, grpcOpts ...grpc.ServerOption) (*grpc.Server, error) {
	logger := zerolog.New(os.Stderr)

	opts := []grpc_logging.Option{
		grpc_logging.WithLogOnEvents(grpc_logging.StartCall, grpc_logging.FinishCall),
	}

	grpcOpts = append(grpcOpts,
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_logging.StreamServerInterceptor(InterceptorLogger(logger), opts...),
			)), grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
		)),
	)

	gsrv := grpc.NewServer(grpcOpts...)
	srv, err := newgrpcServer(id, cs)
	if err != nil {
		return nil, err
	}
	storageproto.RegisterStorageServer(gsrv, srv)
	return gsrv, nil
}

func (s *grpcServer) Ping(ctx context.Context, req *storageproto.PingRequest) (*storageproto.PingResponse, error) {
	return &storageproto.PingResponse{Msg: "PONG"}, nil
}

func (s *grpcServer) CreateIngestionJob(ctx context.Context, req *storageproto.CreateIngestionJobRequest) (*storageproto.CreateIngestionJobResponse, error) {
	cfg, err := s.consul.GetIngestionJobConfiguration(req.JobID)
	if err != nil {
		return nil, err
	}

	if err := s.ingestionJobManager.Execute(cfg); err != nil {
		return nil, err
	}

	return &storageproto.CreateIngestionJobResponse{StorageID: s.storageID, Message: "Ingestion Job created successfully"}, nil
}

func (s *grpcServer) QueryTable(ctx context.Context, req *storageproto.QueryTableRequest) (*storageproto.QueryTableResponse, error) {
	return nil, nil
}
