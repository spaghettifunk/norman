package broker

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	apiName    = "Norman - Broker APIs"
	apiVersion = "v0.0.1"
)

func (b *Broker) APIVersion(ctx *fiber.Ctx) error {
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"API Name":    apiName,
		"API Version": apiVersion,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
	})
	return nil
}

func (b *Broker) CreateQuery(ctx *fiber.Ctx) error {
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"result":    "nice job!",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
	return nil
}
