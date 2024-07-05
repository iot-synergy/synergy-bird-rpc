package headlines

import (
	"context"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
	"github.com/iot-synergy/synergy-common/utils/pointy"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateHeadlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateHeadlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHeadlineLogic {
	return &UpdateHeadlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateHeadlineLogic) UpdateHeadline(in *bird.Headline) (*bird.Headline, error) {
	result, err := l.svcCtx.HeadlineModel.FindOne(l.ctx, in.GetId())

	if err != nil || result == nil {
		return nil, err
	}

	// update := &model.Headline{
	// 	Url:         in.GetUrl(),
	// 	Site:        in.GetSite(),
	// 	Cover:       in.GetCover(),
	// 	Title:       in.GetTitle(),
	// 	Description: in.GetDescription(),
	// 	Image:       in.GetImage(),
	// }

	// update.ID =

	if in.Cover != nil {
		result.Cover = in.GetCover()
	}

	if in.GetUrl() != "" {
		result.Url = in.GetUrl()
	}

	if in.GetSite() != "" {
		result.Site = in.GetSite()
	}

	if in.GetTitle() != "" {
		result.Title = in.GetTitle()
	}

	if in.GetDescription() != "" {
		result.Description = in.GetDescription()
	}

	if in.GetImage() != "" {
		result.Image = in.GetImage()
	}

	_, err = l.svcCtx.HeadlineModel.Update(l.ctx, result)

	if err != nil {
		return nil, err
	}

	result, err = l.svcCtx.HeadlineModel.FindOne(l.ctx, in.GetId())

	if err != nil || result == nil {
		return nil, err
	}

	return &bird.Headline{
		Id:          pointy.GetPointer(result.ID.Hex()),
		Url:         &result.Url,
		Site:        &result.Site,
		Cover:       &result.Cover,
		Title:       &result.Title,
		Description: &result.Description,
		Image:       &result.Image,
		CreateAt:    pointy.GetPointer(result.CreateAt.Unix()),
		UpdateAt:    pointy.GetPointer(result.UpdateAt.Unix()),
	}, nil
}
