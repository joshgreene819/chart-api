package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshgreene819/chart-api/configs"
	"github.com/joshgreene819/chart-api/resources"
	"go.mongodb.org/mongo-driver/mongo"
)

var appConfiguration *configs.Configuration = configs.LoadConfiguration()
var dbClient *mongo.Client = configs.ConnectDB(appConfiguration.MongoURI)

func main() {
	router := gin.Default()
	// Temporary
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	// Dataset Template
	router.POST("/dataset_template", resources.CreateDatasetTemplate)
	router.GET("/dataset_template/:id", resources.GetDatasetTemplate)
	router.GET("/dataset_template", resources.GetMatchingDatasetTemplates)
	// router.GET("/dataset_template", resources.GetAllDatasetTemplates)
	router.PATCH("/dataset_template/:id", resources.UpdateDatasetTemplate)
	router.DELETE("/dataset_template/:id", resources.DeleteDatasetTemplate)
	router.DELETE("/dataset_template", resources.DeleteMatchingDatasetTemplate)
	// router.DELETE("/dataset_template", resources.DeleteAllDatasetTemplates)

	// Dataset
	// ...

	// Chart
	// ...

	router.Run()
}
