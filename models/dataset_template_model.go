package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DatasetTemplate struct {
	ID             primitive.ObjectID     `json:"id,omitempty"`
	Title          string                 `json:"name,omitempty" validate:"required"`
	AssignDefaults bool                   `json:"assignDefaults,omitempty" validate:"required"`
	RequiredKeys   map[string]interface{} `json:"requiredKeys,omitempty" validate:"required"`
}
