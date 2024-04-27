package label

import (
	"context"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindLabelByPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindLabelByPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindLabelByPageLogic {
	return &FindLabelByPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindLabelByPageLogic) FindLabelByPage(in *bird.LabelListReq) (*bird.LabelListVo, error) {
	data, count, err := l.svcCtx.LabelModel.FindListByParamAndPage(l.ctx, in.GetTypee(), in.GetParentId(), in.Page, in.PageSize, 2)
	if err != nil {
		return &bird.LabelListVo{
			Code:    -1,
			Message: "Read failure",
			Data:    nil,
		}, err
	}
	resps := make([]*bird.LabelResp, 0)
	for _, label := range *data {
		resps = append(resps, &bird.LabelResp{
			Id:          label.ID.Hex(),
			RecordState: int32(label.RecordState),
			CreateTime:  label.CreateAt.UnixMilli(),
			Name:        label.Name,
			Typee:       label.Type,
			ParentId:    label.ParentId,
		})
	}
	return &bird.LabelListVo{
		Code:    0,
		Message: "successful",
		Data: &bird.LabelListVoData{
			Total: count,
			Data:  resps,
		},
	}, nil
}
