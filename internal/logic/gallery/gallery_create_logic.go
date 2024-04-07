package gallery

import (
	"context"
	"fmt"
	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/gallery"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
	"github.com/iot-synergy/synergy-event-rpc/synergyeventclient"
	"google.golang.org/grpc/metadata"
	"regexp"
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

	// 查询ai事件
	aiEvent, err := l.svcCtx.EventRpc.QueryAiEventByTraceId(l.ctx, &synergyeventclient.StringBase{Id: in.TraceId})
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	if aiEvent == nil {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "没有查询到关联的ai事件",
			Data: nil,
		}, nil
	}
	fmt.Println("ai事件：")
	fmt.Println(aiEvent)
	fmt.Println(aiEvent.GetOwnerId())
	// 判断事件是否是用户的
	if forein_id != regexp.MustCompile("^peckperk-").ReplaceAllLiteralString(aiEvent.GetOwnerId(), "") {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "ai事件不属于当前用户",
			Data: nil,
		}, nil
	}
	//判断ai事件是否已经创建的成就
	gall, err := l.svcCtx.GalleryModel.FindOneByTraceId(l.ctx, in.TraceId)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	if gall != nil {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "ai事件已经创建过成就了",
			Data: nil,
		}, nil
	}

	// 获取相似度最高的鸟
	names := strings.Split(aiEvent.GetName(), ",")
	confidences := strings.Split(aiEvent.GetConfidence(), ",")
	name := ""        // 鸟的名字
	confidence := 0.0 // 相似度
	for i := 0; i < len(names) && i < len(confidences); i++ {
		float, e := strconv.ParseFloat(confidences[i], 64)
		if e == nil && float > confidence && names[i] != "" { //将相似度最高的鸟的名字赋予name
			confidence = float
			name = names[i]
		}
	}
	if name == "" {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "ai事件当中没有关联鸟类信息",
			Data: nil,
		}, nil
	}
	// 根据鸟的名字查询图鉴
	illustration, err := l.svcCtx.IllustrationModel.FindOneByTitle(l.ctx, name)
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

	//保存
	gallery := model.Gallery{
		UpdateAt:       time.Time{},
		CreateAt:       time.Time{},
		Name:           illustration.Title,
		UserId:         forein_id,
		IllustrationId: illustration.ID.Hex(),
		TraceId:        in.TraceId,
		ImageUrl:       aiEvent.CoverImageUrl,
		RecordState:    2,
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
