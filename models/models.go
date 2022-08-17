package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Struct containing properties that compliant datasets will have. DatasetTemplate.Schema is simply
// json schema unmarshaled into map[string]interface{}. A Dataset's compliance to a DatasetTemplate
// is determined by validation via a json-schema library
type DatasetTemplate struct {
	ID     primitive.ObjectID     `json:"id,omitempty"`
	Title  string                 `json:"title,omitempty" validate:"required"`
	Schema map[string]interface{} `json:"schema,omitempty" validate:"required"`
}

// Arbitrary set of data used to build customizable Charts. Expectations & restrictions can be set
// by assigning a parent template via its ID to the Dataset instance
type Dataset struct {
	ID              primitive.ObjectID     `json:"id,omitempty"`
	Title           string                 `json:"title,omitempty" validate:"required"`
	ParentTemplates []primitive.ObjectID   `json:"parentTemplates,omitempty" validate:"required"`
	Data            map[string]interface{} `json:"data,omitempty" validate:"required"`
}
