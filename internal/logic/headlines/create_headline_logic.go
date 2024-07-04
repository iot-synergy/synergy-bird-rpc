package headlines

import (
	"context"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/headline"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
	"github.com/iot-synergy/synergy-common/utils/pointy"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateHeadlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateHeadlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHeadlineLogic {
	return &CreateHeadlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateHeadlineLogic) CreateHeadline(in *bird.Headline) (*bird.Headline, error) {
	dbHeadline := &model.Headline{
		Url:         in.GetUrl(),
		Site:        in.GetSite(),
		Cover:       in.GetCover(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		Image:       in.GetImage(),
	}

	result, err := l.svcCtx.HeadlineModel.InsertOne(l.ctx, dbHeadline)

	if err != nil {
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
