package classes

import (
	"bufio"
	"context"
	"fmt"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/classes"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitClassesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitClassesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitClassesLogic {
	return &InitClassesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitClassesLogic) InitClasses(in *bird.NullReq) (*bird.BooleanResp, error) {
	classesListFile, err := os.Open("./txt/classes_list.txt")
	if err != nil {
		return &bird.BooleanResp{
			Code:    -1,
			Message: "打开文件失败",
			Data:    false,
		}, err
	}
	defer classesListFile.Close()
	classesChineseListFile, err := os.Open("./txt/classes_Chinese_list.txt")
	if err != nil {
		return &bird.BooleanResp{
			Code:    -1,
			Message: "打开文件失败",
			Data:    false,
		}, err
	}
	defer classesChineseListFile.Close()
	classesEnglishListFile, err := os.Open("./txt/classes_English_list.txt")
	if err != nil {
		return &bird.BooleanResp{
			Code:    -1,
			Message: "打开文件失败",
			Data:    false,
		}, err
	}
	defer classesEnglishListFile.Close()
	err = l.svcCtx.ClassesModel.BatchDelete(l.ctx)
	if err != nil {
		return &bird.BooleanResp{
			Code:    -1,
			Message: "清空旧数据失败",
			Data:    false,
		}, err
	}
	classesMap := make(map[int64]model.Classes, 0)
	classesScanner := bufio.NewScanner(classesListFile)
	for classesScanner.Scan() {
		split := strings.SplitN(classesScanner.Text(), " ", 2)
		if split == nil || len(split) != 2 {
			continue
		}
		id, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			continue
		}
		classesMap[id] = model.Classes{
			ID:          primitive.NewObjectID(),
			UpdateAt:    time.Now(),
			CreateAt:    time.Now(),
			ClassesId:   id,
			ClassesName: split[1],
		}
	}
	chineseScanner := bufio.NewScanner(classesChineseListFile)
	for chineseScanner.Scan() {
		split := strings.SplitN(chineseScanner.Text(), " ", 2)
		if split == nil || len(split) != 2 {
			continue
		}
		id, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			continue
		}
		classes, ok := classesMap[id]
		if ok {
			classes.ChineseName = split[1]
			classesMap[id] = classes
		}
	}
	englishScanner := bufio.NewScanner(classesEnglishListFile)
	for englishScanner.Scan() {
		split := strings.SplitN(englishScanner.Text(), " ", 2)
		if split == nil || len(split) != 2 {
			continue
		}
		id, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			continue
		}
		classes, ok := classesMap[id]
		if ok {
			classes.EnglishName = split[1]
			classesMap[id] = classes
		}
	}
	classesList := make([]model.Classes, 0)
	for _, classes := range classesMap {
		classesList = append(classesList, classes)
	}
	l.svcCtx.ClassesModel.BatchInsert(l.ctx, &classesList)
	return &bird.BooleanResp{
		Code:    0,
		Message: "成功",
		Data:    true,
	}, nil
}

func InitClassesOne(l *InitClassesLogic) error {
	classesListFile, err := os.Open("./txt/classes_list.txt")
	if err != nil {
		return err
	}
	defer classesListFile.Close()
	scanner := bufio.NewScanner(classesListFile)
	for scanner.Scan() {
		split := strings.SplitN(scanner.Text(), " ", 2)
		if split == nil || len(split) != 2 {
			continue
		}
		id, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			continue
		}
		data, err := l.svcCtx.ClassesModel.FindOneByClassesIdOrClassesName(l.ctx, id, split[1])
		if err != nil || data != nil {
			if err != nil {
				fmt.Println("1:" + err.Error())
			}
			continue
		}
		err = l.svcCtx.ClassesModel.Insert(l.ctx, &model.Classes{
			ID:          primitive.ObjectID{},
			UpdateAt:    time.Time{},
			CreateAt:    time.Time{},
			ClassesId:   id,
			ClassesName: split[1],
			ChineseName: "",
			EnglishName: "",
		})
		if err != nil {
			fmt.Println("2:" + err.Error())
		}
	}
	return nil
}
func InitClassesChineseName(l *InitClassesLogic) error {
	classesChineseListFile, err := os.Open("./txt/classes_Chinese_list.txt")
	if err != nil {
		return err
	}
	defer classesChineseListFile.Close()
	scanner := bufio.NewScanner(classesChineseListFile)
	for scanner.Scan() {
		split := strings.SplitN(scanner.Text(), " ", 2)
		if split == nil || len(split) != 2 {
			continue
		}
		id, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			continue
		}
		data, err := l.svcCtx.ClassesModel.FindOneByClassesId(l.ctx, id)
		if data == nil || err != nil {
			if err != nil {
				fmt.Println("3:" + err.Error())
			}
			continue
		}
		data.ChineseName = split[1]
		_, err = l.svcCtx.ClassesModel.Update(l.ctx, data)
		if err != nil {
			fmt.Println("4:" + err.Error())
		}
	}
	return nil
}
func InitClassesEnglishName(l *InitClassesLogic) error {
	classesEnglishListFile, err := os.Open("./txt/classes_English_list.txt")
	if err != nil {
		return err
	}
	defer classesEnglishListFile.Close()
	scanner := bufio.NewScanner(classesEnglishListFile)
	for scanner.Scan() {
		split := strings.SplitN(scanner.Text(), " ", 2)
		if split == nil || len(split) != 2 {
			continue
		}
		id, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			continue
		}
		data, err := l.svcCtx.ClassesModel.FindOneByClassesId(l.ctx, id)
		if data == nil || err != nil {
			if err != nil {
				fmt.Println("5:" + err.Error())
			}
			continue
		}
		data.EnglishName = split[1]
		_, err = l.svcCtx.ClassesModel.Update(l.ctx, data)
		if err != nil {
			fmt.Println("6:" + err.Error())
		}
	}
	return nil
}
