package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ GalleryModel = (*customGalleryModel)(nil)

type (
	// GalleryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGalleryModel.
	GalleryModel interface {
		galleryModel
		FindOneByNameAndUserId(ctx context.Context, name, userId string) (*Gallery, error)
		FindListByParamAndPage(ctx context.Context, userId string, labels []string, favorite int32,
			page, pageSize uint64) (*[]Gallery, int64, error)
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

func (m *customGalleryModel) FindListByParamAndPage(ctx context.Context, userId string, labels []string, favorite int32,
	page, pageSize uint64) (*[]Gallery, int64, error) {
	data := make([]Gallery, 0)

	filterDate := make(map[string]interface{}) //查询条件data
	if labels != nil && len(labels) != 0 {
		filterLabels := make(map[string][]string)
		filterLabels["$in"] = labels
		filterDate["labels"] = filterLabels
	}
	if userId != "" {
		filterDate["userId"] = userId
	}
	if favorite != 0 {
		filterDate["favorite"] = favorite
	}
	marshal, err := bson.Marshal(filterDate)
	if err != nil {
		logx.Error(err.Error())
		return nil, 0, err
	}
	filter := bson.M{} //查询条件
	err = bson.Unmarshal(marshal, filter)
	if err != nil {
		logx.Error(err.Error())
		return nil, 0, err
	}
	count, err := m.conn.CountDocuments(ctx, filter)
	if err != nil {
		logx.Error(err.Error())
		return nil, 0, err
	}
	findoptions := new(options.FindOptions)
	findoptions.SetLimit(int64(pageSize))
	findoptions.SetSkip(int64(page-1) * int64(pageSize))
	findoptions.SetSort(bson.D{bson.E{"updateAt", -1}})

	err = m.conn.Find(ctx, &data, filter, findoptions)
	if err != nil {
		return nil, 0, err
	}
	return &data, count, nil
}
