package illustration

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteIllustrationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteIllustrationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteIllustrationLogic {
	return &DeleteIllustrationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteIllustrationLogic) DeleteIllustration(in *bird.IdReq) (*bird.BooleanResp, error) {
	illustration, err := l.svcCtx.IllustrationModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, mon.ErrNotFound) {
			return nil, errors.New("没有对应记录")
		}
		return &bird.BooleanResp{
			Code:    -2,
			Message: "没有对应记录",
			Data:    false,
		}, err
	}
	if illustration == nil {
		return &bird.BooleanResp{
			Code:    -2,
			Message: "没有对应记录",
			Data:    false,
		}, errors.New("没有对应记录")
	}
	illustration.RecordState = 4
	_, err = l.svcCtx.IllustrationModel.Update(l.ctx, illustration)
	if err != nil {
		logx.Error(err.Error())
		return &bird.BooleanResp{
			Code:    -1,
			Message: err.Error(),
			Data:    false,
		}, err
	}

	return &bird.BooleanResp{
		Code:    0,
		Message: "成功",
		Data:    true,
	}, nil
}
