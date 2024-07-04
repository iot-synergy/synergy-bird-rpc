package headlines

import (
	"context"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
	"github.com/iot-synergy/synergy-common/utils/pointy"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryAllHeadlinesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryAllHeadlinesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAllHeadlinesLogic {
	return &QueryAllHeadlinesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryAllHeadlinesLogic) QueryAllHeadlines(in *bird.HeadlineQueryPageReq) (*bird.HeadlineListResp, error) {
	data, count, err := l.svcCtx.HeadlineModel.FindListByParam(l.ctx, in.GetSite(), in.Page, in.PageSize)
	if err != nil {
		l.Logger.Error(err)
		return &bird.HeadlineListResp{
			Code: -1,
			Msg:  err.Error(),
			Data: nil,
		}, err
	}

	headlineListData := bird.HeadlineListData{
		Total: count,
		Data:  []*bird.Headline{},
	}
	for _, item := range *data {
		headlineListData.Data = append(headlineListData.Data, &bird.Headline{
			Id:          pointy.GetPointer(item.ID.Hex()),
			Url:         pointy.GetPointer(item.Url),
			Site:        pointy.GetPointer(item.Site),
			Cover:       pointy.GetPointer(item.Cover),
			Title:       pointy.GetPointer(item.Title),
			Description: pointy.GetPointer(item.Description),
			Image:       pointy.GetPointer(item.Image),
			CreateAt:    pointy.GetPointer(item.CreateAt.Unix()),
			UpdateAt:    pointy.GetPointer(item.UpdateAt.Unix()),
		})
	}

	return &bird.HeadlineListResp{
		Code: 0,
		Msg:  "successful",
		Data: &headlineListData,
	}, nil

}
