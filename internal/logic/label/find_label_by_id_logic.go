package label

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindLabelByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindLabelByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindLabelByIdLogic {
	return &FindLabelByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindLabelByIdLogic) FindLabelById(in *bird.IdReq) (*bird.LabelVo, error) {
	data, err := l.svcCtx.LabelModel.FindOne(l.ctx, in.Id)
	if errors.Is(err, mon.ErrNotFound) {
		if err != nil {
			return &bird.LabelVo{
				Code:    -2,
				Message: "Data null",
			}, nil
		}
	}
	if err != nil {
		return &bird.LabelVo{
			Code:    -1,
			Message: "fail",
		}, err
	}
	if data.RecordState != 2 {
		return &bird.LabelVo{
			Code:    -1,
			Message: "fail",
		}, err
	}

	return &bird.LabelVo{
		Code:    0,
		Message: "successful",
		Data: &bird.LabelResp{
			Id:          data.ID.Hex(),
			RecordState: int32(data.RecordState),
			CreateTime:  data.CreateAt.UnixMilli(),
			Name:        data.Name,
			Typee:       data.Type,
			ParentId:    data.ParentId,
		},
	}, nil
}
