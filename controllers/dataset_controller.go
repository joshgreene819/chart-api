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

	if complianceErr := validate.Struct(&dataset); complianceErr != nil {
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

func makeCompliantDataset(d *models.Dataset) error {
	/*
		This helper function has to...
			> Verify that every entry in d.ParentTemplates slice is a valid template
				* Fail if any of them aren't valid
				* CAN succeed if len(d.ParentTemplates) == 0
			> For each template in d.ParentTemplates...
				* d must have all of the template's required keys
					- if not AND template has AssignDefaults: true - assign the default (first come first
						serve basis, no overwriting with another templates defaults)
					- if not AND template has AssignDefaults: false - keep going under the assumption that
						a template not yet processed with have the key and AssignDefaults: true. If the end
						has been reached and the required key has not been assigned through defaults, fail
						with error
	*/
	return nil
}
