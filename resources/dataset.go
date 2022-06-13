package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Dataset struct {
	ID             uint                   `json:"id" gorm:"primary_key" bson:"id"`
	Title          string                 `json:"title" bson:"title"`
	ParentTemplate int                    `json:"parentTemplate" bson:"parentTemplate,omitempty"`
	Data           map[string]interface{} `json:"data" bson:"data"`
}

// User-submitted dataset info (without ID)
type DatasetInput struct {
	Title          string                 `json:"title"`
	ParentTemplate int                    `json:"parentTemplate"`
	Data           map[string]interface{} `json:"data"`
}

func CreateDataset(ctx *gin.Context) {
	var input Dataset
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dataset := Dataset{Title: input.Title, ParentTemplate: input.ParentTemplate, Data: input.Data}
	// TODO - database transaction
	// ...
	ctx.JSON(http.StatusOK, gin.H{"data": dataset})
}

// Get by ID
func GetDataset(ctx *gin.Context) {
	return
}

// Get by criteria
func GetMatchingDatasets(ctx *gin.Context) {
	return
}

// Get all
func GetAllDatasets(ctx *gin.Context) {
	return
}

func UpdateDataset(ctx *gin.Context) {
	return
}

// Delete by ID
func DeleteDataset(ctx *gin.Context) {
	return
}

// Delete by criteria
func DeleteMatchingDatasets(ctx *gin.Context) {
	return
}

// Delete all
func DeleteAllDatasets(ctx *gin.Context) {
	return
}
