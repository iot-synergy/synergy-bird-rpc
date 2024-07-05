package headlines

import (
	"context"
	"strconv"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
	"github.com/iot-synergy/synergy-common/utils/pointy"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryHeadlineListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryHeadlineListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryHeadlineListLogic {
	return &QueryHeadlineListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryHeadlineListLogic) QueryHeadlineList(in *bird.HeadlineQueryReq) (*bird.HeadlineListResp, error) {
	list := []*bird.Headline{}
	curIndex := int64(0)
	curIndex, err := strconv.ParseInt(in.CurId, 10, 64)

	if err != nil {
		curIndex = 0
	}

	_list, err := l.svcCtx.HeadlineModel.FindListByIndex(l.ctx, curIndex)

	if err != nil {
		curIndex = 0
	}

	if err != nil {
		return &bird.HeadlineListResp{
			Code: -1,
			Msg:  "error",
			Data: &bird.HeadlineListData{
				Total: 0,
				Data:  []*bird.Headline{},
			},
		}, nil
	}

	for index, item := range *_list {
		id := strconv.FormatInt(curIndex+int64(index)+1, 10)
		list = append(list, &bird.Headline{
			Id:          &id,
			Url:         pointy.GetPointer(item.Url),
			Site:        pointy.GetPointer(item.Site),
			Cover:       pointy.GetPointer(item.Cover),
			Title:       pointy.GetPointer(item.Title),
			Description: pointy.GetPointer(item.Description),
			Image:       pointy.GetPointer(item.Image),
			CreateAt:    pointy.GetPointer(item.CreateAt.Unix()),
			UpdateAt:    pointy.GetPointer(item.UpdateAt.Unix()),
		})
	}

	return &bird.HeadlineListResp{
		Code: 0,
		Msg:  "OK",
		Data: &bird.HeadlineListData{
			Total: int64(len(list)),
			Data:  list,
		},
	}, nil
}
