package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IllustrationJoinGalleryCountDO struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UpdateAt     time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt     time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
	Title        string             `bson:"title,omitempty" json:"tittle,omitempty"`
	Score        float64            `bson:"score,omitempty" json:"score,omitempty"`
	WikiUrl      string             `bson:"wikiUrl,omitempty" json:"wikiUrl,omitempty"`
	ImagePath    string             `bson:"imagePath,omitempty" json:"imagePath,omitempty"`
	IconPath     string             `bson:"iconPath,omitempty" json:"iconPath,omitempty"`
	MoreImages   []string           `bson:"moreImages,omitempty" json:"moreImages,omitempty"`
	Type         string             `bson:"type,omitempty" json:"type,omitempty"`
	Labels       []string           `bson:"labels,omitempty" json:"labels,omitempty"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	RecordState  int8               `bson:"recordState,omitempty" json:"recordState,omitempty"`
	GalleryCount []GalleryCountDO   `bson:"galleryCount,omitempty" json:"galleryCount,omitempty"`
}

type GalleryCountDO struct {
	UnlockTime time.Time `bson:"unlockTime,omitempty" json:"unlockTime,omitempty"`
	UserId     string    `bson:"userId,omitempty" json:"userId,omitempty"`
}

type CountDO struct {
	ID    string `bson:"_id,omitempty" json:"_id,omitempty"`
	Count int64  `bson:"count,omitempty" json:"count,omitempty"`
}
