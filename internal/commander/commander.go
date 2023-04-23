package commander

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/internal/commander/api"
	configuration "github.com/spaghettifunk/norman/internal/common"
)

type Commander struct {
	Name   string
	config *configuration.Configuration
	server *http.Server
}

func New(config *configuration.Configuration) *Commander {
	addr := fmt.Sprintf("%s:%d", config.Commander.Address, config.Commander.Port)
	return &Commander{
		Name:   "commander",
		config: config,
		server: initServer(addr),
	}
}

func initServer(address string) *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/", api.Version)

	return &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func (c *Commander) StartServer() error {
	log.Info().Msg("Commander Server is ready to handle requests")
	if err := c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (c *Commander) ShutdownServer() error {
	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c.server.SetKeepAlivesEnabled(false)
	return c.server.Shutdown(ctx)
}
