package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BehaviorInList struct {
	InEachEntry bool `json:"inEachEntry"`
}

type Options struct {
	AssignDefaults bool           `json:"assignDefaults"`
	AnyDepth       bool           `json:"anyDepth"`
	BehaviorInList BehaviorInList `json:"behaviorInList"`
}

type DatasetTemplate struct {
	ID           primitive.ObjectID     `json:"id,omitempty"`
	Title        string                 `json:"title,omitempty" validate:"required"`
	Options      Options                `json:"options,omitempty" validate:"required"`
	RequiredKeys map[string]interface{} `json:"requiredKeys,omitempty" validate:"required"`
}

type Dataset struct {
	ID              primitive.ObjectID     `json:"id,omitempty"`
	Title           string                 `json:"title,omitempty" validate:"required"`
	ParentTemplates []primitive.ObjectID   `json:"parentTemplates,omitempty" validate:"required"`
	Data            map[string]interface{} `json:"data,omitempty" validate:"required"`
}
