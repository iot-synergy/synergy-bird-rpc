package headlines

import (
	"context"
	"testing"

	"github.com/iot-synergy/synergy-bird-rpc/internal/config"
	"github.com/iot-synergy/synergy-bird-rpc/internal/svc"
	headline_model "github.com/iot-synergy/synergy-bird-rpc/storage/headline"
	"github.com/iot-synergy/synergy-bird-rpc/types/bird"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestUpdateHeadlineLogic_UpdateHeadline(t *testing.T) {
	type args struct {
		in *bird.Headline
	}

	cover := false
	url := "3333333333"
	id := "66869e3f62a93dbc89213006"

	tests := []struct {
		name    string
		l       *UpdateHeadlineLogic
		args    args
		want    *bird.Headline
		wantErr bool
	}{
		{
			name: "",
			l: &UpdateHeadlineLogic{
				ctx: context.Background(),
				svcCtx: &svc.ServiceContext{
					Config:            config.Config{},
					HeadlineModel:     headline_model.NewHeadlineModel("mongodb://admin:123456@192.168.10.181:27017/", "peckperk", "headline"),
					GalleryModel:      nil,
					IllustrationModel: nil,
					LabelModel:        nil,
					GalleryCountModel: nil,
					ClassesModel:      nil,
					EventRpc:          nil,
				},
				Logger: logx.WithContext(context.Background()),
			},
			args: args{
				in: &bird.Headline{
					Id:    &id,
					Url:   &url,
					Cover: &cover,
				},
			},
			want:    &bird.Headline{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.UpdateHeadline(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateHeadlineLogic.UpdateHeadline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.l.Logger.Infow("UpdateHeadline success!", logx.LogField{Key: "result", Value: got})
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("UpdateHeadlineLogic.UpdateHeadline() = %v, want %v", got, tt.want)
			// }
		})
	}
}
