package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
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
		FindIllustrationIdList(ctx context.Context, userId string, illustrationIds []string) (*[]string, error)
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

func (m *customGalleryCountModel) FindIllustrationIdList(ctx context.Context, userId string, illustrationIds []string) (*[]string, error) {
	data := make([]GalleryCount, 0)
	filterDate := make(map[string]interface{}) //查询条件data
	filterDate["recordState"] = 2
	filterDate["count"] = bson.M{"$gt": 0}
	filterDate["illustrationId"] = bson.M{"$in": illustrationIds}
	marshal, err := bson.Marshal(filterDate)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	filter := bson.M{} //查询条件
	err = bson.Unmarshal(marshal, filter)
	m.conn.Find(ctx, &data, filter)
	resp := make([]string, 0)
	for _, datum := range data {
		resp = append(resp, datum.IllustrationId)
	}
	return &resp, nil
}
