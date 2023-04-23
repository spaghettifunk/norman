package broker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	configuration "github.com/spaghettifunk/norman/internal/common"
)

type Broker struct {
	Name   string
	config configuration.Configuration
	server *http.Server
}

func New(config configuration.Configuration) *Broker {
	addr := fmt.Sprintf("%s:%d", config.Broker.Address, config.Broker.Port)
	return &Broker{
		Name:   "storage",
		config: config,
		server: initServer(addr),
	}
}

func initServer(address string) *http.Server {
	router := http.NewServeMux()

	// subscribe routes
	// router.Handle("/", middleware.WithLogging(api.Version()))

	return &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func (b *Broker) StartServer() error {
	log.Info().Msg("Storage Server is ready to handle requests")
	if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (b *Broker) ShutdownServer() error {
	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	b.server.SetKeepAlivesEnabled(false)
	return b.server.Shutdown(ctx)
}
