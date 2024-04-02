package gallery

import (
	"context"
	"errors"
	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/gallery"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
	"github.com/iot-synergy/synergy-event-rpc/synergyeventclient"
	"strconv"
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
	aiEvent, err := l.svcCtx.EventRpc.QueryAiEventByTraceId(l.ctx, &synergyeventclient.StringBase{Id: in.TraceId})
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	if aiEvent == nil {
		return nil, errors.New("没有查询到关联的ai事件")
	}
	if in.UserId != aiEvent.GetOwnerId() {
		return nil, errors.New("ai事件不属于当前用户")
	}
	names := strings.Split(aiEvent.GetName(), ",")
	confidences := strings.Split(aiEvent.GetConfidence(), ",")
	name := ""
	confidence := 0.0
	for i := 0; i < len(names) && i < len(confidences); i++ {
		float, e := strconv.ParseFloat(confidences[i], 64)
		if e == nil && float > confidence && names[i] != "" {
			name = names[i]
		}
	}
	if name == "" {
		return nil, errors.New("ai事件当中没有关联鸟类信息")
	}
	illustration, err := l.svcCtx.IllustrationModel.FindOneByTitle(l.ctx, name)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	if illustration == nil {
		return nil, errors.New("ai事件当中关联鸟类信息没有保存到图鉴里")
	}

	gallery := model.Gallery{
		UpdateAt:       time.Time{},
		CreateAt:       time.Time{},
		Name:           illustration.Title,
		UserId:         in.UserId,
		IllustrationId: illustration.ID.Hex(),
		TraceId:        in.TraceId,
		RecordState:    2,
	}
	err = l.svcCtx.GalleryModel.Insert(l.ctx, &gallery)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	// todo:用户解锁图鉴
	return &bird.GalleryResp{
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
	}, nil
}
