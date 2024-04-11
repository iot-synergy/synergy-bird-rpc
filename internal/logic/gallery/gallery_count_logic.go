package gallery

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strings"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type GalleryCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGalleryCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GalleryCountLogic {
	return &GalleryCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GalleryCountLogic) GalleryCount(in *bird.NullReq) (*bird.GalleryCount, error) {
	// 获取用户id
	value := metadata.ValueFromIncomingContext(l.ctx, "gateway-firebaseid")
	if len(value) <= 0 {
		return &bird.GalleryCount{
			Code: -1,
			Msg:  "用户未登录",
			Data: nil,
		}, nil
	}
	forein_id := strings.Join(value, "")
	unlock, lock, err := l.svcCtx.IllustrationModel.StatisticLock(l.ctx, forein_id)
	if err != nil {
		return &bird.GalleryCount{
			Code: -2,
			Msg:  err.Error(),
			Data: nil,
		}, err
	}
	return &bird.GalleryCount{
		Code: 0,
		Msg:  "成功",
		Data: &bird.GalleryCountData{
			Unlock: unlock,
			Lock:   lock,
			Count:  unlock + lock,
		},
	}, nil
}
