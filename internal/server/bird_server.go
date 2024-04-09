// Code generated by goctl. DO NOT EDIT.
// Source: bird.proto

package server

import (
	"context"

	"github.com/iot-synergy/synergy-bird-rpc/internal/logic/gallery"
	"github.com/iot-synergy/synergy-bird-rpc/internal/logic/illustration"
	"github.com/iot-synergy/synergy-bird-rpc/internal/logic/label"
	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
)

type BirdServer struct {
	svcCtx *svc.ServiceContext
	bird.UnimplementedBirdServer
}

func NewBirdServer(svcCtx *svc.ServiceContext) *BirdServer {
	return &BirdServer{
		svcCtx: svcCtx,
	}
}

func (s *BirdServer) GalleryCreate(ctx context.Context, in *bird.GalleryCreateReq) (*bird.GalleryResp, error) {
	l := gallery.NewGalleryCreateLogic(ctx, s.svcCtx)
	return l.GalleryCreate(in)
}

func (s *BirdServer) GalleryDelete(ctx context.Context, in *bird.IdReq) (*bird.GalleryResp, error) {
	l := gallery.NewGalleryDeleteLogic(ctx, s.svcCtx)
	return l.GalleryDelete(in)
}

func (s *BirdServer) GalleryList(ctx context.Context, in *bird.GalleryListReq) (*bird.GalleryListResp, error) {
	l := gallery.NewGalleryListLogic(ctx, s.svcCtx)
	return l.GalleryList(in)
}

func (s *BirdServer) IllustrationCreate(ctx context.Context, in *bird.IllustrationsCreateReq) (*bird.IllustrationsResp, error) {
	l := illustration.NewIllustrationCreateLogic(ctx, s.svcCtx)
	return l.IllustrationCreate(in)
}

func (s *BirdServer) IllustrationUpdate(ctx context.Context, in *bird.IllustrationsUpdateReq) (*bird.IllustrationsResp, error) {
	l := illustration.NewIllustrationUpdateLogic(ctx, s.svcCtx)
	return l.IllustrationUpdate(in)
}

func (s *BirdServer) IllustrationList(ctx context.Context, in *bird.IllustrationsListReq) (*bird.IllustrationsListResp, error) {
	l := illustration.NewIllustrationListLogic(ctx, s.svcCtx)
	return l.IllustrationList(in)
}

func (s *BirdServer) FindIllustrationByPage(ctx context.Context, in *bird.IllustrationsPageReq) (*bird.IllustrationsListVo, error) {
	l := illustration.NewFindIllustrationByPageLogic(ctx, s.svcCtx)
	return l.FindIllustrationByPage(in)
}

func (s *BirdServer) PublishIllustration(ctx context.Context, in *bird.IdReq) (*bird.IllustrationsResp, error) {
	l := illustration.NewPublishIllustrationLogic(ctx, s.svcCtx)
	return l.PublishIllustration(in)
}

func (s *BirdServer) UnpublishIllustration(ctx context.Context, in *bird.IdReq) (*bird.IllustrationsResp, error) {
	l := illustration.NewUnpublishIllustrationLogic(ctx, s.svcCtx)
	return l.UnpublishIllustration(in)
}

func (s *BirdServer) DeleteIllustration(ctx context.Context, in *bird.IdReq) (*bird.BooleanResp, error) {
	l := illustration.NewDeleteIllustrationLogic(ctx, s.svcCtx)
	return l.DeleteIllustration(in)
}

func (s *BirdServer) LabelCreate(ctx context.Context, in *bird.LabelCreateReq) (*bird.LabelResp, error) {
	l := label.NewLabelCreateLogic(ctx, s.svcCtx)
	return l.LabelCreate(in)
}

func (s *BirdServer) LabelUpdate(ctx context.Context, in *bird.LabelUpdateReq) (*bird.LabelResp, error) {
	l := label.NewLabelUpdateLogic(ctx, s.svcCtx)
	return l.LabelUpdate(in)
}

func (s *BirdServer) LabelList(ctx context.Context, in *bird.LabelListReq) (*bird.LabelListResp, error) {
	l := label.NewLabelListLogic(ctx, s.svcCtx)
	return l.LabelList(in)
}

func (s *BirdServer) FindLabelByPage(ctx context.Context, in *bird.LabelListReq) (*bird.LabelListVo, error) {
	l := label.NewFindLabelByPageLogic(ctx, s.svcCtx)
	return l.FindLabelByPage(in)
}

func (s *BirdServer) FindLabelById(ctx context.Context, in *bird.IdReq) (*bird.LabelVo, error) {
	l := label.NewFindLabelByIdLogic(ctx, s.svcCtx)
	return l.FindLabelById(in)
}

func (s *BirdServer) PublishLabel(ctx context.Context, in *bird.IdReq) (*bird.LabelResp, error) {
	l := label.NewPublishLabelLogic(ctx, s.svcCtx)
	return l.PublishLabel(in)
}

func (s *BirdServer) UnpublishLabel(ctx context.Context, in *bird.IdReq) (*bird.LabelResp, error) {
	l := label.NewUnpublishLabelLogic(ctx, s.svcCtx)
	return l.UnpublishLabel(in)
}

func (s *BirdServer) DeleteLabel(ctx context.Context, in *bird.IdReq) (*bird.BooleanResp, error) {
	l := label.NewDeleteLabelLogic(ctx, s.svcCtx)
	return l.DeleteLabel(in)
}
