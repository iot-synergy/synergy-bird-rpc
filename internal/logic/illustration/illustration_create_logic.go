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
	illustration := model.Illustration{
		UpdateAt:    time.Time{},
		CreateAt:    time.Time{},
		Title:       in.Title,
		Score:       in.Score,
		WikiUrl:     in.WikiUrl,
		ImagePath:   in.ImagePath,
		IconPath:    in.IconPath,
		MoreImages:  in.MoreImages,
		Type:        in.Typee,
		Labels:      make([]string, 0),
		Description: in.Description,
	}
	if in.RecordState == 2 {
		illustration.RecordState = 2
	} else {
		illustration.RecordState = 1
	}
	labels, err := l.svcCtx.LabelModel.FindListByIds(l.ctx, in.Labels)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	for _, label := range *labels {
		illustration.Labels = append(illustration.Labels, label.ID.Hex())
	}
	err = l.svcCtx.IllustrationModel.Insert(l.ctx, &illustration)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	return &bird.IllustrationsResp{
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
		Labels:      illustration.Labels,
		Description: illustration.Description,
	}, nil
}
