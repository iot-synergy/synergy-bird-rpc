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
	if in.GetClassesId() != 0 && in.GetClassesId() != illustration.ClassesId {
		classes, err := l.svcCtx.ClassesModel.FindOneByClassesId(l.ctx, in.GetClassesId())
		if err != nil {
			logx.Error(err.Error())
			return nil, err
		}
		data, err := l.svcCtx.IllustrationModel.FindOneByTitle(l.ctx, classes.ClassesName)
		if err != nil {
			logx.Error(err.Error())
			return nil, err
		}
		if data != nil && data.Title != "" && data.RecordState != 4 {
			return nil, errors.New("图鉴已创建过了")
		} else if data != nil && data.Title != "" && illustration.RecordState == 4 {
			_, err = l.svcCtx.IllustrationModel.Delete(l.ctx, illustration.ID.Hex())
			if err != nil {
				logx.Error(err.Error())
				return nil, err
			}
		}
		illustration.Title = classes.ClassesName
		illustration.ClassesId = classes.ClassesId
		illustration.ChineseName = classes.ChineseName
		illustration.EnglishName = classes.EnglishName
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
	if in.GetRecordState() > 0 && in.GetRecordState() < 5 {
		illustration.RecordState = int8(in.GetRecordState())
	}
	labels, err := l.svcCtx.LabelModel.FindListByIds(l.ctx, in.Labels)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	labelIds := make([]string, 0)
	for _, label := range *labels {
		labelIds = append(labelIds, label.ID.Hex())
	}
	illustration.Labels = labelIds
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
