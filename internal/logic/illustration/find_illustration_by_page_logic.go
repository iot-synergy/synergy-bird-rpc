package illustration

import (
	"context"
	"github.com/iot-synergy/synergy-bird-rpc/common"
	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	model2 "github.com/iot-synergy/synergy-bird-rpc/storage/galleryCount"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/label"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
	"strings"
)

type FindIllustrationByPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindIllustrationByPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindIllustrationByPageLogic {
	return &FindIllustrationByPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindIllustrationByPageLogic) FindIllustrationByPage(in *bird.IllustrationsPageReq) (*bird.IllustrationsListVo, error) {
	// 获取用户id
	value := metadata.ValueFromIncomingContext(l.ctx, "gateway-firebaseid")
	if len(value) <= 0 {
		return &bird.IllustrationsListVo{
			Code:    -1,
			Message: "用户未登录",
			Data:    nil,
		}, nil
	}
	forein_id := strings.Join(value, "")
	data, count, err :=
		l.svcCtx.IllustrationModel.FindPageJoinGallery(l.ctx, in.Labels, forein_id, in.GetTypee(), in.GetKeyword(), in.IsUnlock, 2, in.Page, in.PageSize)
	//data, count, err := l.svcCtx.IllustrationModel.FindListByParamAndPage(l.ctx, in.Labels, in.GetTypee(), in.GetKeyword(), 2, in.Page, in.PageSize)
	if err != nil {
		return &bird.IllustrationsListVo{
			Code:    -1,
			Message: err.Error(),
		}, err
	}
	//获得去重的标签id
	labelIdMap := make(map[string]string)
	//获得图鉴id集合
	illustrationIds := make([]string, 0)
	labelIds := make([]string, 0)
	for _, illustration := range *data {
		illustrationIds = append(illustrationIds, illustration.ID.Hex())
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
	//判断图鉴是否是用户已解锁的
	galleryCounts, _ := l.svcCtx.GalleryCountModel.FindByIllustrationIdList(l.ctx, forein_id, illustrationIds)
	illustrationIdList := make([]string, 0)
	galleryCountMap := make(map[string]model2.GalleryCount)
	for _, galleryCount := range *galleryCounts {
		illustrationIdList = append(illustrationIdList, galleryCount.IllustrationId)
		galleryCountMap[galleryCount.IllustrationId] = galleryCount
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
		illustrationId := illustration.ID.Hex()
		isUnlock := common.ListContainApi(&illustrationIdList, illustrationId)
		var unlockTime int64
		if isUnlock {
			unlockTime = galleryCountMap[illustrationId].UnlockTime.UnixMilli()
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
			IsUnlock:    &isUnlock,
			UnlockTime:  &unlockTime,
		})
	}
	return &bird.IllustrationsListVo{
		Code:    0,
		Message: "成功",
		Data: &bird.IllustrationsListVoData{
			Total: count,
			Data:  resps,
		},
	}, nil
}
