package label

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnpublishLabelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnpublishLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnpublishLabelLogic {
	return &UnpublishLabelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnpublishLabelLogic) UnpublishLabel(in *bird.IdReq) (*bird.LabelResp, error) {
	label, err := l.svcCtx.LabelModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, mon.ErrNotFound) {
			return nil, errors.New("No record")
		}
		return nil, err
	}
	if label == nil {
		return nil, errors.New("No record")
	}
	if label.RecordState != 2 {
		return nil, errors.New("unpublish")
	}
	label.RecordState = 3
	_, err = l.svcCtx.LabelModel.Update(l.ctx, label)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}

	return &bird.LabelResp{
		Id:          label.ID.Hex(),
		RecordState: 1,
		CreateTime:  label.CreateAt.UnixMilli(),
		Name:        label.Name,
		Typee:       label.Type,
		ParentId:    label.ParentId,
	}, nil
}
