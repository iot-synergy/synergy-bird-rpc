package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ IllustrationModel = (*customIllustrationModel)(nil)

type (
	// IllustrationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customIllustrationModel.
	IllustrationModel interface {
		illustrationModel
		FindListByParamAndPage(ctx context.Context, labels []string, typee, keyword string,
			state int32, page, pageSize uint64) (*[]Illustration, int64, error)
	}

	customIllustrationModel struct {
		*defaultIllustrationModel
	}
)

// NewIllustrationModel returns a model for the mongo.
func NewIllustrationModel(url, db, collection string) IllustrationModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customIllustrationModel{
		defaultIllustrationModel: newDefaultIllustrationModel(conn),
	}
}

func (m *customIllustrationModel) FindListByParamAndPage(ctx context.Context, labels []string,
	typee, keyword string, state int32, page, pageSize uint64) (*[]Illustration, int64, error) {
	data := make([]Illustration, 0)

	filterDate := make(map[string]interface{}) //查询条件data
	if keyword != "" {
		var filterKeyword [2]map[string]interface{}
		filterKeyTitle := make(map[string]string)
		filterKeyTitle["$regex"] = keyword
		filterKeyDesc := make(map[string]string)
		filterKeyDesc["$regex"] = keyword
		filterKeyword[0] = make(map[string]interface{})
		filterKeyword[0]["title"] = filterKeyTitle
		filterKeyword[1] = make(map[string]interface{})
		filterKeyword[1]["description"] = filterKeyDesc
		filterDate["$or"] = filterKeyword

		//filterDate["$or"] = bson.A{
		//	bson.M{"title": bson.M{"regex": primitive.Regex{Pattern: ".*" + keyword + ".*", Options: "i"}}},
		//	bson.M{"description": bson.M{"regex": primitive.Regex{Pattern: ".*" + keyword + ".*", Options: "i"}}},
		//}
	}
	if labels != nil && len(labels) > 0 {
		filterLabels := make(map[string][]string)
		filterLabels["$in"] = labels
		filterDate["labels"] = filterLabels
	}
	if typee != "" {
		filterDate["type"] = typee
	}
	if state != 0 {
		filterDate["recordState"] = state
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
	m.conn.Find(ctx, &data, filter, findoptions)
	return &data, count, nil
}
