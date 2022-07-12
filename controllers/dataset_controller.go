package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
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
	if len(d.ParentTemplates) == 0 {
		return nil
	}
	pulledParentTemplates, err := getDatasetTemplates(c, d.ParentTemplates)
	if err != nil {
		return err
	}

	for _, parent := range pulledParentTemplates {
		err := complyWithParent(d, parent)
		if err != nil {
			return err
		}
	}
	return nil
}

func complyWithParent(d *models.Dataset, parent models.DatasetTemplate) error {
	for key, parentValue := range parent.OneTimeData.Data {
		if !(parent.OneTimeData.Metadata.AnyDepth) {
			switch parentValueType := parentValue.(type) {
			case nil:
				if _, ok := d.Data[key]; !ok {
					if parent.OneTimeData.Metadata.AssignDefaults {
						d.Data[key] = parentValue
					} else {
						msg := fmt.Sprintf("dataset does not contain key '%s' and parent template %s (%s) does not assign default values", key, parent.Title, parent.ID.String())
						return errors.New(msg)
					}
				}
				// TEMPORARY Just to suppress error
				fmt.Println(parentValueType)

			case map[string]interface{}:
				// object
				continue

			case []interface{}:
				// list
				continue

			default:
				// primitive type
				datasetValue, ok := d.Data[key]
				if !ok {
					if parent.OneTimeData.Metadata.AssignDefaults {
						d.Data[key] = parentValue
					} else {
						msg := fmt.Sprintf("dataset does not contain key '%s' and parent template %s (%s) does not assign default values", key, parent.Title, parent.ID.String())
						return errors.New(msg)
					}
				} else {
					if reflect.TypeOf(datasetValue) != reflect.TypeOf(parentValue) {
						if parent.OneTimeData.Metadata.AssignDefaults {
							d.Data[key] = parentValue
						} else {
							msg := fmt.Sprintf("key '%s' in dataset but has a different type than the data in parent template '%s' (%s)\ndataset: %v\tparent template: %v", key, parent.Title, parent.ID.String(), reflect.TypeOf(datasetValue).String(), reflect.TypeOf(parentValue).String())
							return errors.New(msg)
						}
					}
				}
			}
		} else {
			// nested stuff
			// ...
		}
		if !(parent.IteratedData.Metadata.AnyDepth) {
			// direct checks for each entry in the list
			// ...
		} else {
			// nested stuff
			// ...
		}
	}
	return nil
}
