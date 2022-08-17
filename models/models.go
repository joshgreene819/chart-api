package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Definition of keys that a Dataset will need in order to comply with the
// DatasetTemplate. Complying with a DatasetTemplate guarantees that queried
// Datasets will have a certain set of keys and characteristics shared by other
// Datasets that comply
type DatasetTemplate struct {
	ID     primitive.ObjectID     `json:"id,omitempty"`
	Title  string                 `json:"title,omitempty" validate:"required"`
	Schema map[string]interface{} `json:"schema,omitempty" validate:"required"`
}

// Arbitrary set of data used to build customizable Charts
type Dataset struct {
	ID              primitive.ObjectID     `json:"id,omitempty"`
	Title           string                 `json:"title,omitempty" validate:"required"`
	ParentTemplates []primitive.ObjectID   `json:"parentTemplates,omitempty" validate:"required"`
	Data            map[string]interface{} `json:"data,omitempty" validate:"required"`
}
