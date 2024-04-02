package gallery

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"time"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type GalleryUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGalleryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GalleryUpdateLogic {
	return &GalleryUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GalleryUpdateLogic) GalleryUpdate(in *bird.GalleryUpdateReq) (*bird.GalleryResp, error) {
	// todo: add your logic here and delete this line
	gallery, err := l.svcCtx.GalleryModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, mon.ErrNotFound) {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	if gallery == nil {
		return nil, errors.New("记录不存在")
	}
	if gallery.UserId != in.GetUserId() {
		return nil, errors.New("此纪录不属于该用户")
	}
	gallery.UpdateAt = time.Now()

	if in.GetRecordState() != 0 {
		gallery.RecordState = int8(in.GetRecordState())
	}

	_, err = l.svcCtx.GalleryModel.Update(l.ctx, gallery)

	return &bird.GalleryResp{
		Id:          gallery.ID.Hex(),
		RecordState: int32(gallery.RecordState),
		CreateTime:  gallery.CreateAt.UnixMilli(),
		Name:        gallery.Name,
		UserId:      gallery.UserId,
	}, nil
}
