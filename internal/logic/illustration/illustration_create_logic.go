package illustration

import (
	"context"
	"errors"
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
	classes, err := l.svcCtx.ClassesModel.FindOneByClassesId(l.ctx, in.ClassesId)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	illustration, err := l.svcCtx.IllustrationModel.FindOneByTitle(l.ctx, classes.ClassesName)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	if illustration != nil && illustration.RecordState != 4 {
		return nil, errors.New("图鉴已创建过了")
	} else if illustration != nil && illustration.RecordState == 4 {
		_, err = l.svcCtx.IllustrationModel.Delete(l.ctx, illustration.ID.Hex())
		if err != nil {
			logx.Error(err.Error())
			return nil, err
		}
	}

	illustration = &model.Illustration{
		UpdateAt:    time.Now(),
		CreateAt:    time.Now(),
		Title:       classes.ClassesName,
		Score:       in.Score,
		WikiUrl:     in.WikiUrl,
		ImagePath:   in.ImagePath,
		IconPath:    in.IconPath,
		MoreImages:  in.MoreImages,
		Type:        in.Typee,
		Labels:      make([]string, 0),
		Description: in.Description,
		ClassesId:   in.ClassesId,
		ChineseName: classes.ChineseName,
		EnglishName: classes.EnglishName,
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
	err = l.svcCtx.IllustrationModel.Insert(l.ctx, illustration)
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
