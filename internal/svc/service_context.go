package svc

import (
	"github.com/iot-synergy/synergy-bird-rpc/internal/config"
	model "github.com/iot-synergy/synergy-bird-rpc/storage/gallery"
	model4 "github.com/iot-synergy/synergy-bird-rpc/storage/galleryCount"
	model2 "github.com/iot-synergy/synergy-bird-rpc/storage/illustration"
	model3 "github.com/iot-synergy/synergy-bird-rpc/storage/label"
	"github.com/iot-synergy/synergy-event-rpc/synergyeventclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	GalleryModel      model.GalleryModel
	IllustrationModel model2.IllustrationModel
	LabelModel        model3.LabelModel
	GalleryCountModel model4.GalleryCountModel
	EventRpc          synergyeventclient.SynergyEvent
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:            c,
		GalleryModel:      model.NewGalleryModel(c.MonDb.Url, c.MonDb.DbName, "gallery"),
		IllustrationModel: model2.NewIllustrationModel(c.MonDb.Url, c.MonDb.DbName, "illustration"),
		LabelModel:        model3.NewLabelModel(c.MonDb.Url, c.MonDb.DbName, "label"),
		GalleryCountModel: model4.NewGalleryCountModel(c.MonDb.Url, c.MonDb.DbName, "gallery_count"),
		EventRpc:          synergyeventclient.NewSynergyEvent(zrpc.MustNewClient(c.EventRpc)),
	}
}
