package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Illustration struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	UpdateAt    time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt    time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
	Title       string    `bson:"title,omitempty" json:"tittle,omitempty"`
	Score       float64   `bson:"score,omitempty" json:"score,omitempty"`
	WikiUrl     string    `bson:"wikiUrl,omitempty" json:"wikiUrl,omitempty"`
	ImagePath   string    `bson:"imagePath,omitempty" json:"imagePath,omitempty"`
	MoreImages  []string  `bson:"moreImages,omitempty" json:"moreImages,omitempty"`
	Type        string    `bson:"type,omitempty" json:"type,omitempty"`
	Labels      []string  `bson:"labels,omitempty" json:"labels,omitempty"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	// 1:Created
	// 2:Normal
	// 3:Deleted
	// 4:Forbidden
	RecordState int8 `bson:"recordState,omitempty" json:"recordState,omitempty"`
}
