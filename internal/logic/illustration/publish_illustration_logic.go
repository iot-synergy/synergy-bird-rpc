package illustration

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishIllustrationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishIllustrationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishIllustrationLogic {
	return &PublishIllustrationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishIllustrationLogic) PublishIllustration(in *bird.IdReq) (*bird.IllustrationsResp, error) {
	illustration, err := l.svcCtx.IllustrationModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, mon.ErrNotFound) {
			return nil, errors.New("没有对应记录")
		}
		return nil, err
	}
	if illustration == nil {
		return nil, errors.New("没有对应记录")
	}
	if illustration.RecordState == 2 {
		return nil, errors.New("已发布")
	}
	illustration.RecordState = 2
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
	}, nil
}
