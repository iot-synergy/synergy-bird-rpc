package illustration

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"time"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type IllustrationUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIllustrationUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IllustrationUpdateLogic {
	return &IllustrationUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IllustrationUpdateLogic) IllustrationUpdate(in *bird.IllustrationsUpdateReq) (*bird.IllustrationsResp, error) {
	// todo: add your logic here and delete this line
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
	illustration.UpdateAt = time.Now()
	if in.GetRecordState() != 0 {
		illustration.RecordState = int8(in.GetRecordState())
	}
	if in.GetTitle() != "" {
		illustration.Title = in.GetTitle()
	}
	if in.Score != nil {
		illustration.Score = in.GetScore()
	}
	if in.GetWikiUrl() != "" {
		illustration.WikiUrl = in.GetWikiUrl()
	}
	if in.GetImagePath() != "" {
		illustration.ImagePath = in.GetImagePath()
	}
	if in.GetIconPath() != "" {
		illustration.IconPath = in.GetIconPath()
	}
	if in.MoreImages != nil && len(in.MoreImages) != 0 {
		illustration.MoreImages = in.MoreImages
	}
	if in.GetTypee() != "" {
		illustration.Type = in.GetTypee()
	}
	if in.Labels != nil && len(in.Labels) != 0 {
		illustration.Labels = in.Labels
	}
	if in.GetDescription() != "" {
		illustration.Description = in.GetDescription()
	}
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
