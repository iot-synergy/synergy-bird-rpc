package label

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLabelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLabelLogic {
	return &DeleteLabelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLabelLogic) DeleteLabel(in *bird.IdReq) (*bird.BooleanResp, error) {
	label, err := l.svcCtx.LabelModel.FindOne(l.ctx, in.Id)
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
	if label == nil {
		return &bird.BooleanResp{
			Code:    -2,
			Message: "没有对应记录",
			Data:    false,
		}, errors.New("没有对应记录")
	}
	label.RecordState = 4
	_, err = l.svcCtx.LabelModel.Update(l.ctx, label)
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
