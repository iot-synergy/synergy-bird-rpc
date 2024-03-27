package illustration

import (
	"context"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/illustration"
	"time"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type IllustrationCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIllustrationCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IllustrationCreateLogic {
	return &IllustrationCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IllustrationCreateLogic) IllustrationCreate(in *bird.IllustrationsCreateReq) (*bird.IllustrationsResp, error) {
	// todo: add your logic here and delete this line
	illustration := model.Illustration{
		UpdateAt:    time.Time{},
		CreateAt:    time.Time{},
		Title:       in.Title,
		Score:       in.Score,
		WikiUrl:     in.WikiUrl,
		ImagePath:   in.ImagePath,
		IconPath:    in.IconPath,
		MoreImages:  in.MoreImages,
		Type:        in.Type,
		Labels:      in.Labels,
		Description: in.Description,
		RecordState: int8(in.RecordState),
	}
	err := l.svcCtx.IllustrationModel.Insert(l.ctx, &illustration)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	return &bird.IllustrationsResp{
		Id:          illustration.ID.String(),
		RecordState: int32(illustration.RecordState),
		CreateTime:  illustration.CreateAt.UnixMilli(),
		Title:       illustration.Title,
		Score:       illustration.Score,
		WikiUrl:     illustration.WikiUrl,
		ImagePath:   illustration.ImagePath,
		IconPath:    illustration.IconPath,
		MoreImages:  illustration.MoreImages,
		Type:        illustration.Type,
		Labels:      illustration.Labels,
		Description: illustration.Description,
	}, nil
}
