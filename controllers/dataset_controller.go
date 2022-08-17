package controllers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joshgreene819/chart-api/models"
	"github.com/joshgreene819/chart-api/responses"
	"github.com/xeipuuv/gojsonschema"
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
	if len(d.ParentTemplates) == 0 {
		return nil
	}
	pulledParentTemplates, err := getDatasetTemplates(c, d.ParentTemplates)
	if err != nil {
		return err
	}

	var response bytes.Buffer
	documentLoader := gojsonschema.NewGoLoader(d.Data)
	for _, parent := range pulledParentTemplates {
		// json-schema validation
		schemaLoader := gojsonschema.NewGoLoader(parent.Schema)
		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			return err
		}

		if !(result.Valid()) {
			errorEntry := fmt.Sprintf("The dataset does not comply with parent %s (%s)\n", parent.ID.String(), parent.Title)
			response.WriteString(errorEntry)
			for _, desc := range result.Errors() {
				errorEntry = fmt.Sprintf("- %s\n", desc)
				response.WriteString(errorEntry)
			}
		}
		response.WriteString("\n")
	}
	responseStr := response.String()
	if len(responseStr) == 0 {
		return nil
	}
	return errors.New(responseStr)
}
