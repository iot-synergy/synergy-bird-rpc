package gallery

import (
	"context"
	"errors"
	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/gallery"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
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
	// todo: add your logic here and delete this line
	dupRecord, err := l.svcCtx.GalleryModel.FindOneByNameAndUserId(l.ctx, in.Name, in.UserId)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	if dupRecord != nil {
		logx.Error("存在相同记录")
		return nil, errors.New("存在相同记录")
	}

	gallery := model.Gallery{
		UpdateAt:    time.Time{},
		CreateAt:    time.Time{},
		Name:        in.Name,
		UserId:      in.UserId,
		Favorite:    in.Favorite,
		Labels:      in.Labels,
		RecordState: int8(in.RecordState),
	}
	err = l.svcCtx.GalleryModel.Insert(l.ctx, &gallery)

	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}

	return &bird.GalleryResp{
		Id:          gallery.ID.String(),
		RecordState: int32(gallery.RecordState),
		CreateTime:  gallery.CreateAt.UnixMilli(),
		Name:        gallery.Name,
		UserId:      gallery.UserId,
		Favorite:    gallery.Favorite,
		Labels:      gallery.Labels,
	}, nil
}
