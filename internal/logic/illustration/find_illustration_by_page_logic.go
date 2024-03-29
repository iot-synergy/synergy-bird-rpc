package illustration

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindIllustrationByPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindIllustrationByPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindIllustrationByPageLogic {
	return &FindIllustrationByPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindIllustrationByPageLogic) FindIllustrationByPage(in *bird.IllustrationsListReq) (*bird.IllustrationsListVo, error) {
	data, count, err := l.svcCtx.IllustrationModel.FindListByParamAndPage(l.ctx, in.Labels, in.GetTypee(), in.GetKeyword(), in.GetState(), in.Page, in.PageSize)
	if err != nil {
		return &bird.IllustrationsListVo{
			Code:    -1,
			Message: err.Error(),
		}, err
	}
	resps := make([]*bird.IllustrationsRespVo, 0)
	for _, illustration := range *data {
		labelResps := make([]*bird.LabelResp, 0)
		for _, label := range illustration.Labels {
			labelData, err := l.svcCtx.LabelModel.FindOne(l.ctx, label)
			if err == nil || errors.Is(err, mon.ErrNotFound) {
				labelResps = append(labelResps, &bird.LabelResp{
					Id:          labelData.ID.Hex(),
					RecordState: int32(labelData.RecordState),
					CreateTime:  labelData.CreateAt.UnixMilli(),
					Name:        labelData.Name,
					Typee:       labelData.Type,
					ParentId:    labelData.ParentId,
				})
			}
		}
		resps = append(resps, &bird.IllustrationsRespVo{
			Id:          illustration.ID.Hex(),
			RecordState: int32(illustration.RecordState),
			CreateTime:  illustration.CreateAt.UnixMilli(),
			Title:       illustration.Title,
			Score:       illustration.Score,
			WikiUrl:     illustration.WikiUrl,
			ImagePath:   illustration.ImagePath,
			IconPath:    illustration.IconPath,
			MoreImages:  illustration.MoreImages,
			Typee:       illustration.Type,
			Labels:      labelResps,
			Description: illustration.Description,
		})
	}
	return &bird.IllustrationsListVo{
		Code:    0,
		Message: "成功",
		Data: &bird.IllustrationsListVoData{
			Total: count,
			Data:  resps,
		},
	}, nil
}
