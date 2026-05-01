package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"gorm.io/gorm"
)

type systemSettingRepository struct {
	db *gorm.DB
}

func NewSystemSettingRepository(db *gorm.DB) contracts.SystemSettingRepository {
	return &systemSettingRepository{
		db: db,
	}
}

func (r *systemSettingRepository) Get(ctx context.Context) (*entity.SystemSetting, error) {
	var setting entity.SystemSetting
	err := r.db.WithContext(ctx).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *systemSettingRepository) Update(ctx context.Context, setting *entity.SystemSetting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}
