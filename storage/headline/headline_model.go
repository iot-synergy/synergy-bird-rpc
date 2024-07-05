package model

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ HeadlineModel = (*customHeadlineModel)(nil)

type (
	// HeadlineModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHeadlineModel.
	HeadlineModel interface {
		headlineModel
		InsertOne(ctx context.Context, data *Headline) (*Headline, error)
		FindListByParam(ctx context.Context, site string, page, pageSize uint64) (*[]Headline, int64, error)
		FindListByIndex(ctx context.Context, lastIndex int64) (*[]Headline, error)
	}

	customHeadlineModel struct {
		*defaultHeadlineModel
	}
)

// NewHeadlineModel returns a model for the mongo.
func NewHeadlineModel(url, db, collection string) HeadlineModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customHeadlineModel{
		defaultHeadlineModel: newDefaultHeadlineModel(conn),
	}
}

func (m *defaultHeadlineModel) InsertOne(ctx context.Context, data *Headline) (*Headline, error) {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	_, err := m.conn.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	var result Headline

	err = m.conn.FindOne(ctx, &result, bson.M{"_id": data.ID})
	switch err {
	case nil:
		return &result, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customHeadlineModel) FindListByParam(ctx context.Context, site string, page, pageSize uint64) (*[]Headline, int64, error) {
	data := make([]Headline, 0)
	filterDate := make(map[string]interface{}) //查询条件data
	if site != "" {
		filterDate["$or"] = bson.A{
			bson.M{"site": bson.M{"$regex": site, "$options": "i"}},
		}
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
	findoptions.SetSort(bson.D{bson.E{Key: "updateAt", Value: -1}})

	err = m.conn.Find(ctx, &data, filter, findoptions)
	if err != nil {
		return nil, 0, err
	}
	return &data, count, nil
}

func (m *customHeadlineModel) FindListByIndex(ctx context.Context, lastIndex int64) (*[]Headline, error) {
	data := make([]Headline, 0)
	// filterDate := make(map[string]interface{}) //查询条件data
	// if lastIndex != "" {
	// 	id, err := strconv.ParseInt(lastIndex, 10, 64)
	// 	filterDate["$or"] = bson.A{
	// 		bson.M{"site": bson.M{"$regex": site, "$options": "i"}},
	// 	}
	// }

	findoptions := new(options.FindOptions)
	findoptions.SetLimit(20)
	findoptions.SetSkip(lastIndex)
	findoptions.SetSort(bson.D{bson.E{Key: "updateAt", Value: -1}})

	err := m.conn.Find(ctx, &data, bson.M{}, findoptions)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
