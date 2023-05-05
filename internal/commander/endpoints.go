package commander

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	apiName    = "Norman - Commander APIs"
	apiVersion = "v0.0.1"
)

func (c *Commander) APIVersion(ctx *fiber.Ctx) error {
	ctx.Status(200).JSON(fiber.Map{
		"API Name":    apiName,
		"API Version": apiVersion,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
	})
	return nil
}

/*
Tenant routes
*/
func (c *Commander) GetTenants(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) GetTenant(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) CreateTenant(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) UpdateTenant(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) PatchTenant(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) DeleteTenant(ctx *fiber.Ctx) error {
	return nil
}

/*
Table routes
*/
func (c *Commander) GetTables(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) GetTable(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) CreateTable(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) UpdateTable(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) PatchTable(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) DeleteTable(ctx *fiber.Ctx) error {
	return nil
}

/*
Segment routes
*/
func (c *Commander) GetSegments(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) GetSegment(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) CreateSegment(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) UpdateSegment(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) PatchSegment(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) DeleteSegment(ctx *fiber.Ctx) error {
	return nil
}

/*
Schema routes
*/
func (c *Commander) GetSchemas(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) GetSchema(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) CreateSchema(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) UpdateSchema(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) PatchSchema(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) DeleteSchema(ctx *fiber.Ctx) error {
	return nil
}
