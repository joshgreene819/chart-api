package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/joshgreene819/chart-api/configs"
	"go.mongodb.org/mongo-driver/mongo"
)

var datasetTemplateCollection *mongo.Collection = configs.GetCollection(configs.DB, "dataset_templates")
var datasetCollection *mongo.Collection = configs.GetCollection(configs.DB, "datasets")
var validate = validator.New()
