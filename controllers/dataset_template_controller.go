package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joshgreene819/chart-api/models"
	"github.com/joshgreene819/chart-api/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateDatasetTemplate(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var datasetTemplate models.DatasetTemplate
	defer cancel()

	// Validate request body
	if err := c.BodyParser(&datasetTemplate); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestResponse{
			Status:   http.StatusBadRequest,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	if validationErr := validate.Struct(&datasetTemplate); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestResponse{
			Status:   http.StatusBadRequest,
			Message:  "error",
			Response: &fiber.Map{"data": validationErr.Error()},
		})
	}

	newDatasetTemplate := models.DatasetTemplate{
		ID:           primitive.NewObjectID(),
		Title:        datasetTemplate.Title,
		Options:      datasetTemplate.Options,
		RequiredKeys: datasetTemplate.RequiredKeys,
	}

	_, err := datasetTemplateCollection.InsertOne(ctx, newDatasetTemplate)
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
		Response: &fiber.Map{"id": newDatasetTemplate.ID},
	})
}

func GetDatasetTemplate(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	var datasetTemplate models.DatasetTemplate
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	err := datasetTemplateCollection.FindOne(ctx, bson.M{"id": objectID}).Decode(&datasetTemplate)
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
		Response: &fiber.Map{"data": datasetTemplate},
	})
}

func GetAllDatasetTemplates(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var datasetTemplates []models.DatasetTemplate
	defer cancel()

	results, err := datasetTemplateCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestResponse{
			Status:   http.StatusInternalServerError,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleDatasetTemplate models.DatasetTemplate
		if err = results.Decode(&singleDatasetTemplate); err != nil {
			c.Status(http.StatusInternalServerError).JSON(responses.RequestResponse{
				Status:   http.StatusInternalServerError,
				Message:  "error",
				Response: &fiber.Map{"data": err.Error()},
			})
		}
		datasetTemplates = append(datasetTemplates, singleDatasetTemplate)
	}

	return c.Status(http.StatusOK).JSON(responses.RequestResponse{
		Status:   http.StatusOK,
		Message:  "success",
		Response: &fiber.Map{"data": datasetTemplates},
	})
}

func EditDatasetTemplate(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	var datasetTemplate models.DatasetTemplate
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)

	// validate the request body
	if err := c.BodyParser(&datasetTemplate); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestResponse{
			Status:   http.StatusBadRequest,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	// validator library to test required fields
	if validationErr := validate.Struct(&datasetTemplate); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.RequestResponse{
			Status:   http.StatusBadRequest,
			Message:  "error",
			Response: &fiber.Map{"data": validationErr.Error()},
		})
	}

	update := bson.M{
		"title":        datasetTemplate.Title,
		"options":      datasetTemplate.Options,
		"requiredKeys": datasetTemplate.RequiredKeys,
	}
	result, err := datasetTemplateCollection.UpdateOne(ctx, bson.M{"id": objectID}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RequestResponse{
			Status:   http.StatusInternalServerError,
			Message:  "error",
			Response: &fiber.Map{"data": err.Error()},
		})
	}

	// Get updated DatasetTemplate details
	var updatedDatasetTemplate models.DatasetTemplate
	if result.MatchedCount == 1 {
		err := datasetTemplateCollection.FindOne(ctx, bson.M{"id": objectID}).Decode(&updatedDatasetTemplate)
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
		Response: &fiber.Map{"data": updatedDatasetTemplate},
	})
}

func DeleteDatasetTemplate(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Params("id")
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)

	result, err := datasetTemplateCollection.DeleteOne(ctx, bson.M{"id": objectID})
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
			Response: &fiber.Map{"data": "Dataset Template with specified ID not found."},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.RequestResponse{
		Status:   http.StatusOK,
		Message:  "success",
		Response: &fiber.Map{"data": "Dataset Template successfully deleted."},
	})
}

func getDatasetTemplates(c *fiber.Ctx, datasetTemplateIDs []primitive.ObjectID) ([]models.DatasetTemplate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var templates []models.DatasetTemplate
	defer cancel()

	cursor, err := datasetTemplateCollection.Find(ctx, bson.M{"$in": datasetTemplateIDs})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var currentTemplate models.DatasetTemplate
		if err = cursor.Decode(&currentTemplate); err != nil {
			return templates, err
		}
		templates = append(templates, currentTemplate)
	}
	return templates, nil
}
