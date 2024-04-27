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
			return nil, errors.New("No record")
		}
		return &bird.BooleanResp{
			Code:    -2,
			Message: "No record",
			Data:    false,
		}, err
	}
	if label == nil {
		return &bird.BooleanResp{
			Code:    -2,
			Message: "No record",
			Data:    false,
		}, errors.New("No record")
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
		Message: "successful",
		Data:    true,
	}, nil
}
