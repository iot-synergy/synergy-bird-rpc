package svc

import (
	"github.com/iot-synergy/synergy-bird-rpc/internal/config"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/gallery"
	model2 "github.com/iot-synergy/synergy-bird-rpc/storage/illustration"
	model3 "github.com/iot-synergy/synergy-bird-rpc/storage/label"
)

type ServiceContext struct {
	Config            config.Config
	GalleryModel      model.GalleryModel
	IllustrationModel model2.IllustrationModel
	LabelModel        model3.LabelModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:            c,
		GalleryModel:      model.NewGalleryModel(c.MonDb.Url, c.MonDb.DbName, "gallery"),
		IllustrationModel: model2.NewIllustrationModel(c.MonDb.Url, c.MonDb.DbName, "illustration"),
		LabelModel:        model3.NewLabelModel(c.MonDb.Url, c.MonDb.DbName, "label"),
	}
}
