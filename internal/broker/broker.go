package broker

import (
	"fmt"
	"os"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spaghettifunk/norman/pkg/consul"
)

type Broker struct {
	Name     string
	ID       uuid.UUID
	Hostname string
	consul   *consul.Consul
	config   configuration.Configuration
	app      *fiber.App
}

func New(config configuration.Configuration) (*Broker, error) {
	// Create new Fiber application
	app := fiber.New(fiber.Config{
		AppName:           "broker-api-server",
		EnablePrintRoutes: true, // TODO: change this based on logger level -- DEBUG
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
	})
	// add default middleware
	app.Use(recover.New())

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// initialize consul client
	cs := consul.New()
	if err := cs.Init(); err != nil {
		return nil, err
	}

	// get the hostname from the machine
	hn, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	br := &Broker{
		Name:     "broker",
		ID:       id,
		Hostname: hn,
		consul:   cs,
		config:   config,
		app:      app,
	}

	br.setupRoutes()

	return br, nil
}

func (b *Broker) setupRoutes() {
	apiV1 := b.app.Group("/broker/v1")

	// query routes
	queryEndpoints := apiV1.Group("/query")
	queryEndpoints.Post("/", b.CreateQuery)
}

func (b *Broker) StartServer(address string) error {
	// register to consul
	log.Info().Msg("register and declare Commander to Consul")
	if err := b.consul.Start(b); err != nil {
		return err
	}
	if err := b.consul.Declare(b); err != nil {
		return err
	}

	log.Info().Msg("Storage Server is ready to handle requests")
	return b.app.Listen(address)
}

func (b *Broker) ShutdownServer() error {
	// deregister to consul
	log.Info().Msg("deregister Commander to Consul")
	if err := b.consul.Stop(b); err != nil {
		return err
	}

	log.Info().Msg("Shutting down server...")
	return b.app.Shutdown()
}

func (b *Broker) GetHost() string {
	hn, err := os.Hostname()
	if err != nil {
		panic(err.Error())
	}
	return hn
}

func (b *Broker) GetPort() string {
	return fmt.Sprint(b.config.Broker.Port)
}

func (b *Broker) GetName() string {
	return b.Name
}

func (b *Broker) GetID() string {
	return b.ID.String()
}

func (b *Broker) GetMetadata() map[string]string { return nil }
