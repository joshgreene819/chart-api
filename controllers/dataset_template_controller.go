package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joshgreene819/chart-api/configs"
	"github.com/joshgreene819/chart-api/models"
	"github.com/joshgreene819/chart-api/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var datasetTemplateCollection *mongo.Collection = configs.GetCollection(configs.DB, "dataset_templates")
var validate = validator.New()

func CreateDatasetTemplate(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var datasetTemplate models.DatasetTemplate
	defer cancel()

	// Validate request body
	if err := c.BodyParser(&datasetTemplate); err != nil {
		responseObj := responses.DatasetTemplateResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}}
		return c.Status(http.StatusBadRequest).JSON(responseObj)
	}

	if validationErr := validate.Struct(&datasetTemplate); validationErr != nil {
		responseObj := responses.DatasetTemplateResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}}
		return c.Status(http.StatusBadRequest).JSON(responseObj)
	}

	newDatasetTemplate := models.DatasetTemplate{
		ID:             primitive.NewObjectID(),
		Title:          datasetTemplate.Title,
		AssignDefaults: datasetTemplate.AssignDefaults,
		RequiredKeys:   datasetTemplate.RequiredKeys,
	}

	result, err := datasetTemplateCollection.InsertOne(ctx, newDatasetTemplate)
	if err != nil {
		responseObj := responses.DatasetTemplateResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}}
		return c.Status(http.StatusInternalServerError).JSON(responseObj)
	}

	// Success
	return c.Status(http.StatusCreated).JSON(responses.DatasetTemplateResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetDatasetTemplate(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	var datasetTemplate models.DatasetTemplate
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	err := datasetTemplateCollection.FindOne(ctx, bson.M{"id": objectID}).Decode(&datasetTemplate)
	if err != nil {
		responseObj := responses.DatasetTemplateResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}}
		return c.Status(http.StatusInternalServerError).JSON(responseObj)
	}

	return c.Status(http.StatusOK).JSON(responses.DatasetTemplateResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": datasetTemplate}})
}

func EditDatasetTemplate(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	var datasetTemplate models.DatasetTemplate
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)

	// validate the request body
	if err := c.BodyParser(&datasetTemplate); err != nil {
		responseObj := responses.DatasetTemplateResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}}
		return c.Status(http.StatusBadRequest).JSON(responseObj)
	}

	// validator library to test required fields
	if validationErr := validate.Struct(&datasetTemplate); validationErr != nil {
		responseObj := responses.DatasetTemplateResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}}
		return c.Status(http.StatusBadRequest).JSON(responseObj)
	}

	update := bson.M{
		"title":          datasetTemplate.Title,
		"assignDefaults": datasetTemplate.AssignDefaults,
		"requiredKeys":   datasetTemplate.RequiredKeys,
	}
	result, err := datasetTemplateCollection.UpdateOne(ctx, bson.M{"id": objectID}, bson.M{"$set": update})
	if err != nil {
		responseObj := responses.DatasetTemplateResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}}
		return c.Status(http.StatusInternalServerError).JSON(responseObj)
	}

	// Get updated DatasetTemplate details
	var updatedDatasetTemplate models.DatasetTemplate
	if result.MatchedCount == 1 {
		err := datasetTemplateCollection.FindOne(ctx, bson.M{"id": objectID}).Decode(&updatedDatasetTemplate)
		if err != nil {
			responseObj := responses.DatasetTemplateResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}}
			return c.Status(http.StatusInternalServerError).JSON(responseObj)
		}
	}

	return c.Status(http.StatusOK).JSON(responses.DatasetTemplateResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedDatasetTemplate}})
}

func DeleteDatasetTemplate(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)

	result, err := datasetTemplateCollection.DeleteOne(ctx, bson.M{"id": objectID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DatasetTemplateResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(responses.DatasetTemplateResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    &fiber.Map{"data": "Dataset Template with specified ID not found."},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.DatasetTemplateResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": "Dataset Template successfully deleted."},
	})
}