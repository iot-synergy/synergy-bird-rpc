package model

import (
	"context"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/illustration"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

var _ GalleryModel = (*customGalleryModel)(nil)

type (
	// GalleryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGalleryModel.
	GalleryModel interface {
		galleryModel
		FindOneByNameAndUserId(ctx context.Context, name, userId string) (*Gallery, error)
		FindListByParamAndPage(ctx context.Context, userId, illustrationId, name string, startTime, endTime int64,
			page, pageSize uint64, labelIds *[]string) (*[]GalleryJoinIllustrationDO, int64, error)
		FindOneByTraceId(ctx context.Context, traceId string) (*Gallery, error)
		CountByUserIdAndIllustrationId(ctx context.Context, userId, illustrationId string) (int64, error)
	}

	customGalleryModel struct {
		*defaultGalleryModel
	}
)

// NewGalleryModel returns a model for the mongo.
func NewGalleryModel(url, db, collection string) GalleryModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customGalleryModel{
		defaultGalleryModel: newDefaultGalleryModel(conn),
	}
}
func (m *customGalleryModel) FindOneByNameAndUserId(ctx context.Context, name, userId string) (*Gallery, error) {
	var data Gallery

	err := m.conn.FindOne(ctx, &data, bson.M{
		"name":   name,
		"userId": userId})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customGalleryModel) FindListByParamAndPage(ctx context.Context, userId string, illustrationId, name string,
	startTime, endTime int64, page, pageSize uint64, labelIds *[]string) (*[]GalleryJoinIllustrationDO, int64, error) {
	data := make([]GalleryJoinIllustrationDO, 0)
	var lookup bson.M
	if labelIds != nil && len(*labelIds) > 0 {
		lookup = bson.M{"from": "illustration", "localField": "name", "foreignField": "title", "as": "illustration",
			"pipeline": bson.A{bson.M{"$match": bson.M{"labels": bson.M{"$in": *labelIds}, "recordState": 2}}}}
	} else {
		lookup = bson.M{"from": "illustration", "localField": "name", "foreignField": "title", "as": "illustration",
			"pipeline": bson.A{bson.M{"$match": bson.M{"recordState": 2}}}}
	}
	filterDate := make(map[string]interface{}) //查询条件data
	filterDate["userId"] = userId
	if illustrationId != "" {
		filterDate["illustrationId"] = illustrationId
	}
	if name != "" {
		filterDate["name"] = bson.M{"$regex": name}
	}
	if startTime != 0 {
		filterDate["updateAt"] = bson.M{"$gte": time.UnixMilli(startTime)}
	}
	if endTime != 0 {
		filterDate["updateAt"] = bson.M{"$lt": time.UnixMilli(endTime)}
	}
	filterDate["recordState"] = 2
	filterDate["illustration"] = bson.M{"$gt": make([]string, 0)}
	marshal, err := bson.Marshal(filterDate)
	if err != nil {
		logx.Error(err.Error())
		return nil, 0, err
	}
	match := bson.M{} //查询条件
	err = bson.Unmarshal(marshal, match)
	if err != nil {
		logx.Error(err.Error())
		return nil, 0, err
	}

	filter := bson.A{bson.M{"$lookup": lookup}, bson.M{"$match": match}, bson.M{"$limit": pageSize}, bson.M{"$skip": (page - 1) * pageSize}}
	countFilter := bson.A{bson.M{"$lookup": lookup}, bson.M{"$match": match}, bson.M{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}}}
	err = m.conn.Aggregate(ctx, &data, filter)
	if err != nil {
		return nil, 0, err
	}
	countDO := make([]model.CountDO, 0)
	err = m.conn.Aggregate(ctx, &countDO, countFilter)
	if err != nil {
		return nil, 0, err
	}
	var count int64
	for _, do := range countDO {
		count += do.Count
	}
	return &data, count, nil
}

func (m *customGalleryModel) FindOneByTraceId(ctx context.Context, traceId string) (*Gallery, error) {
	var data Gallery

	err := m.conn.FindOne(ctx, &data, bson.M{
		"traceId":     traceId,
		"recordState": 2})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customGalleryModel) CountByUserIdAndIllustrationId(ctx context.Context, userId, illustrationId string) (int64, error) {
	count, err := m.conn.CountDocuments(ctx, bson.M{"userId": userId, "illustrationId": illustrationId, "recordState": 2})
	if err != nil {
		return 0, err
	}
	return count, nil
}
