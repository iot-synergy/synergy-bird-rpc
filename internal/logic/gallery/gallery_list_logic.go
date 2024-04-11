package gallery

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strings"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type GalleryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGalleryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GalleryListLogic {
	return &GalleryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GalleryListLogic) GalleryList(in *bird.GalleryListReq) (*bird.GalleryListResp, error) {
	// 获取用户id
	value := metadata.ValueFromIncomingContext(l.ctx, "gateway-firebaseid")
	if len(value) <= 0 {
		return &bird.GalleryListResp{
			Code: -1,
			Msg:  "用户未登录",
			Data: nil,
		}, nil
	}
	forein_id := strings.Join(value, "")
	data, count, err := l.svcCtx.GalleryModel.FindListByParamAndPage(l.ctx, forein_id, in.GetIllustrationId(),
		in.GetName(), in.GetStartTime(), in.GetEndTime(), in.Page, in.PageSize, &in.LabelIds)
	if err != nil {
		return &bird.GalleryListResp{
			Code: -1,
			Msg:  err.Error(),
			Data: nil,
		}, err
	}

	resps := make([]*bird.GalleryRespData, 0)
	for _, gallery := range *data {
		var resp bird.IllustrationsResp
		if len(gallery.Illustration) > 0 {
			resp = bird.IllustrationsResp{
				Id:          gallery.Illustration[0].ID.Hex(),
				RecordState: int32(gallery.Illustration[0].RecordState),
				CreateTime:  gallery.Illustration[0].CreateAt.UnixMilli(),
				Title:       gallery.Illustration[0].Title,
				Score:       gallery.Illustration[0].Score,
				WikiUrl:     gallery.Illustration[0].WikiUrl,
				ImagePath:   gallery.Illustration[0].ImagePath,
				IconPath:    gallery.Illustration[0].IconPath,
				MoreImages:  gallery.Illustration[0].MoreImages,
				Typee:       gallery.Illustration[0].Type,
				Labels:      gallery.Illustration[0].Labels,
				Description: gallery.Illustration[0].Description,
				ClassesId:   gallery.Illustration[0].ClassesId,
				ChineseName: gallery.Illustration[0].ChineseName,
				EnglishName: gallery.Illustration[0].EnglishName,
			}
		}
		resps = append(resps, &bird.GalleryRespData{
			Id:           gallery.ID.Hex(),
			RecordState:  int32(gallery.RecordState),
			CreateTime:   gallery.CreateAt.UnixMilli(),
			Name:         gallery.Name,
			UserId:       gallery.UserId,
			Illustration: &resp,
		})
	}

	return &bird.GalleryListResp{
		Code: 0,
		Msg:  "成功",
		Data: &bird.GalleryListRespData{
			Data:  resps,
			Total: count,
		},
	}, nil
}
