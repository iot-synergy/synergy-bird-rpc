package classes

import (
	"context"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClassesListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClassesListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClassesListLogic {
	return &ClassesListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ClassesListLogic) ClassesList(in *bird.ClassesListReq) (*bird.ClassesListResp, error) {
	data, count, err := l.svcCtx.ClassesModel.FindListByParam(l.ctx, in.GetKeyword(), in.Page, in.PageSize)
	if err != nil {
		l.Logger.Error(err)
		return &bird.ClassesListResp{
			Code:  -1,
			Msg:   err.Error(),
			Total: 0,
			Data:  nil,
		}, err
	}
	classesData := make([]*bird.ClassesData, 0)
	for _, classes := range *data {
		classesData = append(classesData, &bird.ClassesData{
			Id:          classes.ID.Hex(),
			ClassesId:   classes.ClassesId,
			ClassesName: classes.ClassesName,
			ChineseName: classes.ChineseName,
			EnglishName: classes.EnglishName,
		})
	}
	return &bird.ClassesListResp{
		Code:  0,
		Msg:   "成功",
		Total: count,
		Data:  classesData,
	}, nil
}
