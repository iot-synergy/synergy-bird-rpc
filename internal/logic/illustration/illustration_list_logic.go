package illustration

import (
	"context"
	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/label"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type IllustrationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIllustrationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IllustrationListLogic {
	return &IllustrationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IllustrationListLogic) IllustrationList(in *bird.IllustrationsListReq) (*bird.IllustrationsListResp, error) {
	data, count, err := l.svcCtx.IllustrationModel.FindListByParamAndPage(l.ctx, &in.Labels, nil, in.GetTypee(), in.GetKeyword(), nil, in.GetState(), in.Page, in.PageSize)
	if err != nil {
		return &bird.IllustrationsListResp{
			Results: nil,
			Total:   0,
			Code:    -1,
			Message: err.Error(),
		}, err
	}
	//获得去重的标签id
	labelIdMap := make(map[string]string)
	labelIds := make([]string, 0)
	for _, illustration := range *data {
		for _, label := range illustration.Labels {
			labelIdMap[label] = ""
		}
	}
	for key, _ := range labelIdMap {
		labelIds = append(labelIds, key)
	}
	//通过标签id获取标签，并转化成map
	labels, err := l.svcCtx.LabelModel.FindListByIds(l.ctx, labelIds)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	labelMap := make(map[string]model.Label)
	for _, label := range *labels {
		labelMap[label.ID.Hex()] = label
	}
	resps := make([]*bird.IllustrationsRespVo, 0) //结果集
	for _, illustration := range *data {
		labelResps := make([]*bird.LabelResp, 0) //图鉴的标签列表
		for _, labelId := range illustration.Labels {
			//从labelMap获取标签
			label, ok := labelMap[labelId]
			if ok {
				labelResps = append(labelResps, &bird.LabelResp{
					Id:          label.ID.Hex(),
					RecordState: int32(label.RecordState),
					CreateTime:  label.CreateAt.UnixMilli(),
					Name:        label.Name,
					Typee:       label.Type,
					ParentId:    label.ParentId,
				})
			}
		}
		resps = append(resps, &bird.IllustrationsRespVo{
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
			Labels:      labelResps,
			Description: illustration.Description,
			ClassesId:   illustration.ClassesId,
			ChineseName: illustration.ChineseName,
			EnglishName: illustration.EnglishName,
		})
	}
	return &bird.IllustrationsListResp{
		Results: resps,
		Total:   count,
		Code:    0,
		Message: "成功",
	}, nil
}
