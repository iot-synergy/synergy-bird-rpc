package gallery

import (
	"context"
	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/gallery"
	model2 "github.com/iot-synergy/synergy-bird-rpc/storage/galleryCount"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
	"github.com/iot-synergy/synergy-event-rpc/synergyeventclient"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/metadata"
	"regexp"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GalleryCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGalleryCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GalleryCreateLogic {
	return &GalleryCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GalleryCreateLogic) GalleryCreate(in *bird.GalleryCreateReq) (*bird.GalleryResp, error) {
	// 获取用户id
	value := metadata.ValueFromIncomingContext(l.ctx, "gateway-firebaseid")
	if len(value) <= 0 {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "用户未登录",
			Data: nil,
		}, nil
	}
	forein_id := strings.Join(value, "")

	// 根据鸟的名字查询图鉴
	illustration, err := l.svcCtx.IllustrationModel.FindOneByEnglishName(l.ctx, in.Name)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	if illustration == nil {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "ai事件当中关联鸟类信息没有保存到图鉴里",
			Data: nil,
		}, nil
	}

	var onlyTraceId, imageUrl string
	var isCreate bool
label:
	for _, traceId := range in.TraceIds {
		// 查询ai事件
		aiEvent, err := l.svcCtx.EventRpc.QueryAiEventByTraceId(l.ctx, &synergyeventclient.StringBase{Id: traceId})
		if err != nil {
			logx.Error(err.Error())
			continue label
		}
		if aiEvent == nil {
			return &bird.GalleryResp{
				Code: -1,
				Msg:  "没有查询到关联的ai事件:" + traceId,
				Data: nil,
			}, nil
		}
		// 判断事件是否是用户的
		if forein_id != regexp.MustCompile("^peckperk-").ReplaceAllLiteralString(aiEvent.GetOwnerId(), "") {
			return &bird.GalleryResp{
				Code: -1,
				Msg:  "ai事件不属于当前用户:" + traceId,
				Data: nil,
			}, nil
		}
		//判断ai事件是否已经创建的成就
		gall, err := l.svcCtx.GalleryModel.FindOneByTraceId(l.ctx, traceId)
		if err != nil {
			logx.Error(err.Error())
			continue label
		}
		if gall != nil {
			isCreate = true
			continue label
		}
		names := strings.Split(aiEvent.GetName(), ",")
		for _, name := range names {
			if name == in.Name {
				imageUrl = aiEvent.CoverImageUrl
				onlyTraceId = traceId
				break label
			}
		}
	}

	if onlyTraceId == "" {
		if isCreate {
			return &bird.GalleryResp{
				Code: -1,
				Msg:  "事件已经创建过成就",
				Data: nil,
			}, nil
		}
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "没有在事件中查询到<" + in.Name + ">",
			Data: nil,
		}, nil
	}

	//保存
	gallery := model.Gallery{
		UpdateAt:       time.Now(),
		CreateAt:       time.Now(),
		Name:           illustration.Title,
		UserId:         forein_id,
		IllustrationId: illustration.ID.Hex(),
		TraceId:        onlyTraceId,
		ImageUrl:       imageUrl,
		RecordState:    2,
	}
	//更新图鉴数量
	galleryCountData, err := l.svcCtx.GalleryCountModel.FindOneByUserIdAndIllustrationId(l.ctx, forein_id, illustration.ID.Hex())
	if galleryCountData == nil || galleryCountData.Name == "" {
		galleryCountData = &model2.GalleryCount{
			ID:             primitive.ObjectID{},
			UpdateAt:       time.Time{},
			CreateAt:       time.Time{},
			Name:           illustration.Title,
			UserId:         forein_id,
			IllustrationId: illustration.ID.Hex(),
			Count:          1,
			RecordState:    2,
			UnlockTime:     gallery.CreateAt,
		}
		l.svcCtx.GalleryCountModel.Insert(l.ctx, galleryCountData)
	} else {
		count, _ := l.svcCtx.GalleryModel.CountByUserIdAndIllustrationId(l.ctx, forein_id, illustration.ID.Hex())
		galleryCountData.Count = count
		if count == 1 {
			galleryCountData.UnlockTime = gallery.CreateAt
		}
		l.svcCtx.GalleryCountModel.Update(l.ctx, galleryCountData)
	}

	err = l.svcCtx.GalleryModel.Insert(l.ctx, &gallery)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	return &bird.GalleryResp{
		Code: 0,
		Msg:  "",
		Data: &bird.GalleryRespData{
			Id:          gallery.ID.Hex(),
			RecordState: int32(gallery.RecordState),
			CreateTime:  gallery.CreateAt.UnixMilli(),
			Name:        gallery.Name,
			UserId:      gallery.UserId,
			Illustration: &bird.IllustrationsResp{
				Id:          illustration.ID.Hex(),
				RecordState: int32(illustration.RecordState),
				CreateTime:  illustration.CreateAt.UnixMilli(),
				Title:       illustration.Title,
				Score:       illustration.Score,
				WikiUrl:     illustration.WikiUrl,
				ImagePath:   illustration.ImagePath,
				IconPath:    illustration.IconPath,
				MoreImages:  illustration.MoreImages,
				Typee:       illustration.Type,
				Labels:      illustration.Labels,
				Description: illustration.Description,
			},
		},
	}, nil
}
