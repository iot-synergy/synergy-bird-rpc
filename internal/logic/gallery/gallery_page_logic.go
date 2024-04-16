package gallery

import (
	"context"
	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/illustration"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type GalleryPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGalleryPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GalleryPageLogic {
	return &GalleryPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GalleryPageLogic) GalleryPage(in *bird.GalleryPageReq) (resp *bird.GalleryListResp, err error) {
	//根据标签列表获取图鉴Title列表
	var titles *[]string
	if in.GetLabelIds() != nil && len(in.GetLabelIds()) > 0 {
		titles, err = l.svcCtx.IllustrationModel.FindTitleListByLabelIds(l.ctx, in.GetLabelIds())
		if titles == nil || len(*titles) == 0 {
			return &bird.GalleryListResp{
				Code: 0,
				Msg:  "成功",
				Data: &bird.GalleryListRespData{
					Data:  nil,
					Total: 0,
				},
			}, nil
		}
	}
	//根据Title列表查询成就
	data, count, err := l.svcCtx.GalleryModel.FindListByParamAndPage(l.ctx, in.UserId, in.GetIllustrationId(),
		in.GetName(), in.GetStartTime(), in.GetEndTime(), in.Page, in.PageSize, titles)
	if err != nil {
		return &bird.GalleryListResp{
			Code: -1,
			Msg:  err.Error(),
			Data: nil,
		}, err
	}
	titleList := make([]string, 0)
	for _, gallery := range *data {
		titleList = append(titleList, gallery.Name)
	}
	illustrationList, err := l.svcCtx.IllustrationModel.FindListByTitles(l.ctx, &titleList)
	if err != nil {
		return &bird.GalleryListResp{
			Code: -1,
			Msg:  err.Error(),
			Data: nil,
		}, err
	}
	illustrationMap := make(map[string]model.Illustration)
	for _, illustration := range *illustrationList {
		illustrationMap[illustration.Title] = illustration
	}

	resps := make([]*bird.GalleryRespData, 0)
	for _, gallery := range *data {
		var resp bird.IllustrationsResp
		illustration := illustrationMap[gallery.Name]
		resp = bird.IllustrationsResp{
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
			ClassesId:   illustration.ClassesId,
			ChineseName: illustration.ChineseName,
			EnglishName: illustration.EnglishName,
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
