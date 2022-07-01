package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joshgreene819/chart-api/models"
	"github.com/joshgreene819/chart-api/responses"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateDataset(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var dataset models.Dataset
	defer cancel()

	// Validate request body
	if err := c.BodyParser(&dataset); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestResponse{
			Status:   http.StatusBadRequest,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	if validationErr := validate.Struct(&dataset); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestResponse{
			Status:   http.StatusBadRequest,
			Message:  "error",
			Response: &fiber.Map{"data": validationErr.Error()},
		})
	}

	if complianceErr := makeCompliantDataset(c, &dataset); complianceErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestResponse{
			Status:   http.StatusBadRequest,
			Message:  "error",
			Response: &fiber.Map{"data": complianceErr.Error()},
		})
	}

	newDataset := models.Dataset{
		ID:              primitive.NewObjectID(),
		Title:           dataset.Title,
		ParentTemplates: dataset.ParentTemplates,
		Data:            dataset.Data,
	}

	_, err := datasetCollection.InsertOne(ctx, newDataset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestResponse{
			Status:   http.StatusInternalServerError,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	// Success
	return c.Status(http.StatusCreated).JSON(responses.RequestResponse{
		Status:   http.StatusCreated,
		Message:  "success",
		Response: &fiber.Map{"id": newDataset.ID},
	})
}

func makeCompliantDataset(c *fiber.Ctx, d *models.Dataset) error {
	return nil
}
