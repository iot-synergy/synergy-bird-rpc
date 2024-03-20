package gallery

import (
	"context"

	"synergy-bird-rpc/internal/svc"
	"synergy-bird-rpc/types/bird"

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
	// todo: add your logic here and delete this line
	data, count, err := l.svcCtx.GalleryModel.FindListByParamAndPage(l.ctx, in.UserId, in.Labels, in.Favorite, in.Page, in.PageSize)
	if err != nil {
		return &bird.GalleryListResp{
			Results: nil,
			Total:   0,
			Code:    -1,
			Message: err.Error(),
		}, err
	}

	resps := make([]*bird.GalleryResp, 0)
	for _, gallery := range *data {
		resps = append(resps, &bird.GalleryResp{
			Id:          gallery.ID.String(),
			RecordState: int32(gallery.RecordState),
			CreateTime:  gallery.CreateAt.UnixMilli(),
			Name:        gallery.Name,
			UserId:      gallery.UserId,
			Favorite:    gallery.Favorite,
			Labels:      gallery.Labels,
		})
	}

	return &bird.GalleryListResp{
		Results: resps,
		Total:   count,
		Code:    0,
		Message: "成功",
	}, nil
}
