package label

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLabelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLabelLogic {
	return &PublishLabelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLabelLogic) PublishLabel(in *bird.IdReq) (*bird.LabelResp, error) {
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
	if label.RecordState == 2 {
		return nil, errors.New("Have released")
	}
	label.RecordState = 2
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
