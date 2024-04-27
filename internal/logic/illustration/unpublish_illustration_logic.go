package illustration

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnpublishIllustrationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnpublishIllustrationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnpublishIllustrationLogic {
	return &UnpublishIllustrationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnpublishIllustrationLogic) UnpublishIllustration(in *bird.IdReq) (*bird.IllustrationsResp, error) {
	illustration, err := l.svcCtx.IllustrationModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, mon.ErrNotFound) {
			return nil, errors.New("No record")
		}
		return nil, err
	}
	if illustration == nil {
		return nil, errors.New("No record")
	}
	if illustration.RecordState != 2 {
		return nil, errors.New("unpublish")
	}
	illustration.RecordState = 3
	_, err = l.svcCtx.IllustrationModel.Update(l.ctx, illustration)
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
		ClassesId:   illustration.ClassesId,
		ChineseName: illustration.ChineseName,
		EnglishName: illustration.EnglishName,
	}, nil
}
