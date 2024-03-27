package label

import (
	"context"
	"errors"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/label"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"time"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type LabelCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLabelCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LabelCreateLogic {
	return &LabelCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LabelCreateLogic) LabelCreate(in *bird.LabelCreateReq) (*bird.LabelResp, error) {
	if in.ParentId != "" {
		parent, err := l.svcCtx.LabelModel.FindOne(l.ctx, in.ParentId)
		if err != nil && !errors.Is(err, mon.ErrNotFound) {
			logx.Error(err.Error())
			return nil, err
		}
		if errors.Is(err, mon.ErrNotFound) || parent == nil {
			return nil, errors.New("父节点id不存在")
		}
		if parent.RecordState != 2 {
			return nil, errors.New("父节点不是健康的")
		}
	}
	label, err := l.svcCtx.LabelModel.FindRecord(l.ctx, in.Name, in.Type, in.ParentId)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	if label != nil {
		return nil, errors.New("存在相同记录")
	}
	data := model.Label{
		UpdateAt:    time.Now(),
		CreateAt:    time.Now(),
		Name:        in.Name,
		Type:        in.Type,
		ParentId:    in.ParentId,
		RecordState: int8(in.RecordState),
	}
	err = l.svcCtx.LabelModel.Insert(l.ctx, &data)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	return &bird.LabelResp{
		Id:          data.ID.Hex(),
		RecordState: int32(data.RecordState),
		CreateTime:  data.CreateAt.UnixMilli(),
		Name:        data.Name,
		Type:        data.Type,
		ParentId:    data.ParentId,
	}, nil
}
