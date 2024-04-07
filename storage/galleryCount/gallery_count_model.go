package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
)

var _ GalleryCountModel = (*customGalleryCountModel)(nil)

type (
	// GalleryCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGalleryCountModel.
	GalleryCountModel interface {
		galleryCountModel
		FindOneByUserIdAndIllustrationId(ctx context.Context, userId, illustrationId string) (*GalleryCount, error)
	}

	customGalleryCountModel struct {
		*defaultGalleryCountModel
	}
)

// NewGalleryCountModel returns a model for the mongo.
func NewGalleryCountModel(url, db, collection string) GalleryCountModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customGalleryCountModel{
		defaultGalleryCountModel: newDefaultGalleryCountModel(conn),
	}
}

func (m *customGalleryCountModel) FindOneByUserIdAndIllustrationId(ctx context.Context, userId, illustrationId string) (*GalleryCount, error) {
	var data GalleryCount

	err := m.conn.FindOne(ctx, &data, bson.M{"userId": userId, "illustrationId": illustrationId})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}
