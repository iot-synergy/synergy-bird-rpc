package label

import (
	"context"
	"errors"
	model "synergy-bird-rpc/storage/label"
	"time"

	"synergy-bird-rpc/internal/svc"
	"synergy-bird-rpc/types/bird"

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
	// todo: add your logic here and delete this line
	label, err := l.svcCtx.LabelModel.FindRecord(l.ctx, in.UserId, in.Name, in.Type, in.ParentId)
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
		UserId:      in.UserId,
		RecordState: int8(in.RecordState),
	}
	err = l.svcCtx.LabelModel.Insert(l.ctx, &data)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	return &bird.LabelResp{
		Id:          data.ID.String(),
		RecordState: int32(data.RecordState),
		CreateTime:  data.CreateAt.UnixMilli(),
		Name:        data.Name,
		Type:        data.Type,
		ParentId:    data.ParentId,
		UserId:      data.UserId,
	}, nil
}
