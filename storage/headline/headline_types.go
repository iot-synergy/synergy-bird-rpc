package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Headline struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Url         string `bson:"url,omitempty"  json:"url"`
	Site        string `bson:"site,omitempty" json:"site"`
	Cover       bool   `bson:"cover,omitempty"  json:"cover"`
	Title       string `bson:"title,omitempty"  json:"title"`
	Description string `bson:"description,omitempty"  json:"description"`
	Image       string `bson:"image,omitempty"  json:"image"`
	State       bool   `bson:"state,omitempty"  json:"state"`

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
