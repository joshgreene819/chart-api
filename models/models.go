package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// type BehaviorInList struct {
// 	InEachEntry bool `json:"inEachEntry"`
// }

// type Options struct {
// 	AssignDefaults bool           `json:"assignDefaults"`
// 	AnyDepth       bool           `json:"anyDepth"`
// 	BehaviorInList BehaviorInList `json:"behaviorInList"`
// }

// type DatasetTemplate struct {
// 	ID           primitive.ObjectID     `json:"id,omitempty"`
// 	Title        string                 `json:"title,omitempty" validate:"required"`
// 	Options      Options                `json:"options,omitempty" validate:"required"`
// 	RequiredKeys map[string]interface{} `json:"requiredKeys,omitempty" validate:"required"`
// }

// Provides configuration for how a dataset and its keys can be compliant with
// a DatasetTemplate
type Metadata struct {
	AnyDepth       bool `json:"anyDepth"`
	AssignDefaults bool `json:"assignDefaults"`
	StrictMatch    bool `json:"strictMatch"`
}

// Keys that a DatasetTemplate will require in order to be implemented by a
// dataset
type RequiredKeys struct {
	Metadata Metadata               `json:"metadata"`
	Data     map[string]interface{} `json:"data"`
}

// Definition of keys that a Dataset will need in order to comply with the
// DatasetTemplate. Complying with a DatasetTemplate guarantees that queried
// Datasets will have a certain set of keys and characteristics shared by other
// Datasets that comply
type DatasetTemplate struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	Title        string             `json:"title,omitempty" validate:"required"`
	OneTimeData  RequiredKeys       `json:"oneTimeData,omitempty" validate:"required"`
	IteratedData RequiredKeys       `json:"iteratedData,omitempty" validate:"required"`
}

// Arbitrary set of data used to build customizable Charts
type Dataset struct {
	ID              primitive.ObjectID     `json:"id,omitempty"`
	Title           string                 `json:"title,omitempty" validate:"required"`
	ParentTemplates []primitive.ObjectID   `json:"parentTemplates,omitempty" validate:"required"`
	Data            map[string]interface{} `json:"data,omitempty" validate:"required"`
}
