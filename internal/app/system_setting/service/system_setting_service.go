package service

import (
	"context"
	"time"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
)

type systemSettingService struct {
	settingRepo contracts.SystemSettingRepository
	timeout     time.Duration
}

func NewSystemSettingService(
	settingRepo contracts.SystemSettingRepository,
	timeout time.Duration,
) contracts.SystemSettingService {
	return &systemSettingService{
		settingRepo: settingRepo,
		timeout:     timeout,
	}
}

func (s *systemSettingService) GetPublicSettings(ctx context.Context) (*dto.SystemSettingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	setting, err := s.settingRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return dto.SystemSettingEntityToResponse(setting), nil
}

func (s *systemSettingService) GetSettings(ctx context.Context) (*dto.SystemSettingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	setting, err := s.settingRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return dto.SystemSettingEntityToResponse(setting), nil
}

func (s *systemSettingService) UpdateSettings(ctx context.Context, req *dto.UpdateSystemSettingRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	setting, err := s.settingRepo.Get(ctx)
	if err != nil {
		return err
	}

	setting.CoinRate = req.CoinRate
	setting.ApprovalFee = req.ApprovalFee
	setting.BankName = req.BankName
	setting.AccountNumber = req.AccountNumber
	setting.AccountName = req.AccountName

	return s.settingRepo.Update(ctx, setting)
}
