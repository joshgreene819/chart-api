package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshgreene819/chart-api/resources"
)

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

/* Test connecting to MongoDB container */
// func tempMain() {
// 	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb0.example.com:27017"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Disconnect(ctx)

// 	// List databases
// 	databases, err := client.ListDatabaseNames(ctx, bson.M{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(databases)
// }
