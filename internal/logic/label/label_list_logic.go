package label

import (
	"context"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type LabelListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLabelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LabelListLogic {
	return &LabelListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LabelListLogic) LabelList(in *bird.LabelListReq) (*bird.LabelListResp, error) {
	data, count, err := l.svcCtx.LabelModel.FindListByParamAndPage(l.ctx, in.Typee, in.ParentId, in.Page, in.PageSize)
	if err != nil {
		return nil, err
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
	return &bird.LabelListResp{
		Results: resps,
		Total:   count,
		Code:    0,
		Message: "成功",
	}, nil
}
