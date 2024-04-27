package gallery

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type GalleryDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGalleryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GalleryDeleteLogic {
	return &GalleryDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GalleryDeleteLogic) GalleryDelete(in *bird.IdReq) (*bird.GalleryResp, error) {
	// 获取用户id
	value := metadata.ValueFromIncomingContext(l.ctx, "gateway-firebaseid")
	if len(value) <= 0 {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "User not logged in",
			Data: nil,
		}, nil
	}
	forein_id := strings.Join(value, "")

	gallery, err := l.svcCtx.GalleryModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, mon.ErrNotFound) {
			return &bird.GalleryResp{
				Code: -1,
				Msg:  "Record does not exist",
				Data: nil,
			}, nil
		}
		return nil, err
	}
	if gallery == nil {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "Record does not exist",
			Data: nil,
		}, nil
	}
	if gallery.UserId != forein_id {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "Record does not belong to the user",
			Data: nil,
		}, nil
	}
	gallery.UpdateAt = time.Now()
	gallery.RecordState = 4

	_, err = l.svcCtx.GalleryModel.Update(l.ctx, gallery)
	if err != nil {
		return &bird.GalleryResp{
			Code: -2,
			Msg:  "fail",
			Data: nil,
		}, err
	}

	//更新图鉴数量
	galleryCountData, err := l.svcCtx.GalleryCountModel.FindOneByUserIdAndIllustrationId(l.ctx, forein_id, gallery.IllustrationId)
	if galleryCountData != nil && galleryCountData.Name != "" {
		count, _ := l.svcCtx.GalleryModel.CountByUserIdAndIllustrationId(l.ctx, forein_id, gallery.IllustrationId)
		galleryCountData.Count = count
		l.svcCtx.GalleryCountModel.Update(l.ctx, galleryCountData)
	}

	return &bird.GalleryResp{
		Code: 0,
		Msg:  "successful",
		Data: nil,
	}, nil
}
