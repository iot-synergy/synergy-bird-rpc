package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gallery struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	UpdateAt       time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt       time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
	Name           string    `bson:"name,omitempty" json:"name,omitempty"`
	UserId         string    `bson:"userId,omitempty" json:"userId,omitempty"`
	IllustrationId string    `bson:"illustrationId,omitempty" json:"illustrationId,omitempty"`
	TraceId        string    `bson:"traceId,omitempty" json:"traceId,omitempty"`
	ImageUrl       string    `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	//1:Created
	//2:Normal
	//3:Deleted
	//4:Forbidden
	RecordState int8 `bson:"recordState,omitempty" json:"recordState,omitempty"`
}
