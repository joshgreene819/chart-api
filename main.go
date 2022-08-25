package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joshgreene819/chart-api/configs"
	"github.com/joshgreene819/chart-api/controllers"
)

func main() {
	app := fiber.New()

	// Setup database client connection
	configs.ClientConnectDB()

	// Sample route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello world! (from Fiber & MongoDB)"})
	})

	// DatasetTemplate routes
	app.Post("/datasetTemplate", controllers.CreateDatasetTemplate)
	app.Get("/datasetTemplate", controllers.GetAllDatasetTemplates)
	app.Get("/datasetTemplate/:id", controllers.GetDatasetTemplate)
	app.Put("/datasetTemplate/:id", controllers.EditDatasetTemplate)
	app.Delete("/datasetTemplate/:id", controllers.DeleteDatasetTemplate)

	// Dataset routes
	app.Post("/dataset", controllers.CreateDataset)
	app.Get("/dataset", controllers.GetAllDatasets)
	app.Get("/dataset/:id", controllers.GetDataset)
	app.Put("/dataset/:id", controllers.EditDataset)
	app.Delete("/dataset/:id", controllers.DeleteDatasetTemplate)

	app.Listen(":8080")
}
