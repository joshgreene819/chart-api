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

	// Implemented routes
	app.Post("/datasetTemplate", controllers.CreateDatasetTemplate)
	app.Get("/datasetTemplate/:id", controllers.GetDatasetTemplate)
	app.Put("/datasetTemplate/:id", controllers.EditDatasetTemplate)
	app.Delete("/datasetTemplate/:id", controllers.DeleteDatasetTemplate)

	app.Listen(":8080")
}
