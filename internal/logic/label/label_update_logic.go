package label

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type LabelUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLabelUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LabelUpdateLogic {
	return &LabelUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LabelUpdateLogic) LabelUpdate(in *bird.LabelUpdateReq) (*bird.LabelResp, error) {
	// todo: add your logic here and delete this line
	label, err := l.svcCtx.LabelModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, mon.ErrNotFound) {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	if label == nil {
		return nil, errors.New("记录不存在")
	}
	if in.GetRecordState() != 0 {
		label.RecordState = int8(in.GetRecordState())
	}
	if in.GetName() != "" {
		label.Name = in.GetName()
	}
	if in.GetTypee() != "" {
		label.Type = in.GetTypee()
	}
	if in.GetParentId() != "" {
		label.ParentId = in.GetParentId()
	}
	l.svcCtx.LabelModel.Update(l.ctx, label)
	return &bird.LabelResp{
		Id:          label.ID.Hex(),
		RecordState: int32(label.RecordState),
		CreateTime:  label.CreateAt.UnixMilli(),
		Name:        label.Name,
		Typee:       label.Type,
		ParentId:    label.ParentId,
	}, nil
}
