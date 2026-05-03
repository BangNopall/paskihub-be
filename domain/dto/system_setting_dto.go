package dto

import "github.com/BangNopall/paskihub-be/domain/entity"

type SystemSettingResponse struct {
	CoinRate    float64 `json:"coin_rate"`
	ApprovalFee float64 `json:"approval_fee"`
	BankInfo    struct {
		BankName      string `json:"bank_name"`
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
	} `json:"bank_info"`
}

type UpdateSystemSettingRequest struct {
	CoinRate      float64 `json:"coin_rate" validate:"required,min=1"`
	ApprovalFee   float64 `json:"approval_fee" validate:"required,min=0"`
	BankName      string  `json:"bank_name" validate:"required"`
	AccountNumber string  `json:"account_number" validate:"required"`
	AccountName   string  `json:"account_name" validate:"required"`
}

func SystemSettingEntityToResponse(setting *entity.SystemSetting) *SystemSettingResponse {
	res := &SystemSettingResponse{}
	res.CoinRate = setting.CoinRate
	res.ApprovalFee = setting.ApprovalFee
	res.BankInfo.BankName = setting.BankName
	res.BankInfo.AccountNumber = setting.AccountNumber
	res.BankInfo.AccountName = setting.AccountName
	return res
}
