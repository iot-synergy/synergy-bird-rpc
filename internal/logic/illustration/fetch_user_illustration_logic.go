package illustration

import (
	"context"

	"synergy-bird-rpc/internal/svc"
	"synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type FetchUserIllustrationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFetchUserIllustrationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FetchUserIllustrationLogic {
	return &FetchUserIllustrationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FetchUserIllustrationLogic) FetchUserIllustration(in *bird.FetchUserIllustrationReq) (*bird.FetchUserIllustrationResp, error) {
	// todo: add your logic here and delete this line

	return &bird.FetchUserIllustrationResp{}, nil
}
