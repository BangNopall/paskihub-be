package contracts

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
)

type SystemSettingRepository interface {
	Get(ctx context.Context) (*entity.SystemSetting, error)
	Update(ctx context.Context, setting *entity.SystemSetting) error
}

type SystemSettingService interface {
	GetPublicSettings(ctx context.Context) (*dto.SystemSettingResponse, error)
	GetSettings(ctx context.Context) (*dto.SystemSettingResponse, error)
	UpdateSettings(ctx context.Context, req *dto.UpdateSystemSettingRequest) error
}
