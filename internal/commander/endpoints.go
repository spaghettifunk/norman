package commander

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spaghettifunk/norman/internal/common/model"
)

const (
	apiName    = "Norman - Commander APIs"
	apiVersion = "v0.0.1"
)

func (c *Commander) APIVersion(ctx *fiber.Ctx) error {
	ctx.Status(http.StatusOK).JSON(fiber.Map{
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
Ingestion Job routes
*/
type CreateIngestionJobRequest struct {
	model.IngestionJobConfiguration
}

func (c *Commander) GetJobs(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) GetJob(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) CreateJob(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) UpdateJob(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) PatchJob(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) DeleteJob(ctx *fiber.Ctx) error {
	return nil
}

/*
Schema routes
*/
type CreateSchemaRequest struct {
	model.Schema
}

func (c *Commander) GetSchemas(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) GetSchema(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) CreateSchema(ctx *fiber.Ctx) error {
	// Validate the body payload -- a bit useless for now
	payload := &CreateSchemaRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	// TODO: change this to a better Request struct
	// we pass the body directly for now
	if err := c.schemaManager.Execute(ctx.Body()); err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Schema",
			"error":   err.Error(),
		})
		return err
	}
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message":   "Schema created successfully",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})

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
