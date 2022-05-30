package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DatasetTemplate struct {
	ID             uint                   `json:"id" gorm:"primary_key"`
	Title          string                 `json:"title"`
	AssignDefaults bool                   `json:"assignDefaults"`
	RequiredKeys   map[string]interface{} `json:"requiredKeys"`
}

type DatasetTemplateInput struct {
	Title          string                 `json:"title"`
	AssignDefaults bool                   `json:"assignDefaults"`
	RequiredKeys   map[string]interface{} `json:"requiredKeys"`
}

// POST /dataset_template
func CreateDatasetTemplate(ctx *gin.Context) {
	// Validate input
	var input DatasetTemplateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := DatasetTemplate{
		Title:          input.Title,
		AssignDefaults: input.AssignDefaults,
		RequiredKeys:   input.RequiredKeys}
	// Add template to MongoDB
	// ...

	ctx.JSON(http.StatusOK, gin.H{"data": template})
}

// GET /dataset_template/:id
func GetDatasetTemplate(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "GET /dataset_template/{id} not yet implemented"})
}

// GET /dataset_template
func GetMatchingDatasetTemplates(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "GET /dataset_template not yet implemented"})
}

// GET /dataset_template
func GetAllDatasetTemplates(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "GET /dataset_template not yet implemented"})
}

// PATCH /dataset_template/:id
func UpdateDatasetTemplate(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "PATCH /dataset_template/:id not yet implemented"})
}

// DELETE /dataset_template/:id
func DeleteDatasetTemplate(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "DELETE /dataset_template/:id not yet implemented"})
}

// DELETE /dataset_template
func DeleteMatchingDatasetTemplate(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "DELETE /dataset_template not yet implemented"})
}

// DELETE /dataset_template
func DeleteAllDatasetTemplates(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "DELETE /dataset_template not yet implemented"})
}
