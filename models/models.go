package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DatasetTemplate struct {
	ID             primitive.ObjectID     `json:"id,omitempty"`
	Title          string                 `json:"title,omitempty" validate:"required"`
	AssignDefaults bool                   `json:"assignDefaults,omitempty" validate:"required"`
	RequiredKeys   map[string]interface{} `json:"requiredKeys,omitempty" validate:"required"`
}

type Dataset struct {
	ID              primitive.ObjectID     `json:"id,omitempty"`
	Title           string                 `json:"title,omitempty" validate:"required"`
	ParentTemplates []primitive.ObjectID   `json:"parentTemplates,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty" validate:"required"`
}