package model

import (
	"context"
	"errors"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ ClassesModel = (*customClassesModel)(nil)

type (
	// ClassesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customClassesModel.
	ClassesModel interface {
		classesModel
		FindOneByClassesIdOrClassesName(ctx context.Context, classesId int64, classesName string) (*Classes, error)
		FindOneByClassesId(ctx context.Context, classesId int64) (*Classes, error)
		UpdateOneChineseNameByClassesId(ctx context.Context, classesId int64, chineseName string) (*mongo.UpdateResult, error)
		UpdateOneEnglishNameByClassesId(ctx context.Context, classesId int64, englishName string) (*mongo.UpdateResult, error)
		BatchInsert(ctx context.Context, data *[]Classes) error
		BatchDelete(ctx context.Context) error
		FindListByParam(ctx context.Context, keyword string, page, pageSize uint64) (*[]Classes, int64, error)
	}

	customClassesModel struct {
		*defaultClassesModel
	}
)

// NewClassesModel returns a model for the mongo.
func NewClassesModel(url, db, collection string) ClassesModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customClassesModel{
		defaultClassesModel: newDefaultClassesModel(conn),
	}
}

func (m *defaultClassesModel) FindOneByClassesIdOrClassesName(ctx context.Context, classesId int64, classesName string) (*Classes, error) {
	var data Classes
	err := m.conn.FindOne(ctx, &data, bson.M{"$or": bson.A{bson.M{"classesId": classesId}, bson.M{"classesName": classesName}}})
	switch {
	case errors.Is(err, mon.ErrNotFound):
		return nil, nil
	case err == nil:
		return &data, nil
	default:
		return nil, err
	}
}

func (m *defaultClassesModel) FindOneByClassesId(ctx context.Context, classesId int64) (*Classes, error) {
	var data Classes
	err := m.conn.FindOne(ctx, &data, bson.M{"classesId": classesId})
	switch {
	case errors.Is(err, mon.ErrNotFound):
		return nil, nil
	case err == nil:
		return &data, nil
	default:
		return nil, err
	}
}

func (m *defaultClassesModel) UpdateOneChineseNameByClassesId(ctx context.Context, classesId int64, chineseName string) (*mongo.UpdateResult, error) {
	res, err := m.conn.UpdateMany(ctx, bson.M{"classesId": classesId}, bson.M{"$set": bson.M{"chineseName": chineseName}})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *defaultClassesModel) UpdateOneEnglishNameByClassesId(ctx context.Context, classesId int64, englishName string) (*mongo.UpdateResult, error) {
	res, err := m.conn.UpdateMany(ctx, bson.M{"classesId": classesId}, bson.M{"$set": bson.M{"englishName": englishName}})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *defaultClassesModel) BatchInsert(ctx context.Context, data *[]Classes) error {
	array := bson.A{}
	for _, classes := range *data {
		array = append(array, classes)
	}
	_, err := m.conn.InsertMany(ctx, array)
	return err
}

func (m *defaultClassesModel) BatchDelete(ctx context.Context) error {
	_, err := m.conn.DeleteMany(ctx, bson.M{})
	return err
}

func (m *defaultClassesModel) FindListByParam(ctx context.Context, keyword string, page, pageSize uint64) (*[]Classes, int64, error) {
	data := make([]Classes, 0)
	filterDate := make(map[string]interface{}) //查询条件data
	if keyword != "" {
		id, err := strconv.ParseInt(keyword, 10, 64)
		if err == nil {
			filterDate["classesId"] = id
		} else {
			filterDate["$or"] = bson.A{
				bson.M{"classesName": bson.M{"$regex": keyword, "$options": "i"}},
				bson.M{"chineseName": bson.M{"$regex": keyword, "$options": "i"}},
				bson.M{"englishName": bson.M{"$regex": keyword, "$options": "i"}},
			}
		}
	}

	marshal, err := bson.Marshal(filterDate)
	if err != nil {
		logx.Error(err.Error())
		return nil, 0, err
	}
	filter := bson.M{} //查询条件
	err = bson.Unmarshal(marshal, filter)

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
