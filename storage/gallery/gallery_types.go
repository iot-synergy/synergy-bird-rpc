package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gallery struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
	Name     string    `bson:"name,omitempty" json:"name,omitempty"`
	UserId   string    `bson:"userId,omitempty" json:"userId,omitempty"`
	Favorite int32     `bson:"favorite,omitempty" json:"favorite,omitempty"`
	Labels   []string  `bson:"labels,omitempty" json:"labels,omitempty"`
	//1:Created
	//2:Normal
	//3:Deleted
	//4:Forbidden
	RecordState int8 `bson:"recordState,omitempty" json:"recordState,omitempty"`
}
