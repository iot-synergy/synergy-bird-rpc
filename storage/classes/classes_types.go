package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Classes struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	UpdateAt    time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt    time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
	ClassesId   int64     `bson:"classesId" json:"classesId"`
	ClassesName string    `bson:"classesName,omitempty" json:"classesName,omitempty"`
	ChineseName string    `bson:"chineseName" json:"chineseName"`
	EnglishName string    `bson:"englishName" json:"englishName"`
}
