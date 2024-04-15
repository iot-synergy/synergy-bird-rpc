package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		FindOneByTitle(ctx context.Context, title string) (*Illustration, error)
		FindOneByEnglishName(ctx context.Context, englishName string) (*Illustration, error)
		FindListByIds(ctx context.Context, ids *[]string) (*[]Illustration, error)
		FindPageJoinGallery(ctx context.Context, labels []string, foreinId, typee, keyword string, isUnlock *bool,
			state int32, page, pageSize uint64) (*[]Illustration, int64, error)
		StatisticLock(ctx context.Context, userId string) (int32, int32, error)
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

func (m *customIllustrationModel) FindOneByTitle(ctx context.Context, title string) (*Illustration, error) {
	var data Illustration
	err := m.conn.FindOne(ctx, &data, bson.M{"title": title, "recordState": 2})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customIllustrationModel) FindOneByEnglishName(ctx context.Context, englishName string) (*Illustration, error) {
	var data Illustration
	err := m.conn.FindOne(ctx, &data, bson.M{"englishName": englishName, "recordState": 2})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customIllustrationModel) FindListByIds(ctx context.Context, ids *[]string) (*[]Illustration, error) {
	oids := make([]primitive.ObjectID, 0)
	for _, id := range *ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			oids = append(oids, oid)
		}
	}
	data := make([]Illustration, 0)
	err := m.conn.Find(ctx, &data, bson.M{"_id": bson.M{"$in": oids}})
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (m *customIllustrationModel) FindPageJoinGallery(ctx context.Context, labels []string, foreinId, typee, keyword string, isUnlock *bool,
	state int32, page, pageSize uint64) (*[]Illustration, int64, error) {
	data := make([]Illustration, 0)
	lookup := bson.M{"from": "gallery_count", "localField": "title", "foreignField": "name", "as": "galleryCount"}
	filterDate := make(map[string]interface{}) //查询条件data
	if keyword != "" && (isUnlock == nil || *isUnlock == true) {
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
	}
	if isUnlock != nil {
		if *isUnlock == true {
			filterDate["galleryCount.userId"] = foreinId
			filterDate["galleryCount._id"] = bson.M{"$exists": true}
			filterDate["galleryCount.recordState"] = 2
			filterDate["galleryCount.count"] = bson.M{"$gte": 1}
		} else if keyword == "" {
			var filterGallery [4]map[string]interface{}
			filterGallery[0] = bson.M{"galleryCount._id": bson.M{"$exists": false}}
			filterGallery[1] = bson.M{"galleryCount.userId": bson.M{"$ne": foreinId}}
			filterGallery[2] = bson.M{"galleryCount.recordState": bson.M{"$ne": 2}}
			filterGallery[3] = bson.M{"galleryCount.count": bson.M{"$lt": 1}}
			filterDate["$or"] = filterGallery
		} else {
			var filterKeywordGallery [2]map[string]interface{}
			var filterKeyword [2]map[string]interface{}
			filterKeyTitle := make(map[string]string)
			filterKeyTitle["$regex"] = keyword
			filterKeyDesc := make(map[string]string)
			filterKeyDesc["$regex"] = keyword
			filterKeyword[0] = make(map[string]interface{})
			filterKeyword[0]["title"] = filterKeyTitle
			filterKeyword[1] = make(map[string]interface{})
			filterKeyword[1]["description"] = filterKeyDesc
			var filterGallery [4]map[string]interface{}
			filterGallery[0] = bson.M{"galleryCount._id": bson.M{"$exists": false}}
			filterGallery[1] = bson.M{"galleryCount.userId": bson.M{"$ne": foreinId}}
			filterGallery[2] = bson.M{"galleryCount.recordState": bson.M{"$ne": 2}}
			filterGallery[3] = bson.M{"galleryCount.count": bson.M{"$lt": 1}}
			filterKeywordGallery[0] = make(map[string]interface{})
			filterKeywordGallery[0]["$or"] = filterKeyword
			filterKeywordGallery[1] = make(map[string]interface{})
			filterKeywordGallery[1]["$or"] = filterGallery
			filterDate["$and"] = filterKeywordGallery
		}
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
	} else {
		filterDate["recordState"] = bson.M{"$ne": 4}
	}
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

	project := bson.M{
		"_id":         1,
		"updateAt":    1,
		"createAt":    1,
		"title":       1,
		"score":       1,
		"wikiUrl":     1,
		"imagePath":   1,
		"iconPath":    1,
		"moreImages":  1,
		"type":        1,
		"labels":      1,
		"description": 1,
		"recordState": 1,
	}

	filter := bson.A{bson.M{"$lookup": lookup}, bson.M{"$match": match}, bson.M{"$project": project}, bson.M{"$limit": pageSize}, bson.M{"$skip": (page - 1) * pageSize}}
	countFilter := bson.A{bson.M{"$lookup": lookup}, bson.M{"$match": match}, bson.M{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}}}

	err = m.conn.Aggregate(ctx, &data, filter)
	if err != nil {
		logx.Error(err.Error())
		return nil, 0, err
	}
	countDO := make([]CountDO, 0)
	err = m.conn.Aggregate(ctx, &countDO, countFilter)
	if err != nil {
		logx.Error(err.Error())
		return nil, 0, err
	}
	var count int64
	for _, do := range countDO {
		count += do.Count
	}
	return &data, count, err
}

func (m *customIllustrationModel) StatisticLock(ctx context.Context, userId string) (unlock int32, lock int32, err error) {
	lockCountDO := make([]CountDO, 0)
	lockFilter := bson.A{
		bson.M{"$lookup": bson.M{"from": "gallery_count", "localField": "title", "foreignField": "name", "as": "galleryCount"}},
		bson.M{"$match": bson.M{
			"galleryCount.userId": bson.M{"$ne": userId},
			"recordState":         2,
		}},
		bson.M{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}},
	}
	err = m.conn.Aggregate(ctx, &lockCountDO, lockFilter)
	if err != nil {
		return
	}
	for _, do := range lockCountDO {
		lock += int32(do.Count)
	}
	unlockCountDO := make([]CountDO, 0)
	unlockFilter := bson.A{
		bson.M{"$lookup": bson.M{"from": "gallery_count", "localField": "title", "foreignField": "name", "as": "galleryCount"}},
		bson.M{"$match": bson.M{
			"galleryCount.userId": userId,
			"recordState":         2,
		}},
		bson.M{"$group": bson.M{"_id": "", "count": bson.M{"$sum": 1}}},
	}
	err = m.conn.Aggregate(ctx, &unlockCountDO, unlockFilter)
	if err != nil {
		return
	}
	for _, do := range unlockCountDO {
		unlock += int32(do.Count)
	}
	return
}
