package headlines

import (
	"context"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteHeadlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteHeadlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteHeadlineLogic {
	return &DeleteHeadlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteHeadlineLogic) DeleteHeadline(in *bird.Headline) (*bird.BaseResp, error) {
	_, err := l.svcCtx.HeadlineModel.Delete(l.ctx, in.GetId())

	if err != nil {
		return nil, err
	}

	return &bird.BaseResp{
		Msg: "OK",
	}, nil
}
