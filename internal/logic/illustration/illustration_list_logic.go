package illustration

import (
	"context"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type IllustrationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIllustrationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IllustrationListLogic {
	return &IllustrationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IllustrationListLogic) IllustrationList(in *bird.IllustrationsListReq) (*bird.IllustrationsListResp, error) {
	// todo: add your logic here and delete this line
	data, count, err := l.svcCtx.IllustrationModel.FindListByParamAndPage(l.ctx, in.Labels, in.Type, in.Keyword, in.Page, in.PageSize)
	if err != nil {
		return &bird.IllustrationsListResp{
			Results: nil,
			Total:   0,
			Code:    -1,
			Message: err.Error(),
		}, err
	}
	resps := make([]*bird.IllustrationsResp, 0)
	for _, illustration := range *data {
		resps = append(resps, &bird.IllustrationsResp{
			Id:          illustration.ID.String(),
			RecordState: int32(illustration.RecordState),
			CreateTime:  illustration.CreateAt.UnixMilli(),
			Title:       illustration.Title,
			Score:       illustration.Score,
			WikiUrl:     illustration.WikiUrl,
			ImagePath:   illustration.ImagePath,
			MoreImages:  illustration.MoreImages,
			Type:        illustration.Type,
			Labels:      illustration.Labels,
			Description: illustration.Description,
		})
	}
	return &bird.IllustrationsListResp{
		Results: resps,
		Total:   count,
		Code:    0,
		Message: "成功",
	}, nil
}
