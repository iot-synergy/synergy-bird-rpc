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
			Msg:  "用户未登录",
			Data: nil,
		}, nil
	}
	forein_id := strings.Join(value, "")

	gallery, err := l.svcCtx.GalleryModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, mon.ErrNotFound) {
			return &bird.GalleryResp{
				Code: -1,
				Msg:  "记录不存在",
				Data: nil,
			}, nil
		}
		return nil, err
	}
	if gallery == nil {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "记录不存在",
			Data: nil,
		}, nil
	}
	if gallery.UserId != forein_id {
		return &bird.GalleryResp{
			Code: -1,
			Msg:  "此纪录不属于该用户",
			Data: nil,
		}, nil
	}
	gallery.UpdateAt = time.Now()
	gallery.RecordState = 4

	_, err = l.svcCtx.GalleryModel.Update(l.ctx, gallery)

	return &bird.GalleryResp{
		Code: 0,
		Msg:  "成功",
		Data: nil,
	}, nil
}
