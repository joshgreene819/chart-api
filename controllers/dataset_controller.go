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
	"go.mongodb.org/mongo-driver/bson"
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

func GetDataset(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	var dataset models.Dataset
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	err := datasetCollection.FindOne(ctx, bson.M{"id": objectID}).Decode(&dataset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestResponse{
			Status:   http.StatusInternalServerError,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.RequestResponse{
		Status:   http.StatusOK,
		Message:  "success",
		Response: &fiber.Map{"data": dataset},
	})
}

func GetAllDatasets(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var datasets []models.Dataset
	defer cancel()

	results, err := datasetCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestResponse{
			Status:   http.StatusInternalServerError,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleDataset models.Dataset
		if err = results.Decode(&singleDataset); err != nil {
			c.Status(http.StatusInternalServerError).JSON(responses.RequestResponse{
				Status:   http.StatusInternalServerError,
				Message:  "error",
				Response: &fiber.Map{"data": err.Error()},
			})
		}
		datasets = append(datasets, singleDataset)
	}

	return c.Status(http.StatusOK).JSON(responses.RequestResponse{
		Status:   http.StatusOK,
		Message:  "success",
		Response: &fiber.Map{"data": datasets},
	})
}

func EditDataset(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	var dataset models.Dataset
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)

	// validate the request body
	if err := c.BodyParser(&dataset); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestResponse{
			Status:   http.StatusBadRequest,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	// validator library to test required fields
	if validationErr := validate.Struct(&dataset); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestResponse{
			Status:   http.StatusBadRequest,
			Message:  "error",
			Response: &fiber.Map{"data": validationErr.Error()},
		})
	}

	// validate dataset against its parent templates
	if complianceErr := makeCompliantDataset(c, &dataset); complianceErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestResponse{
			Status:   http.StatusBadRequest,
			Message:  "error",
			Response: &fiber.Map{"data": complianceErr.Error()},
		})
	}

	update := bson.M{
		"title":           dataset.Title,
		"parentTemplates": dataset.ParentTemplates,
		"data":            dataset.Data,
	}
	result, err := datasetCollection.UpdateOne(ctx, bson.M{"id": objectID}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestResponse{
			Status:   http.StatusInternalServerError,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	// get updated dataset details
	var updatedDataset models.Dataset
	if result.MatchedCount == 1 {
		err := datasetTemplateCollection.FindOne(ctx, bson.M{"id": objectID}).Decode(&updatedDataset)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.RequestResponse{
				Status:   http.StatusInternalServerError,
				Message:  "error",
				Response: &fiber.Map{"data": err.Error()},
			})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.RequestResponse{
		Status:   http.StatusOK,
		Message:  "success",
		Response: &fiber.Map{"data": updatedDataset},
	})
}

func DeleteDataset(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	result, err := datasetCollection.DeleteOne(ctx, bson.M{"id": objectID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestResponse{
			Status:   http.StatusInternalServerError,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(responses.RequestResponse{
			Status:   http.StatusNotFound,
			Message:  "error",
			Response: &fiber.Map{"data": "Dataset with specified ID not found."},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.RequestResponse{
		Status:   http.StatusOK,
		Message:  "success",
		Response: &fiber.Map{"data": "Dataset successfully deleted."},
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
