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

		// check one time data
		if !(parent.OneTimeData.Metadata.AnyDepth) {
			for k, v := range parent.OneTimeData.Data {
				if _, ok := d.Data[k]; !ok {
					// doesn't exist
					if parent.OneTimeData.Metadata.AssignDefaults {
						d.Data[k] = v
					} else {
						return errors.New("TBD")
					}
				}
			}
			// do some direct accessing
			// for each key in parent.OneTimeData.Data ...
			//		d[key] -> should exist, be the same type
		}

		// check iterated data
		if !(parent.IteratedData.Metadata.AnyDepth) {

		}
	}

	return nil
}

func complyWithParent(d *models.Dataset, parent models.DatasetTemplate) error {
	// I hope you turn out recursive
	if !(parent.OneTimeData.Metadata.AnyDepth) {
		for parentKey, parentValue := range parent.OneTimeData.Data {
			switch parentValueType := parentValue.(type) {
			case nil:
				// null type in parent key means that dataset can have whatever kind of data
				continue

			case map[string]interface{}:
				// Make dynamnic strict from v and try and marshal d.Data[k] into v's struct
				// for k, v in parentValue 
				//   DynamicStruct.add(k, v)
				// DynamicStruct.build().new()
				// unmarshal d.Data[parentKey] if it exists
				continue

			case []interface{}:
				continue

			default:
				// direct comparison
				if _, ok := d.Data[parentKey]; !ok {
					if parent.OneTimeData.Metadata.AssignDefaults {
						d.Data[parentKey] = parentValue
					} else {
						msg := fmt.Sprintf("dataset does not contain required key: '%s' in top-level of data", parentKey)
						errors.New(msg)
					}
				} else {
					if reflect.TypeOf(parentValue) != reflect.TypeOf(d.Data[parentKey]) {
						if parent.OneTimeData.Metadata.AssignDefaults {
							d.Data[parentKey] = parentValue
						} else {
							msg := fmt.Sprintf("dataset contains key: '%s' but with an incorrect type. Dataset template (%v) has key of type %T while dataset has key of type %T",
								parentKey, d.ID, parentValue, d.Data[parentKey])
							errors.New(msg)
						}
					}
				}
			}
		}
	}

	return nil
}
