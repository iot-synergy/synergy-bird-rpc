package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ LabelModel = (*customLabelModel)(nil)

type (
	// LabelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLabelModel.
	LabelModel interface {
		labelModel
		FindRecord(ctx context.Context, name, typee, parentId string) (*Label, error)
		FindListByParamAndPage(ctx context.Context, typee, parentId string, page, pageSize uint64, recordState int32,
		) (*[]Label, int64, error)
		FindListByIds(ctx context.Context, ids []string) (*[]Label, error)
	}

	customLabelModel struct {
		*defaultLabelModel
	}
)

// NewLabelModel returns a model for the mongo.
func NewLabelModel(url, db, collection string) LabelModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customLabelModel{
		defaultLabelModel: newDefaultLabelModel(conn),
	}
}

func (m *customLabelModel) FindRecord(ctx context.Context, name, typee, parentId string) (*Label, error) {
	var data Label

	filterDate := make(map[string]interface{}) //查询条件data
	if name != "" {
		filterDate["name"] = name
	}
	if typee != "" {
		filterDate["type"] = typee
	}
	if parentId != "" {
		filterDate["parentId"] = parentId
	}
	filterDate["recordState"] = 2
	marshal, err := bson.Marshal(filterDate)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	filter := bson.M{} //查询条件
	err = bson.Unmarshal(marshal, filter)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}

	err = m.conn.FindOne(ctx, &data, filter)
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customLabelModel) FindListByParamAndPage(ctx context.Context, typee, parentId string, page, pageSize uint64, recordState int32,
) (*[]Label, int64, error) {
	data := make([]Label, 0)

	filterDate := make(map[string]interface{}) //查询条件data
	if typee != "" {
		filterDate["type"] = typee
	}
	if parentId != "" {
		filterDate["parentId"] = parentId
	}
	if recordState != 0 {
		filterDate["recordState"] = recordState
	} else {
		filterDate["recordState"] = bson.M{"$ne": 4}
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
		return nil, 0, err
	}
	findoptions := new(options.FindOptions)
	findoptions.SetLimit(int64(pageSize))
	findoptions.SetSkip(int64(page-1) * int64(pageSize))

	findoptions.SetSort(bson.D{bson.E{"updateAt", -1}})
	m.conn.Find(ctx, &data, filter, findoptions)
	return &data, count, nil
}

func (m *customLabelModel) FindListByIds(ctx context.Context, ids []string) (*[]Label, error) {
	data := make([]Label, 0)
	filterDate := make(map[string]interface{}) //查询条件data
	filterDate["recordState"] = 2
	filterDate["_id"] = bson.M{"$in": ids}
	marshal, err := bson.Marshal(filterDate)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	filter := bson.M{} //查询条件
	err = bson.Unmarshal(marshal, filter)
	m.conn.Find(ctx, &data, filter)
	return &data, nil
}
