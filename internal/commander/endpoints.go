package commander

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spaghettifunk/norman/internal/common/entities"
	cingestion "github.com/spaghettifunk/norman/internal/common/ingestion"
	storage_v1 "github.com/spaghettifunk/norman/proto/v1/storage"
)

const (
	apiName    = "Norman - Commander APIs"
	apiVersion = "v0.0.1"
)

func (c *Commander) APIVersion(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"API Name":    apiName,
		"API Version": apiVersion,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
	})
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
	Job cingestion.IngestionJobConfiguration `json:"job"`
}

func (c *Commander) GetJobs(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) GetJob(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) CreateJob(ctx *fiber.Ctx) error {
	payload := &CreateIngestionJobRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	// parse config and transform into an IngestionJob
	if err := cingestion.NewIngestionJob(&payload.Job); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create ingestion job",
			"error":   err.Error(),
		})
	}

	if err := c.consul.PutIngestionJobConfiguration(&payload.Job); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save ingestion job into Consul",
			"error":   err.Error(),
		})
	}

	// call gRPC function to trigger the ingestion job
	req := &storage_v1.CreateIngestionJobRequest{JobID: payload.Job.ID.String()}
	resp, err := c.storageGRPCClient.CreateIngestionJob(context.Background(), req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create ingestion job",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"storageID": resp.StorageID,
		"message":   resp.Message,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
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
Table routes
*/
type CreateTableRequest struct {
	Table entities.Table `json:"table"`
}

func (c *Commander) GetTables(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) GetTable(ctx *fiber.Ctx) error {
	return nil
}

func (c *Commander) CreateTable(ctx *fiber.Ctx) error {
	payload := &CreateTableRequest{}
	if err := ctx.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.consul.PutTableConfiguration(&payload.Table); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Table",
			"error":   err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message":   "Table created successfully",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
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
